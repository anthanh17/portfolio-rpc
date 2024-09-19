package db

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// CREATE
type CreatePortfolioCategoryTxParams struct {
	Name       string   `json:"name"`
	ProfileIds []string `json:"profile_ids"`
	UserID     string   `json:"user_id"`
}

type CreatePortfolioCategoryTxResult struct {
	CategoryID string `json:"category_id"`
}

func (store *SQLStore) CreatePortfolioCategoryTx(ctx context.Context, arg CreatePortfolioCategoryTxParams) (CreatePortfolioCategoryTxResult, error) {
	var result CreatePortfolioCategoryTxResult
	categoryID := uuid.New().String()

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		argCategory := CreatePortfolioCategoryParams{
			ID:   categoryID,
			Name: arg.Name,
		}

		// Create a new category
		_, err = q.CreatePortfolioCategory(ctx, argCategory)
		if err != nil {
			return errors.New("error CreatePortfolioCategory " + err.Error())
		}

		// table: u_categories
		argUCategory := CreateUCategoryParams{
			CategoryID: pgtype.Text{
				String: categoryID,
				Valid:  true,
			},
			UserID: arg.UserID,
		}
		_, err = q.CreateUCategory(ctx, argUCategory)
		if err != nil {
			return errors.New("error CreatePortfolioCategory - CreateUCategory: " + err.Error())
		}

		// table: p_categories
		if len(arg.ProfileIds) > 0 {
			for _, profileID := range arg.ProfileIds {
				argPCategory := CreatePCategoryParams{
					PortfolioID: profileID,
					CategoryID: pgtype.Text{
						String: categoryID,
						Valid:  true,
					},
				}
				_, err = q.CreatePCategory(ctx, argPCategory)
				if err != nil {
					return errors.New("error CreatePCategory " + err.Error())
				}
			}
		}

		return err
	})

	result.CategoryID = categoryID
	return result, err
}

// UPDATE
type UpdatePortfolioCategoryTxParams struct {
	CategoryID string   `json:"category_id"`
	Name       string   `json:"name"`
	ProfileIds []string `json:"profile_ids"`
}

type UpdatePortfolioCategoryTxResult struct {
	CategoryID string `json:"category_id"`
}

func (store *SQLStore) UpdatePortfolioCategoryTx(ctx context.Context, arg UpdatePortfolioCategoryTxParams) (UpdatePortfolioCategoryTxResult, error) {
	var result UpdatePortfolioCategoryTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// table: portfolio_categories
		argPortfolioCategory := UpdatePortfolioCategoryParams{
			ID:   arg.CategoryID,
			Name: arg.Name,
		}

		_, err = q.UpdatePortfolioCategory(ctx, argPortfolioCategory)
		if err != nil {
			return errors.New("error UpdatePortfolioCategory " + err.Error())
		}

		// table: p_categories
		if len(arg.ProfileIds) > 0 {
			pCategories, err := q.GetPCategoryByCategoryId(ctx, pgtype.Text{
				String: arg.CategoryID,
				Valid:  true,
			})
			if err != nil {
				return errors.New("error GetPCategoryByCategoryId " + err.Error())
			}

			currentPCategories := make([]string, len(pCategories))
			// convert branches to currentBranches
			for i, pc := range pCategories {
				currentPCategories[i] = pc.PortfolioID
			}

			add, remove := findDifferences(currentPCategories, arg.ProfileIds)

			if len(add) > 0 {
				for _, portfolioID := range add {
					argPC := CreatePCategoryParams{
						PortfolioID: portfolioID,
						CategoryID: pgtype.Text{
							String: arg.CategoryID,
							Valid:  true,
						},
					}
					_, err = q.CreatePCategory(ctx, argPC)
					if err != nil {
						return errors.New("error UpdatePortfolioCategory - CreatePCategory " + err.Error())
					}
				}
			}

			if len(remove) > 0 {
				for _, portfolioID := range add {
					argPC := DeletePCategoryParams{
						PortfolioID: portfolioID,
						CategoryID: pgtype.Text{
							String: arg.CategoryID,
							Valid:  true,
						},
					}
					err = q.DeletePCategory(ctx, argPC)
					if err != nil {
						return errors.New("error UpdatePortfolioCategory - DeletePCategory " + err.Error())
					}
				}
			}
		}

		return err
	})

	result.CategoryID = arg.CategoryID
	return result, err
}

// DELETE
type DeletePortfolioCategoryTxParams struct {
	CategoryID string `json:"category_id"`
}

type DeletePortfolioCategoryTxResult struct {
	Status bool `json:"status"`
}

func (store *SQLStore) DeletePortfolioCategoryTx(ctx context.Context, arg DeletePortfolioCategoryTxParams) (DeletePortfolioCategoryTxResult, error) {
	var result DeletePortfolioCategoryTxResult
	result.Status = false

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// table: portfolio_categories
		err = q.DeletePortfolioCategory(ctx, arg.CategoryID)
		if err != nil {
			return errors.New("error DeletePortfolioCategory " + err.Error())
		}

		// table: p_categories
		categories, err := q.GetPCategoryByCategoryId(ctx, pgtype.Text{
			String: arg.CategoryID,
			Valid:  true,
		})
		if err != nil {
			return errors.New("error GetPCategoryByCategoryId " + err.Error())
		}

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

		return err
	})

	if err != nil {
		result.Status = true
	}
	return result, err
}

// Remove portfolio profile in category api
type RemovePortfolioProfileInCategoryTxParams struct {
	CategoryID   string   `json:"category_id"`
	PortfolioIDs []string `json:"portfolio_ids"`
}

type RemovePortfolioProfileInCategoryTxResult struct {
	Status bool `json:"status"`
}

func (store *SQLStore) RemovePortfolioProfileInCategoryTx(ctx context.Context, arg RemovePortfolioProfileInCategoryTxParams) (RemovePortfolioProfileInCategoryTxResult, error) {
	var result RemovePortfolioProfileInCategoryTxResult
	result.Status = false

	err := store.execTx(ctx, func(q *Queries) error {
		errCh := make(chan error, len(arg.PortfolioIDs))

		// Start goroutines for each portfolio ID
		var wg sync.WaitGroup
		for _, portfolioID := range arg.PortfolioIDs {
			wg.Add(1)
			go func() {
				defer wg.Done()

				argPCategory := DeletePCategoryParams{
					PortfolioID: portfolioID,
					CategoryID: pgtype.Text{
						String: arg.CategoryID,
						Valid:  true,
					},
				}

				err := q.DeletePCategory(ctx, argPCategory)
				errCh <- err
			}()
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

	if err != nil {
		result.Status = true
	}
	return result, err
}
