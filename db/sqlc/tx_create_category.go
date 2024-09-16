package db

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreatePortfolioCategoryTxParams struct {
	Name       string   `json:"name"`
	ProfileIds []string `json:"profile_ids"`
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

		if len(arg.ProfileIds) > 0 {
			for _, profileID := range arg.ProfileIds {
				// table: p_categories
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
