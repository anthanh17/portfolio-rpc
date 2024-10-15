package gapi

import (
	cache "portfolio-profile-rpc/caching"
	db "portfolio-profile-rpc/db/sqlc"
	rd_portfolio_rpc "portfolio-profile-rpc/rd_portfolio_profile_rpc"

	"portfolio-profile-rpc/util"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// Server serves gRPC requests for our service.
type Server struct {
	rd_portfolio_rpc.UnimplementedRdPortfolioRpcServer
	config              util.Config
	store               db.Store
	logger              *zap.Logger
	mongoClient         *mongo.Client
	hashtagCache        *cache.HashtagCache
	secretKeyEncryption string
}

// NewServer creates a new gRPC server.
func NewServer(
	config util.Config,
	store db.Store,
	mongoClient *mongo.Client,
	cachier cache.Cachier,
	logger *zap.Logger,
	secretKeyEncryption string) (*Server, error) {
	server := &Server{
		config:              config,
		store:               store,
		logger:              logger,
		mongoClient:         mongoClient,
		hashtagCache:        cache.NewHashtagCache(cachier, logger),
		secretKeyEncryption: secretKeyEncryption,
	}

	return server, nil
}
