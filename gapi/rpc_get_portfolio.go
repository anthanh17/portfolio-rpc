package gapi

import (
	"context"
	"portfolio-profile-rpc/rd_portfolio_rpc"
)

func (s *Server) GetProfileByUserID(ctx context.Context, in *rd_portfolio_rpc.GetProfileByUserIDRequest) (*rd_portfolio_rpc.GetProfileByUserIDResponse, error) {

	// fmt.Printf("\n==> Created portfolioId: %s", txResult.PortfolioID)
	return &rd_portfolio_rpc.GetProfileByUserIDResponse{}, nil
}
