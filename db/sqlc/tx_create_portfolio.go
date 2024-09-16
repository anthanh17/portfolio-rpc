package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
)

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
			Privacy: PortfolioPrivacy(arg.Privacy),
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
