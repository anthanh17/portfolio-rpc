package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
)

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

	result.Status = true
	return result, err
}
