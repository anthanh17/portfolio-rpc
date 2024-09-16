package gapi

import (
	"context"
	"fmt"
	db "portfolio-profile-rpc/db/sqlc"
	"portfolio-profile-rpc/rd_portfolio_rpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
