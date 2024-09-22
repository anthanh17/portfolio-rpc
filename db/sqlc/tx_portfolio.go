package db

import (
	"context"
	"errors"
	"portfolio-profile-rpc/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// CREATE
type PortfolioAsset struct {
	TickerId   int64   `json:"ticker_id,omitempty"`
	Allocation float64 `json:"allocation"`
	Price      float64 `json:"price"`
}

type CreatePortfolioTxParams struct {
	CategoryID     []string          `json:"category_id"`
	PortfolioName  string            `json:"portfolio_name"`
	OrganizationId []string          `json:"organization_id"`
	BranchId       []string          `json:"branch_id"`
	AdvisorId      []string          `json:"advisor_id"`
	Assets         []*PortfolioAsset `json:"assets"`
	Privacy        string            `json:"privacy"`
	AuthorID       string            `json:"author_id"`
}

type CreatePortfolioTxResult struct {
	PortfolioID string `json:"portfolio_id"`
}

func (store *SQLStore) CreatePortfolioTx(ctx context.Context, arg CreatePortfolioTxParams) (CreatePortfolioTxResult, error) {
	var result CreatePortfolioTxResult
	portfolioId := uuid.New().String()

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// table: u_portfolio
		argUPortfolio := CreateUserPortfolioParams{
			UserID: arg.AuthorID,
			PortfolioID: pgtype.Text{
				String: portfolioId,
				Valid:  true,
			},
		}
		_, err = q.CreateUserPortfolio(ctx, argUPortfolio)
		if err != nil {
			return errors.New("error CreateUserPortfolio " + err.Error())
		}

		// table: portfolios
		argPortfolio := CreatePortfolioParams{
			ID:       portfolioId,
			Name:     arg.PortfolioName,
			Privacy:  arg.Privacy,
			AuthorID: arg.AuthorID,
		}

		_, err = q.CreatePortfolio(ctx, argPortfolio)
		if err != nil {
			return errors.New("error CreatePortfolio " + err.Error())
		}

		// table: assets
		if len(arg.Assets) > 0 {
			for _, asset := range arg.Assets {
				argAssest := CreateAssetParams{
					PortfolioID: portfolioId,
					TickerID:    int32(asset.TickerId),
					Price:       asset.Price,
					Allocation:  asset.Allocation,
				}

				_, err := q.CreateAsset(ctx, argAssest)
				if err != nil {
					return errors.New("error CreateAsset " + err.Error())
				}
			}
		}

		// table: p_categories
		if len(arg.CategoryID) > 0 {
			for _, category := range arg.CategoryID {
				argPCategory := CreatePCategoryParams{
					PortfolioID: portfolioId,
					CategoryID: pgtype.Text{
						String: category,
						Valid:  true,
					},
				}
				_, err = q.CreatePCategory(ctx, argPCategory)
				if err != nil {
					return errors.New("error CreatePCategory " + err.Error())
				}
			}
		}

		// table: p_branches
		if len(arg.BranchId) > 0 {
			for _, branch := range arg.BranchId {
				argPBranch := CreatePBranchParams{
					PortfolioID: portfolioId,
					BranchID: pgtype.Text{
						String: branch,
						Valid:  true,
					},
				}
				_, err = q.CreatePBranch(ctx, argPBranch)
				if err != nil {
					return errors.New("error CreatePBranch " + err.Error())
				}
			}
		}

		// table: p_advisors
		if len(arg.AdvisorId) > 0 {
			for _, advisor := range arg.AdvisorId {
				argPAdvisor := CreatePAdvisorParams{
					PortfolioID: portfolioId,
					AdvisorID: pgtype.Text{
						String: advisor,
						Valid:  true,
					},
				}
				_, err = q.CreatePAdvisor(ctx, argPAdvisor)
				if err != nil {
					return errors.New("error CreatePAdvisor " + err.Error())
				}
			}
		}

		// table: p_organizations
		if len(arg.OrganizationId) > 0 {
			for _, organization := range arg.OrganizationId {
				argPOrganization := CreatePOrganizationParams{
					PortfolioID: portfolioId,
					OrganizationID: pgtype.Text{
						String: organization,
						Valid:  true,
					},
				}
				_, err = q.CreatePOrganization(ctx, argPOrganization)
				if err != nil {
					return errors.New("error CreatePOrganization " + err.Error())
				}
			}
		}

		return err
	})

	result.PortfolioID = portfolioId
	return result, err
}

