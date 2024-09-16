package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
)

type UpdatePortfolioCategoryTxParams struct {
	CategoryID string `json:"category_id"`
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
			ID: arg.CategoryID,
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
