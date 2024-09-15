package gapi

import (
	"context"
	"fmt"
	db "portfolio-profile-rpc/db/sqlc"
	"portfolio-profile-rpc/rd_portfolio_rpc"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdatePortfolioProfile(ctx context.Context, in *rd_portfolio_rpc.UpdatePortfolioProfileRequest) (*rd_portfolio_rpc.UpdatePortfolioProfileResponse, error) {
	// table: portfolios
	argPortfolio := db.UpdatePortfolioParams{
		ID:      in.ProfileId,
		Name:    in.Name,
		Privacy: db.PortfolioPrivacy(in.Privacy),
	}

	_, err := s.store.UpdatePortfolio(ctx, argPortfolio)
	if err != nil {
		s.logger.Sugar().Infof("cannot UpdatePortfolio: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to UpdatePortfolio: %s", err)
	}

	// table: assets
	for _, asset := range in.Assets {
		argAssest := db.UpdateAssetParams{
			PortfolioID: in.ProfileId,
			TickerID:    int32(asset.TickerId),
			Price:       asset.Price,
			Allocation:  asset.Allocation,
		}

		_, err := s.store.UpdateAsset(ctx, argAssest)
		if err != nil {
			s.logger.Sugar().Infof("cannot UpdateAsset: %v\n", err)
			return nil, status.Errorf(codes.Internal, "failed to UpdateAsset: %s", err)
		}
	}

	// table: p_categories
	argPCategory := db.UpdatePCategoryParams{
		PortfolioID: in.ProfileId,
		CategoryID: pgtype.Text{
			String: in.CategoryId,
			Valid:  true,
		},
	}
	_, err = s.store.UpdatePCategory(ctx, argPCategory)
	if err != nil {
		s.logger.Sugar().Infof("cannot UpdatePCategory: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to UpdatePCategory: %s", err)
	}

	// table: p_branches
	for _, branch := range in.BranchId {
		argPBranch := db.UpdatePBranchParams{
			PortfolioID: in.ProfileId,
			BranchID: pgtype.Text{
				String: branch,
				Valid:  true,
			},
		}
		_, err = s.store.UpdatePBranch(ctx, argPBranch)
		if err != nil {
			s.logger.Sugar().Infof("cannot UpdatePBranch: %v\n", err)
			return nil, status.Errorf(codes.Internal, "failed to UpdatePBranch: %s", err)
		}
	}

	// table: p_advisors
	for _, advisor := range in.AdvisorId {
		argPAdvisor := db.UpdatePAdvisorParams{
			PortfolioID: in.ProfileId,
			AdvisorID: pgtype.Text{
				String: advisor,
				Valid:  true,
			},
		}
		_, err = s.store.UpdatePAdvisor(ctx, argPAdvisor)
		if err != nil {
			s.logger.Sugar().Infof("cannot UpdatePAdvisor: %v\n", err)
			return nil, status.Errorf(codes.Internal, "failed to UpdatePAdvisor: %s", err)
		}
	}

	// table: p_organizations
	for _, organization := range in.AdvisorId {
		argPOrganization := db.UpdatePOrganizationParams{
			PortfolioID: in.ProfileId,
			OrganizationID: pgtype.Text{
				String: organization,
				Valid:  true,
			},
		}
		_, err = s.store.UpdatePOrganization(ctx, argPOrganization)
		if err != nil {
			s.logger.Sugar().Infof("cannot UpdatePOrganization: %v\n", err)
			return nil, status.Errorf(codes.Internal, "failed to UpdatePOrganization: %s", err)
		}
	}

	fmt.Printf("\n==> Updated portfolioId: %s", in.ProfileId)
	return &rd_portfolio_rpc.UpdatePortfolioProfileResponse{
		Status: true,
	}, nil
}
