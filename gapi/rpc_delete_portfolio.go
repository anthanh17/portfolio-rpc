package gapi

import (
	"context"
	"fmt"
	db "portfolio-profile-rpc/db/sqlc"
	"portfolio-profile-rpc/rd_portfolio_rpc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeletePortfolioProfile(ctx context.Context, in *rd_portfolio_rpc.DeletePortfolioProfileRequest) (*rd_portfolio_rpc.DeletePortfolioProfileResponse, error) {
	// table: portfolios
	portfolioId := uuid.New().String()

	argPortfolio := db.CreatePortfolioParams{
		ID:      portfolioId,
		Name:    in.Name,
		Privacy: db.PortfolioPrivacy(in.Privacy),
	}

	_, err := s.store.CreatePortfolio(ctx, argPortfolio)
	if err != nil {
		s.logger.Sugar().Infof("cannot CreatePortfolio: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to CreatePortfolio: %s", err)
	}

	// table: assets
	for _, asset := range in.Assets {
		argAssest := db.CreateAssetParams{
			PortfolioID: portfolioId,
			TickerID:    int32(asset.TickerId),
			Price:       asset.Price,
			Allocation:  asset.Allocation,
		}

		_, err := s.store.CreateAsset(ctx, argAssest)
		if err != nil {
			s.logger.Sugar().Infof("cannot CreateAsset: %v\n", err)
			return nil, status.Errorf(codes.Internal, "failed to CreateAsset: %s", err)
		}
	}

	// table: p_categories
	argPCategory := db.CreatePCategoryParams{
		PortfolioID: portfolioId,
		CategoryID: pgtype.Text{
			String: in.CategoryId,
			Valid:  true,
		},
	}
	_, err = s.store.CreatePCategory(ctx, argPCategory)
	if err != nil {
		s.logger.Sugar().Infof("cannot CreatePCategory: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to CreatePCategory: %s", err)
	}

	// table: p_branches
	for _, branch := range in.BranchId {
		argPBranch := db.CreatePBranchParams{
			PortfolioID: portfolioId,
			BranchID: pgtype.Text{
				String: branch,
				Valid:  true,
			},
		}
		_, err = s.store.CreatePBranch(ctx, argPBranch)
		if err != nil {
			s.logger.Sugar().Infof("cannot CreatePBranche: %v\n", err)
			return nil, status.Errorf(codes.Internal, "failed to CreatePBranche: %s", err)
		}
	}

	// table: p_advisors
	for _, advisor := range in.AdvisorId {
		argPAdvisor := db.CreatePAdvisorParams{
			PortfolioID: portfolioId,
			AdvisorID: pgtype.Text{
				String: advisor,
				Valid:  true,
			},
		}
		_, err = s.store.CreatePAdvisor(ctx, argPAdvisor)
		if err != nil {
			s.logger.Sugar().Infof("cannot CreatePAdvisor: %v\n", err)
			return nil, status.Errorf(codes.Internal, "failed to CreatePAdvisor: %s", err)
		}
	}

	// table: p_organizations
	for _, organization := range in.AdvisorId {
		argPOrganization := db.CreatePOrganizationParams{
			PortfolioID: portfolioId,
			OrganizationID: pgtype.Text{
				String: organization,
				Valid:  true,
			},
		}
		_, err = s.store.CreatePOrganization(ctx, argPOrganization)
		if err != nil {
			s.logger.Sugar().Infof("cannot CreatePOrganization: %v\n", err)
			return nil, status.Errorf(codes.Internal, "failed to CreatePOrganization: %s", err)
		}
	}

	fmt.Printf("==> Created portfolioId: %s", portfolioId)
	return &rd_portfolio_rpc.DeletePortfolioProfileResponse{
		Status: true,
	}, nil
}
