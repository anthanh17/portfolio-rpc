package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
)

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
