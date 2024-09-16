package gapi

import (
	"context"
	"fmt"
	db "portfolio-profile-rpc/db/sqlc"
	"portfolio-profile-rpc/rd_portfolio_rpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
