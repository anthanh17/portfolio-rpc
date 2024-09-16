package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
)

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
			Privacy: PortfolioPrivacy(arg.Privacy),
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
