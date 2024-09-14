package gapi

import (
	db "portfolio-profile-rpc/db/sqlc"
	"portfolio-profile-rpc/rd_portfolio_rpc"

	"portfolio-profile-rpc/util"

	"go.uber.org/zap"
)

// Server serves gRPC requests for our service.
type Server struct {
	rd_portfolio_rpc.UnimplementedRdPortfolioRpcServer
	config util.Config
	store  db.Store
	logger *zap.Logger
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store, logger *zap.Logger) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
		logger: logger,
	}

	return server, nil
}
