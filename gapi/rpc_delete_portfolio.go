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

func (s *Server) DeletePortfolioProfile(ctx context.Context, in *rd_portfolio_rpc.DeletePortfolioProfileRequest) (*rd_portfolio_rpc.DeletePortfolioProfileResponse, error) {
	// table: portfolios
	err := s.store.DeletePortfolio(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("cannot DeletePortfolio: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to DeletePortfolio: %s", err)
	}

	// table: assets
	assets, err := s.store.GetAssetsByPortfolioId(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("cannot GetAssetsByPortfolioId: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to GetAssetsByPortfolioId: %s", err)
	}
	for _, asset := range assets {
		argAssest := db.DeleteAssetParams{
			PortfolioID: asset.PortfolioID,
			TickerID:    int32(asset.TickerID),
		}

		err := s.store.DeleteAsset(ctx, argAssest)
		if err != nil {
			s.logger.Sugar().Infof("cannot CreateAsset: %v\n", err)
			return nil, status.Errorf(codes.Internal, "failed to CreateAsset: %s", err)
		}
	}

	// table: p_categories
	categories, err := s.store.GetPCategoryByPortfolioId(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("cannot GetPCategoryByPortfolioId: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to GetPCategoryByPortfolioId: %s", err)
	}

	for _, category := range categories {
		argPCategory := db.DeletePCategoryParams{
			PortfolioID: category.PortfolioID,
			CategoryID: pgtype.Text{
				String: category.CategoryID.String,
				Valid:  true,
			},
		}
		err = s.store.DeletePCategory(ctx, argPCategory)
		if err != nil {
			s.logger.Sugar().Infof("cannot DeletePCategory: %v\n", err)
			return nil, status.Errorf(codes.Internal, "failed to DeletePCategory: %s", err)
		}
	}

	// table: p_branches
	branches, err := s.store.GetPBranchByPortfolioId(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("cannot GetPBranchByPortfolioId: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to GetPBranchByPortfolioId: %s", err)
	}

	for _, branch := range branches {
		argPBranch := db.DeletePBranchParams{
			PortfolioID: branch.PortfolioID,
			BranchID: pgtype.Text{
				String: branch.BranchID.String,
				Valid:  true,
			},
		}
		err = s.store.DeletePBranch(ctx, argPBranch)
		if err != nil {
			s.logger.Sugar().Infof("cannot DeletePBranch: %v\n", err)
			return nil, status.Errorf(codes.Internal, "failed to DeletePBranch: %s", err)
		}
	}

	// table: p_advisors
	advisors, err := s.store.GetPAdvisorByPortfolioId(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("cannot GetPAdvisorByPortfolioId: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to GetPAdvisorByPortfolioId: %s", err)
	}

	for _, advisor := range advisors {
		argPAdvisor := db.DeletePAdvisorParams{
			PortfolioID: advisor.PortfolioID,
			AdvisorID: pgtype.Text{
				String: advisor.AdvisorID.String,
				Valid:  true,
			},
		}
		err = s.store.DeletePAdvisor(ctx, argPAdvisor)
		if err != nil {
			s.logger.Sugar().Infof("cannot DeletePAdvisor: %v\n", err)
			return nil, status.Errorf(codes.Internal, "failed to DeletePAdvisor: %s", err)
		}
	}

	// table: p_organizations
	organizations, err := s.store.GetPOrganizationByPortfolioId(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("cannot GetPOrganizationByPortfolioId: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to GetPOrganizationByPortfolioId: %s", err)
	}

	for _, organization := range organizations {
		argPOrganization := db.DeletePOrganizationParams{
			PortfolioID: organization.PortfolioID,
			OrganizationID: pgtype.Text{
				String: organization.OrganizationID.String,
				Valid:  true,
			},
		}
		err = s.store.DeletePOrganization(ctx, argPOrganization)
		if err != nil {
			s.logger.Sugar().Infof("cannot DeletePOrganization: %v\n", err)
			return nil, status.Errorf(codes.Internal, "failed to DeletePOrganization: %s", err)
		}
	}

	fmt.Printf("==> Deleted portfolioId: %s", in.ProfileId)
	return &rd_portfolio_rpc.DeletePortfolioProfileResponse{
		Status: true,
	}, nil
}