// UPDATE
type UpdatePortfolioTxParams struct {
	PortfolioID    string            `json:"portfolio_id"`
	PortfolioName  string            `json:"portfolio_name"`
	CategoryID     []string          `json:"category_id"`
	OrganizationId []string          `json:"organization_id"`
	BranchId       []string          `json:"branch_id"`
	AdvisorId      []string          `json:"advisor_id"`
	Assets         []*PortfolioAsset `json:"assets"`
	Privacy        string            `json:"privacy"`
}

type UpdatePortfolioTxResult struct {
	PortfolioID string `json:"portfolio_id"`
}

func (store *SQLStore) UpdatePortfolioTx(ctx context.Context, arg UpdatePortfolioTxParams) (UpdatePortfolioTxResult, error) {
	var result UpdatePortfolioTxResult

	/*
		- error channel hold err:
			- p_categories
			- p_branches
			- p_advisors
			- p_organizations
	*/
	errGet := make(chan error, 4)

	// GetPCategoryByPortfolioId
	addPCategoriesCh := make(chan []string)
	removePCategoriesCh := make(chan []string)
	go func() {
		categories, err := store.GetPCategoryByPortfolioId(ctx, arg.PortfolioID)
		errGet <- err

		currentCategories := make([]string, len(categories))
		// convert categories to currentCategories
		for i, c := range categories {
			currentCategories[i] = c.CategoryID.String
		}

		add, remove := util.FindDifferences(currentCategories, arg.CategoryID)

		// push data to channel
		addPCategoriesCh <- add
		removePCategoriesCh <- remove
	}()

	// GetPBranchByPortfolioId
	addPBranchesCh := make(chan []string)
	removePBranchesCh := make(chan []string)
	go func() {
		branches, err := store.GetPBranchByPortfolioId(ctx, arg.PortfolioID)
		errGet <- err

		currentBranches := make([]string, len(branches))
		// convert branches to currentBranches
		for i, b := range branches {
			currentBranches[i] = b.BranchID.String
		}

		add, remove := util.FindDifferences(currentBranches, arg.BranchId)
		// push data to channel
		addPBranchesCh <- add
		removePBranchesCh <- remove
	}()

	// GetPAdvisorByPortfolioId
	addPAdvisorsCh := make(chan []string)
	removePAdvisorsCh := make(chan []string)
	go func() {
		advisors, err := store.GetPAdvisorByPortfolioId(ctx, arg.PortfolioID)
		errGet <- err

		currentAdvisors := make([]string, len(advisors))
		// convert advisors to currentAdvisors
		for i, a := range advisors {
			currentAdvisors[i] = a.AdvisorID.String
		}

		add, remove := util.FindDifferences(currentAdvisors, arg.AdvisorId)
		// push data to channel
		addPAdvisorsCh <- add
		removePAdvisorsCh <- remove
	}()

	// GetPOrganizationByPortfolioId
	addPOrganizationsCh := make(chan []string)
	removePOrganizationsCh := make(chan []string)
	go func() {
		organizations, err := store.GetPOrganizationByPortfolioId(ctx, arg.PortfolioID)
		errGet <- err

		currentOrganizations := make([]string, len(organizations))
		// convert organizations to currentOrganizations
		for i, o := range organizations {
			currentOrganizations[i] = o.OrganizationID.String
		}

		add, remove := util.FindDifferences(currentOrganizations, arg.OrganizationId)
		// push data to channel
		addPOrganizationsCh <- add
		removePOrganizationsCh <- remove
	}()

	addPCategories, removePCategories := <-addPCategoriesCh, <-removePCategoriesCh
	addPBranches, removePBranches := <-addPBranchesCh, <-removePBranchesCh
	addPAdvisors, removePAdvisors := <-addPAdvisorsCh, <-removePAdvisorsCh
	addPOrganizations, removePOrganizations := <-addPOrganizationsCh, <-removePOrganizationsCh

	close(errGet)
	// Collect and handle errors
	for err := range errGet {
		if err != nil {
			return result, err
		}
	}

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// table: portfolios
		argPortfolio := UpdatePortfolioParams{
			ID:      arg.PortfolioID,
			Name:    arg.PortfolioName,
			Privacy: arg.Privacy,
		}

		_, err = q.UpdatePortfolio(ctx, argPortfolio)
		if err != nil {
			return errors.New("error UpdatePortfolio " + err.Error())
		}

		// table: assets
		if len(arg.Assets) > 0 {
			for _, asset := range arg.Assets {
				argAssest := UpdateAssetParams{
					PortfolioID: arg.PortfolioID,
					TickerID:    int32(asset.TickerId),
					Price:       asset.Price,
					Allocation:  asset.Allocation,
				}

				_, err := q.UpdateAsset(ctx, argAssest)
				if err != nil {
					return errors.New("error UpdateAsset " + err.Error())
				}
			}
		}

		// table: p_categories
		if len(addPCategories) > 0 {
			for _, value := range addPCategories {
				argPOrganization := CreatePCategoryParams{
					PortfolioID: arg.PortfolioID,
					CategoryID: pgtype.Text{
						String: value,
						Valid:  true,
					},
				}
				_, err = q.CreatePCategory(ctx, argPOrganization)
				if err != nil {
					return errors.New("error UpdatePCategory - CreatePCategory " + err.Error())
				}
			}
		}

		if len(removePCategories) > 0 {
			for _, value := range removePCategories {
				argPOrganization := DeletePCategoryParams{
					PortfolioID: arg.PortfolioID,
					CategoryID: pgtype.Text{
						String: value,
						Valid:  true,
					},
				}

				err = q.DeletePCategory(ctx, argPOrganization)
				if err != nil {
					return errors.New("error UpdatePCategory - DeletePCategory " + err.Error())
				}
			}
		}

		// table: p_branches
		if len(addPBranches) > 0 {
			for _, value := range addPBranches {
				argPBranch := CreatePBranchParams{
					PortfolioID: arg.PortfolioID,
					BranchID: pgtype.Text{
						String: value,
						Valid:  true,
					},
				}

				_, err = q.CreatePBranch(ctx, argPBranch)
				if err != nil {
					return errors.New("error UpdatePBranch - CreatePBranch " + err.Error())
				}
			}
		}

		if len(removePBranches) > 0 {
			for _, value := range removePBranches {
				argPBranch := DeletePBranchParams{
					PortfolioID: arg.PortfolioID,
					BranchID: pgtype.Text{
						String: value,
						Valid:  true,
					},
				}

				err = q.DeletePBranch(ctx, argPBranch)
				if err != nil {
					return errors.New("error UpdatePBranch - DeletePBranch " + err.Error())
				}
			}
		}

		// table: p_advisors
		if len(addPAdvisors) > 0 {
			for _, value := range addPAdvisors {
				argPAdvisor := CreatePAdvisorParams{
					PortfolioID: arg.PortfolioID,
					AdvisorID: pgtype.Text{
						String: value,
						Valid:  true,
					},
				}

				_, err = q.CreatePAdvisor(ctx, argPAdvisor)
				if err != nil {
					return errors.New("error UpdatePAdvisor - CreatePAdvisor " + err.Error())
				}
			}
		}

		if len(removePAdvisors) > 0 {
			for _, value := range removePAdvisors {
				argPAdvisor := DeletePAdvisorParams{
					PortfolioID: arg.PortfolioID,
					AdvisorID: pgtype.Text{
						String: value,
						Valid:  true,
					},
				}

				err = q.DeletePAdvisor(ctx, argPAdvisor)
				if err != nil {
					return errors.New("error UpdatePAdvisor - DeletePAdvisor " + err.Error())
				}
			}
		}

		// table: p_organizations
		if len(addPOrganizations) > 0 {
			for _, value := range addPOrganizations {
				argPOrganization := CreatePOrganizationParams{
					PortfolioID: arg.PortfolioID,
					OrganizationID: pgtype.Text{
						String: value,
						Valid:  true,
					},
				}

				_, err = q.CreatePOrganization(ctx, argPOrganization)
				if err != nil {
					return errors.New("error UpdatePOrganization - CreatePOrganization " + err.Error())
				}
			}
		}

		if len(removePOrganizations) > 0 {
			for _, value := range removePOrganizations {
				argPOrganization := DeletePOrganizationParams{
					PortfolioID: arg.PortfolioID,
					OrganizationID: pgtype.Text{
						String: value,
						Valid:  true,
					},
				}

				err = q.DeletePOrganization(ctx, argPOrganization)
				if err != nil {
					return errors.New("error UpdatePOrganization - DeletePOrganization " + err.Error())
				}
			}
		}

		return err
	})

	result.PortfolioID = arg.PortfolioID
	return result, err
}

