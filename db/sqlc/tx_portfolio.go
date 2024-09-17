package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
)

// CREATE
type PortfolioAsset struct {
	TickerId   int64   `json:"ticker_id,omitempty"`
	Allocation float64 `json:"allocation"`
	Price      float64 `json:"price"`
}

type CreatePortfolioTxParams struct {
	PortfolioID    string            `json:"portfolio_id"`
	CategoryID     string            `json:"category_id"`
	PortfolioName  string            `json:"portfolio_name"`
	OrganizationId []string          `json:"organization_id"`
	BranchId       []string          `json:"branch_id"`
	AdvisorId      []string          `json:"advisor_id"`
	Assets         []*PortfolioAsset `json:"assets"`
	Privacy        string            `json:"privacy"`
	AuthorID        string            `json:"author_id"`
}

type CreatePortfolioTxResult struct {
	PortfolioID string `json:"portfolio_id"`
}

func (store *SQLStore) CreatePortfolioTx(ctx context.Context, arg CreatePortfolioTxParams) (CreatePortfolioTxResult, error) {
	var result CreatePortfolioTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// table: portfolios
		argPortfolio := CreatePortfolioParams{
			ID:      arg.PortfolioID,
			Name:    arg.PortfolioName,
			Privacy: arg.Privacy,
			AuthorID: arg.AuthorID,
		}

		_, err = q.CreatePortfolio(ctx, argPortfolio)
		if err != nil {
			return errors.New("error CreatePortfolio " + err.Error())
		}

		// table: assets
		for _, asset := range arg.Assets {
			argAssest := CreateAssetParams{
				PortfolioID: arg.PortfolioID,
				TickerID:    int32(asset.TickerId),
				Price:       asset.Price,
				Allocation:  asset.Allocation,
			}

			_, err := q.CreateAsset(ctx, argAssest)
			if err != nil {
				return errors.New("error CreateAsset " + err.Error())
			}
		}

		// table: p_categories
		argPCategory := CreatePCategoryParams{
			PortfolioID: arg.PortfolioID,
			CategoryID: pgtype.Text{
				String: arg.CategoryID,
				Valid:  true,
			},
		}
		_, err = q.CreatePCategory(ctx, argPCategory)
		if err != nil {
			return errors.New("error CreatePCategory " + err.Error())
		}

		// table: p_branches
		for _, branch := range arg.BranchId {
			argPBranch := CreatePBranchParams{
				PortfolioID: arg.PortfolioID,
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

		// table: p_advisors
		for _, advisor := range arg.AdvisorId {
			argPAdvisor := CreatePAdvisorParams{
				PortfolioID: arg.PortfolioID,
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

		// table: p_organizations
		for _, organization := range arg.OrganizationId {
			argPOrganization := CreatePOrganizationParams{
				PortfolioID: arg.PortfolioID,
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

		return err
	})

	result.PortfolioID = arg.PortfolioID
	return result, err
}

// UPDATE
type UpdatePortfolioTxParams struct {
	PortfolioID    string            `json:"portfolio_id"`
	CategoryID     string            `json:"category_id"`
	PortfolioName  string            `json:"portfolio_name"`
	OrganizationId []string          `json:"organization_id"`
	BranchId       []string          `json:"branch_id"`
	AdvisorId      []string          `json:"advisor_id"`
	Assets         []*PortfolioAsset `json:"assets"`
	Privacy        string            `json:"privacy"`
}

type UpdatePortfolioTxResult struct {
	PortfolioID string `json:"portfolio_id"`
}

func findDifferences(cur, new []string) (add, remove []string) {
	// create map store "cur" element
	curMap := make(map[string]bool)
	for _, v := range cur {
		curMap[v] = true
	}

	// Finds elements in "new" but not in "cur"
	for _, v := range new {
		if _, ok := curMap[v]; !ok {
			add = append(add, v)
		}
	}

	// create map store "new" element
	newMap := make(map[string]bool)
	for _, v := range new {
		newMap[v] = true
	}

	// Finds elements in "cur" but not in "new"
	for _, v := range cur {
		if _, ok := newMap[v]; !ok {
			remove = append(remove, v)
		}
	}

	return add, remove
}

func (store *SQLStore) UpdatePortfolioTx(ctx context.Context, arg UpdatePortfolioTxParams) (UpdatePortfolioTxResult, error) {
	var result UpdatePortfolioTxResult

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

		// table: p_categories
		argPCategory := UpdatePCategoryParams{
			PortfolioID: arg.PortfolioID,
			CategoryID: pgtype.Text{
				String: arg.CategoryID,
				Valid:  true,
			},
		}
		_, err = q.UpdatePCategory(ctx, argPCategory)
		if err != nil {
			return errors.New("error UpdatePCategory " + err.Error())
		}

		// table: p_branches
		branches, err := q.GetPBranchByPortfolioId(ctx, arg.PortfolioID)
		if err != nil {
			return errors.New("error GetPBranchByPortfolioId " + err.Error())
		}

		currentBranches := make([]string, len(branches))
		// convert branches to currentBranches
		for i, b := range branches {
			currentBranches[i] = b.BranchID.String
		}

		add, remove := findDifferences(currentBranches, arg.BranchId)
		if len(add) > 0 {
			for _, value := range add {
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

		if len(remove) > 0 {
			for _, value := range remove {
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
		advisors, err := q.GetPAdvisorByPortfolioId(ctx, arg.PortfolioID)
		if err != nil {
			return errors.New("error GetPAdvisorByPortfolioId " + err.Error())
		}

		currentAdvisors := make([]string, len(advisors))
		// convert advisors to currentAdvisors
		for i, a := range advisors {
			currentAdvisors[i] = a.AdvisorID.String
		}

		add, remove = findDifferences(currentAdvisors, arg.AdvisorId)
		if len(add) > 0 {
			for _, value := range add {
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

		if len(remove) > 0 {
			for _, value := range remove {
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
		organizations, err := q.GetPOrganizationByPortfolioId(ctx, arg.PortfolioID)
		if err != nil {
			return errors.New("error GetPOrganizationByPortfolioId " + err.Error())
		}

		currentOrganizations := make([]string, len(organizations))
		// convert organizations to currentOrganizations
		for i, o := range organizations {
			currentOrganizations[i] = o.OrganizationID.String
		}

		add, remove = findDifferences(currentOrganizations, arg.OrganizationId)
		if len(add) > 0 {
			for _, value := range add {
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

		if len(remove) > 0 {
			for _, value := range remove {
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

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// table: portfolios
		err = q.DeletePortfolio(ctx, arg.PortfolioID)
		if err != nil {
			return errors.New("error DeletePortfolio " + err.Error())
		}

		// table: assets
		assets, err := q.GetAssetsByPortfolioId(ctx, arg.PortfolioID)
		if err != nil {
			return errors.New("error GetAssetsByPortfolioId " + err.Error())
		}
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

		// table: p_categories
		categories, err := q.GetPCategoryByPortfolioId(ctx, arg.PortfolioID)
		if err != nil {
			return errors.New("error GetPCategoryByPortfolioId " + err.Error())
		}

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

		// table: p_branches
		branches, err := q.GetPBranchByPortfolioId(ctx, arg.PortfolioID)
		if err != nil {
			return errors.New("error GetPBranchByPortfolioId " + err.Error())
		}

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

		// table: p_advisors
		advisors, err := q.GetPAdvisorByPortfolioId(ctx, arg.PortfolioID)
		if err != nil {
			return errors.New("error GetPAdvisorByPortfolioId " + err.Error())
		}

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

		// table: p_organizations
		organizations, err := q.GetPOrganizationByPortfolioId(ctx, arg.PortfolioID)
		if err != nil {
			return errors.New("error GetPOrganizationByPortfolioId " + err.Error())
		}

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

		return err
	})

	result.PortfolioID = arg.PortfolioID
	return result, err
}
