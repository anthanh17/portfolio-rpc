package gapi

import (
	"context"
	"fmt"
	db "portfolio-profile-rpc/db/sqlc"
	"portfolio-profile-rpc/rd_portfolio_rpc"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreatePortfolioProfile(ctx context.Context, in *rd_portfolio_rpc.CreatePortfolioProfileRequest) (*rd_portfolio_rpc.CreatePortfolioProfileResponse, error) {
	portfolioId := uuid.New().String()

	// Default authorID = userID
	authorID := in.UserId

	// convert assests
	assesConvert := make([]*db.PortfolioAsset, len(in.Assets))
	for i, asset := range in.Assets {
		assesConvert[i] = &db.PortfolioAsset{
			TickerId:   asset.TickerId,
			Allocation: asset.Allocation,
			Price:      asset.Price,
		}
	}

	// This portfolio belongs to someone else.
	if in.AuthorId != "" {
		authorID = in.AuthorId
	}

	arg := db.CreatePortfolioTxParams{
		PortfolioID:    portfolioId,
		CategoryID:     in.CategoryId,
		PortfolioName:  in.Name,
		OrganizationId: in.OrganizationId,
		BranchId:       in.BranchId,
		AdvisorId:      in.AdvisorId,
		Assets:         assesConvert,
		Privacy:        in.Privacy,
		AuthorID:       authorID,
	}

	// Add transaction - create a new portfolio
	txResult, err := s.store.CreatePortfolioTx(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			s.logger.Sugar().Infof("\ncannot CreatePortfolioTx: %v\n", err)
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		s.logger.Sugar().Infof("\ncannot CreatePortfolioTx: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to create portfolio: %s", err)
	}

	fmt.Printf("\n==> Created portfolioId: %s", txResult.PortfolioID)
	return &rd_portfolio_rpc.CreatePortfolioProfileResponse{
		Status: true,
	}, nil
}

func (s *Server) UpdatePortfolioProfile(ctx context.Context, in *rd_portfolio_rpc.UpdatePortfolioProfileRequest) (*rd_portfolio_rpc.UpdatePortfolioProfileResponse, error) {
	// convert assests
	assesConvert := make([]*db.PortfolioAsset, len(in.Assets))
	for i, asset := range in.Assets {
		assesConvert[i] = &db.PortfolioAsset{
			TickerId:   asset.TickerId,
			Allocation: asset.Allocation,
			Price:      asset.Price,
		}
	}

	arg := db.UpdatePortfolioTxParams{
		PortfolioID:    in.ProfileId,
		CategoryID:     in.CategoryId,
		PortfolioName:  in.Name,
		OrganizationId: in.OrganizationId,
		BranchId:       in.BranchId,
		AdvisorId:      in.AdvisorId,
		Assets:         assesConvert,
		Privacy:        in.Privacy,
	}

	// Add transaction - update a portfolio
	txResult, err := s.store.UpdatePortfolioTx(ctx, arg)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot UpdatePortfolioTx: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to update portfolio: %s", err)
	}

	fmt.Printf("\n==> Updated portfolioId: %s", txResult.PortfolioID)
	return &rd_portfolio_rpc.UpdatePortfolioProfileResponse{
		Status: true,
	}, nil
}

func (s *Server) DeletePortfolioProfile(ctx context.Context, in *rd_portfolio_rpc.DeletePortfolioProfileRequest) (*rd_portfolio_rpc.DeletePortfolioProfileResponse, error) {
	arg := db.DeletePortfolioTxParams{
		PortfolioID: in.ProfileId,
	}

	// Add transaction - delete a portfolio
	txResult, err := s.store.DeletePortfolioTx(ctx, arg)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot DeletePortfolioTx: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to delete portfolio: %s", err)
	}

	fmt.Printf("\n==> Deleted portfolioId: %s", txResult.PortfolioID)
	return &rd_portfolio_rpc.DeletePortfolioProfileResponse{
		Status: true,
	}, nil
}

// TODO
func (s *Server) GetProfileByUserID(ctx context.Context, in *rd_portfolio_rpc.GetProfileByUserIDRequest) (*rd_portfolio_rpc.GetProfileByUserIDResponse, error) {

	// fmt.Printf("\n==> Created portfolioId: %s", txResult.PortfolioID)
	return &rd_portfolio_rpc.GetProfileByUserIDResponse{}, nil
}
