package db

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

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
		// Start goroutines
		var wg sync.WaitGroup
		errBuffer := 1 +
			len(arg.Assets) +
			len(arg.CategoryID) +
			len(arg.BranchId) +
			len(arg.AdvisorId) +
			len(arg.OrganizationId)

		errCh := make(chan error, errBuffer)

		// table: portfolios
		wg.Add(1)
		go func() {
			defer wg.Done()
			argPortfolio := CreatePortfolioParams{
				ID:       portfolioId,
				Name:     arg.PortfolioName,
				Privacy:  arg.Privacy,
				AuthorID: arg.AuthorID,
			}

			_, err := q.CreatePortfolio(ctx, argPortfolio)
			errCh <- err // errors.New("error CreatePortfolio " + err.Error())
		}()

		// table: assets
		if len(arg.Assets) > 0 {
			for _, asset := range arg.Assets {
				wg.Add(1)
				go func() {
					defer wg.Done()
					argAssest := CreateAssetParams{
						PortfolioID: portfolioId,
						TickerID:    int32(asset.TickerId),
						Price:       asset.Price,
						Allocation:  asset.Allocation,
					}

					_, err := q.CreateAsset(ctx, argAssest)
					errCh <- err // errors.New("error CreateAsset " + err.Error())
				}()
			}
		}

		// // table: p_categories
		// if len(arg.CategoryID) > 0 {
		// 	for _, category := range arg.CategoryID {
		// 		wg.Add(1)
		// 		go func() {
		// 			defer wg.Done()
		// 			argPCategory := CreatePCategoryParams{
		// 				PortfolioID: portfolioId,
		// 				CategoryID: pgtype.Text{
		// 					String: category,
		// 					Valid:  true,
		// 				},
		// 			}
		// 			_, err := q.CreatePCategory(ctx, argPCategory)
		// 			errCh <- err // errors.New("error CreatePCategory " + err.Error())
		// 		}()
		// 	}
		// }

		// // table: p_branches
		// if len(arg.BranchId) > 0 {
		// 	for _, branch := range arg.BranchId {
		// 		wg.Add(1)
		// 		go func() {
		// 			defer wg.Done()
		// 			argPBranch := CreatePBranchParams{
		// 				PortfolioID: portfolioId,
		// 				BranchID: pgtype.Text{
		// 					String: branch,
		// 					Valid:  true,
		// 				},
		// 			}
		// 			_, err := q.CreatePBranch(ctx, argPBranch)
		// 			errCh <- err // errors.New("error CreatePBranch " + err.Error())
		// 		}()
		// 	}
		// }

		// // table: p_advisors
		// if len(arg.AdvisorId) > 0 {
		// 	for _, advisor := range arg.AdvisorId {
		// 		wg.Add(1)
		// 		go func() {
		// 			defer wg.Done()
		// 			argPAdvisor := CreatePAdvisorParams{
		// 				PortfolioID: portfolioId,
		// 				AdvisorID: pgtype.Text{
		// 					String: advisor,
		// 					Valid:  true,
		// 				},
		// 			}
		// 			_, err := q.CreatePAdvisor(ctx, argPAdvisor)
		// 			errCh <- err // errors.New("error CreatePAdvisor " + err.Error())
		// 		}()
		// 	}
		// }

		// // table: p_organizations
		// if len(arg.OrganizationId) > 0 {
		// 	for _, organization := range arg.OrganizationId {
		// 		wg.Add(1)
		// 		go func() {
		// 			defer wg.Done()
		// 			argPOrganization := CreatePOrganizationParams{
		// 				PortfolioID: portfolioId,
		// 				OrganizationID: pgtype.Text{
		// 					String: organization,
		// 					Valid:  true,
		// 				},
		// 			}
		// 			_, err := q.CreatePOrganization(ctx, argPOrganization)
		// 			errCh <- err // errors.New("error CreatePOrganization " + err.Error())
		// 		}()
		// 	}
		// }

		// Wait for all goroutines to finish
		wg.Wait()
		close(errCh)

		// Collect and handle errors
		for err := range errCh {
			if err != nil {
				return err
			}
		}

		return nil
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

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

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
			categories, err := q.GetPCategoryByPortfolioId(ctx, arg.PortfolioID)
			errGet <- errors.New("error GetPCategoryByPortfolioId " + err.Error())

			currentCategories := make([]string, len(categories))
			// convert categories to currentCategories
			for i, c := range categories {
				currentCategories[i] = c.CategoryID.String
			}

			add, remove := findDifferences(currentCategories, arg.CategoryID)

			// push data to channel
			addPCategoriesCh <- add
			removePCategoriesCh <- remove
		}()

		// GetPBranchByPortfolioId
		addPBranchesCh := make(chan []string)
		removePBranchesCh := make(chan []string)
		go func() {
			branches, err := q.GetPBranchByPortfolioId(ctx, arg.PortfolioID)
			errGet <- errors.New("error GetPBranchByPortfolioId " + err.Error())

			currentBranches := make([]string, len(branches))
			// convert branches to currentBranches
			for i, b := range branches {
				currentBranches[i] = b.BranchID.String
			}

			add, remove := findDifferences(currentBranches, arg.BranchId)
			// push data to channel
			addPBranchesCh <- add
			removePBranchesCh <- remove
		}()

		// GetPAdvisorByPortfolioId
		addPAdvisorsCh := make(chan []string)
		removePAdvisorsCh := make(chan []string)
		go func() {
			advisors, err := q.GetPAdvisorByPortfolioId(ctx, arg.PortfolioID)
			errGet <- errors.New("error GetPAdvisorByPortfolioId " + err.Error())

			currentAdvisors := make([]string, len(advisors))
			// convert advisors to currentAdvisors
			for i, a := range advisors {
				currentAdvisors[i] = a.AdvisorID.String
			}

			add, remove := findDifferences(currentAdvisors, arg.AdvisorId)
			// push data to channel
			addPAdvisorsCh <- add
			removePAdvisorsCh <- remove
		}()

		// GetPOrganizationByPortfolioId
		addPOrganizationsCh := make(chan []string)
		removePOrganizationsCh := make(chan []string)
		go func() {
			organizations, err := q.GetPOrganizationByPortfolioId(ctx, arg.PortfolioID)
			errGet <- errors.New("error GetPOrganizationByPortfolioId " + err.Error())

			currentOrganizations := make([]string, len(organizations))
			// convert organizations to currentOrganizations
			for i, o := range organizations {
				currentOrganizations[i] = o.OrganizationID.String
			}

			add, remove := findDifferences(currentOrganizations, arg.OrganizationId)
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
				return err
			}
		}

		// -------------   Start goroutines --------
		var wg sync.WaitGroup
		errBuffer := 1 +
			len(arg.Assets) +
			len(addPCategories) +
			len(removePCategories) +
			len(addPBranches) +
			len(removePBranches) +
			len(addPAdvisors) +
			len(removePAdvisors) +
			len(addPOrganizations) +
			len(removePOrganizations)

		errCh := make(chan error, errBuffer)

		// table: portfolios
		wg.Add(1)
		go func() {
			defer wg.Done()
			argPortfolio := UpdatePortfolioParams{
				ID:      arg.PortfolioID,
				Name:    arg.PortfolioName,
				Privacy: arg.Privacy,
			}

			_, err = q.UpdatePortfolio(ctx, argPortfolio)
			errCh <- errors.New("error UpdatePortfolio " + err.Error())
		}()

		// table: assets
		if len(arg.Assets) > 0 {
			for _, asset := range arg.Assets {
				wg.Add(1)
				go func() {
					defer wg.Done()
					argAssest := UpdateAssetParams{
						PortfolioID: arg.PortfolioID,
						TickerID:    int32(asset.TickerId),
						Price:       asset.Price,
						Allocation:  asset.Allocation,
					}

					_, err := q.UpdateAsset(ctx, argAssest)
					errCh <- errors.New("error UpdateAsset " + err.Error())
				}()
			}
		}

		// table: p_categories
		if len(addPCategories) > 0 {
			for _, value := range addPCategories {
				wg.Add(1)
				go func() {
					defer wg.Done()
					argPOrganization := CreatePCategoryParams{
						PortfolioID: arg.PortfolioID,
						CategoryID: pgtype.Text{
							String: value,
							Valid:  true,
						},
					}
					_, err = q.CreatePCategory(ctx, argPOrganization)
					errCh <- errors.New("error UpdatePCategory - CreatePCategory " + err.Error())
				}()
			}
		}

		if len(removePCategories) > 0 {
			for _, value := range removePCategories {
				wg.Add(1)
				go func() {
					defer wg.Done()
					argPOrganization := DeletePCategoryParams{
						PortfolioID: arg.PortfolioID,
						CategoryID: pgtype.Text{
							String: value,
							Valid:  true,
						},
					}
					err = q.DeletePCategory(ctx, argPOrganization)
					errCh <- errors.New("error UpdatePCategory - DeletePCategory " + err.Error())
				}()
			}
		}

		// table: p_branches
		if len(addPBranches) > 0 {
			for _, value := range addPBranches {
				wg.Add(1)
				go func() {
					defer wg.Done()
					argPBranch := CreatePBranchParams{
						PortfolioID: arg.PortfolioID,
						BranchID: pgtype.Text{
							String: value,
							Valid:  true,
						},
					}
					_, err = q.CreatePBranch(ctx, argPBranch)
					errCh <- errors.New("error UpdatePBranch - CreatePBranch " + err.Error())
				}()
			}
		}

		if len(removePBranches) > 0 {
			for _, value := range removePBranches {
				wg.Add(1)
				go func() {
					defer wg.Done()
					argPBranch := DeletePBranchParams{
						PortfolioID: arg.PortfolioID,
						BranchID: pgtype.Text{
							String: value,
							Valid:  true,
						},
					}
					err = q.DeletePBranch(ctx, argPBranch)
					errCh <- errors.New("error UpdatePBranch - DeletePBranch " + err.Error())
				}()
			}
		}

		// table: p_advisors
		if len(addPAdvisors) > 0 {
			for _, value := range addPAdvisors {
				wg.Add(1)
				go func() {
					defer wg.Done()
					argPAdvisor := CreatePAdvisorParams{
						PortfolioID: arg.PortfolioID,
						AdvisorID: pgtype.Text{
							String: value,
							Valid:  true,
						},
					}
					_, err = q.CreatePAdvisor(ctx, argPAdvisor)
					errCh <- errors.New("error UpdatePAdvisor - CreatePAdvisor " + err.Error())
				}()
			}
		}

		if len(removePAdvisors) > 0 {
			for _, value := range removePAdvisors {
				wg.Add(1)
				go func() {
					defer wg.Done()
					argPAdvisor := DeletePAdvisorParams{
						PortfolioID: arg.PortfolioID,
						AdvisorID: pgtype.Text{
							String: value,
							Valid:  true,
						},
					}
					err = q.DeletePAdvisor(ctx, argPAdvisor)
					errCh <- errors.New("error UpdatePAdvisor - DeletePAdvisor " + err.Error())
				}()
			}
		}

		// table: p_organizations
		if len(addPOrganizations) > 0 {
			for _, value := range addPOrganizations {
				wg.Add(1)
				go func() {
					defer wg.Done()
					argPOrganization := CreatePOrganizationParams{
						PortfolioID: arg.PortfolioID,
						OrganizationID: pgtype.Text{
							String: value,
							Valid:  true,
						},
					}
					_, err = q.CreatePOrganization(ctx, argPOrganization)
					errCh <- errors.New("error UpdatePOrganization - CreatePOrganization " + err.Error())
				}()
			}
		}

		if len(removePOrganizations) > 0 {
			for _, value := range removePOrganizations {
				wg.Add(1)
				go func() {
					defer wg.Done()
					argPOrganization := DeletePOrganizationParams{
						PortfolioID: arg.PortfolioID,
						OrganizationID: pgtype.Text{
							String: value,
							Valid:  true,
						},
					}
					err = q.DeletePOrganization(ctx, argPOrganization)
					errCh <- errors.New("error UpdatePOrganization - DeletePOrganization " + err.Error())
				}()
			}
		}

		// Wait for all goroutines to finish
		wg.Wait()
		close(errCh)

		// Collect and handle errors
		for err := range errCh {
			if err != nil {
				return err
			}
		}

		return nil
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
			assets, err := q.GetAssetsByPortfolioId(ctx, arg.PortfolioID)
			errGet <- errors.New("error GetAssetsByPortfolioId " + err.Error())

			// push data to channel
			assetsCh <- assets
		}()

		// GetPCategoryByPortfolioId
		categoriesCh := make(chan []HamonixBusinessPCategory)
		go func() {
			categories, err := q.GetPCategoryByPortfolioId(ctx, arg.PortfolioID)
			errGet <- errors.New("error GetPCategoryByPortfolioId " + err.Error())

			// push data to channel
			categoriesCh <- categories
		}()

		// GetPBranchByPortfolioId
		branchesCh := make(chan []HamonixBusinessPBranch)
		go func() {
			branches, err := q.GetPBranchByPortfolioId(ctx, arg.PortfolioID)
			errGet <- errors.New("error GetPBranchByPortfolioId " + err.Error())

			// push data to channel
			branchesCh <- branches
		}()

		// GetPAdvisorByPortfolioId
		advisorsCh := make(chan []HamonixBusinessPAdvisor)
		go func() {
			advisors, err := q.GetPAdvisorByPortfolioId(ctx, arg.PortfolioID)
			errGet <- errors.New("error GetPAdvisorByPortfolioId " + err.Error())

			// push data to channel
			advisorsCh <- advisors
		}()

		// GetPOrganizationByPortfolioId
		organizationsCh := make(chan []HamonixBusinessPOrganization)
		go func() {
			organizations, err := q.GetPOrganizationByPortfolioId(ctx, arg.PortfolioID)
			errGet <- errors.New("error GetPOrganizationByPortfolioId " + err.Error())

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
				return err
			}
		}

		// -------------   Start goroutines --------
		var wg sync.WaitGroup
		errBuffer := 1 +
			len(assets) +
			len(categories) +
			len(branches) +
			len(advisors) +
			len(organizations)

		errCh := make(chan error, errBuffer)

		// table: portfolios
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = q.DeletePortfolio(ctx, arg.PortfolioID)
			errCh <- errors.New("error DeletePortfolio " + err.Error())
		}()

		// table: assets
		if len(assets) > 0 {
			for _, asset := range assets {
				wg.Add(1)
				go func() {
					defer wg.Done()
					argAssest := DeleteAssetParams{
						PortfolioID: asset.PortfolioID,
						TickerID:    int32(asset.TickerID),
					}

					err := q.DeleteAsset(ctx, argAssest)
					errCh <- errors.New("error DeleteAsset " + err.Error())
				}()
			}
		}

		// table: p_categories
		if len(categories) > 0 {
			for _, category := range categories {
				wg.Add(1)
				go func() {
					defer wg.Done()
					argPCategory := DeletePCategoryParams{
						PortfolioID: category.PortfolioID,
						CategoryID: pgtype.Text{
							String: category.CategoryID.String,
							Valid:  true,
						},
					}

					err = q.DeletePCategory(ctx, argPCategory)
					errCh <- errors.New("error DeletePCategory " + err.Error())
				}()
			}
		}

		// table: p_branches
		if len(branches) > 0 {
			for _, branch := range branches {
				wg.Add(1)
				go func() {
					defer wg.Done()
					argPBranch := DeletePBranchParams{
						PortfolioID: branch.PortfolioID,
						BranchID: pgtype.Text{
							String: branch.BranchID.String,
							Valid:  true,
						},
					}
					err = q.DeletePBranch(ctx, argPBranch)
					errCh <- errors.New("error DeletePBranch " + err.Error())
				}()
			}
		}

		// table: p_advisors
		if len(advisors) > 0 {
			for _, advisor := range advisors {
				wg.Add(1)
				go func() {
					defer wg.Done()
					argPAdvisor := DeletePAdvisorParams{
						PortfolioID: advisor.PortfolioID,
						AdvisorID: pgtype.Text{
							String: advisor.AdvisorID.String,
							Valid:  true,
						},
					}
					err = q.DeletePAdvisor(ctx, argPAdvisor)
					errCh <- errors.New("error DeletePAdvisor " + err.Error())
				}()
			}
		}

		// table: p_organizations
		if len(organizations) > 0 {
			for _, organization := range organizations {
				wg.Add(1)
				go func() {
					defer wg.Done()
					argPOrganization := DeletePOrganizationParams{
						PortfolioID: organization.PortfolioID,
						OrganizationID: pgtype.Text{
							String: organization.OrganizationID.String,
							Valid:  true,
						},
					}
					err = q.DeletePOrganization(ctx, argPOrganization)
					errCh <- errors.New("error DeletePOrganization " + err.Error())
				}()
			}
		}

		// Wait for all goroutines to finish
		wg.Wait()
		close(errCh)

		// Collect and handle errors
		for err := range errCh {
			if err != nil {
				return err
			}
		}

		return nil
	})

	result.PortfolioID = arg.PortfolioID
	return result, err
}