// DELETE
type DeletePortfolioTxParams struct {
	PortfolioID string `json:"portfolio_id"`
}

type DeletePortfolioTxResult struct {
	PortfolioID string `json:"portfolio_id"`
}

func (store *SQLStore) DeletePortfolioTx(ctx context.Context, arg DeletePortfolioTxParams) (DeletePortfolioTxResult, error) {
	var result DeletePortfolioTxResult

	/*
		- assets
		- p_categories
		- p_branches
		- p_advisors
		- p_organizations
	*/
	errGet := make(chan error, 5)

	// GetAssetsByPortfolioId
	assetsCh := make(chan []HamonixBusinessAsset)
	go func() {
		assets, err := store.GetAssetsByPortfolioId(ctx, arg.PortfolioID)
		errGet <- err

		// push data to channel
		assetsCh <- assets
	}()

	// GetPCategoryByPortfolioId
	categoriesCh := make(chan []HamonixBusinessPCategory)
	go func() {
		categories, err := store.GetPCategoryByPortfolioId(ctx, arg.PortfolioID)
		errGet <- err

		// push data to channel
		categoriesCh <- categories
	}()

	// GetPBranchByPortfolioId
	branchesCh := make(chan []HamonixBusinessPBranch)
	go func() {
		branches, err := store.GetPBranchByPortfolioId(ctx, arg.PortfolioID)
		errGet <- err

		// push data to channel
		branchesCh <- branches
	}()

	// GetPAdvisorByPortfolioId
	advisorsCh := make(chan []HamonixBusinessPAdvisor)
	go func() {
		advisors, err := store.GetPAdvisorByPortfolioId(ctx, arg.PortfolioID)
		errGet <- err

		// push data to channel
		advisorsCh <- advisors
	}()

	// GetPOrganizationByPortfolioId
	organizationsCh := make(chan []HamonixBusinessPOrganization)
	go func() {
		organizations, err := store.GetPOrganizationByPortfolioId(ctx, arg.PortfolioID)
		errGet <- err

		// push data to channel
		organizationsCh <- organizations
	}()

	assets := <-assetsCh
	categories := <-categoriesCh
	branches := <-branchesCh
	advisors := <-advisorsCh
	organizations := <-organizationsCh

	close(errGet)
	// Collect and handle errors
	for err := range errGet {
		if err != nil {
			return result, err
		}
	}

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// table: portfolios
		err = q.DeletePortfolio(ctx, arg.PortfolioID)
		if err != nil {
			return errors.New("error DeletePortfolio " + err.Error())
		}

		// table: assets
		if len(assets) > 0 {
			for _, asset := range assets {
				argAssest := DeleteAssetParams{
					PortfolioID: asset.PortfolioID,
					TickerID:    int32(asset.TickerID),
				}

				err := q.DeleteAsset(ctx, argAssest)
				if err != nil {
					return errors.New("error DeleteAsset " + err.Error())
				}
			}
		}

		// table: p_categories
		if len(categories) > 0 {
			for _, category := range categories {
				argPCategory := DeletePCategoryParams{
					PortfolioID: category.PortfolioID,
					CategoryID: pgtype.Text{
						String: category.CategoryID.String,
						Valid:  true,
					},
				}

				err = q.DeletePCategory(ctx, argPCategory)
				if err != nil {
					return errors.New("error DeletePCategory " + err.Error())
				}
			}
		}

		// table: p_branches
		if len(branches) > 0 {
			for _, branch := range branches {
				argPBranch := DeletePBranchParams{
					PortfolioID: branch.PortfolioID,
					BranchID: pgtype.Text{
						String: branch.BranchID.String,
						Valid:  true,
					},
				}

				err = q.DeletePBranch(ctx, argPBranch)
				if err != nil {
					return errors.New("error DeletePBranch " + err.Error())
				}
			}
		}

		// table: p_advisors
		if len(advisors) > 0 {
			for _, advisor := range advisors {
				argPAdvisor := DeletePAdvisorParams{
					PortfolioID: advisor.PortfolioID,
					AdvisorID: pgtype.Text{
						String: advisor.AdvisorID.String,
						Valid:  true,
					},
				}

				err = q.DeletePAdvisor(ctx, argPAdvisor)
				if err != nil {
					return errors.New("error DeletePAdvisor " + err.Error())
				}
			}
		}

		// table: p_organizations
		if len(organizations) > 0 {
			for _, organization := range organizations {
				argPOrganization := DeletePOrganizationParams{
					PortfolioID: organization.PortfolioID,
					OrganizationID: pgtype.Text{
						String: organization.OrganizationID.String,
						Valid:  true,
					},
				}

				err = q.DeletePOrganization(ctx, argPOrganization)
				if err != nil {
					return errors.New("error DeletePOrganization " + err.Error())
				}
			}
		}

		return err
	})

	result.PortfolioID = arg.PortfolioID
	return result, err
}
