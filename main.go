package main

import (
	"context"
	"fmt"
	"log"
	"net"
	db "portfolio-profile-rpc/db/sqlc"
	"portfolio-profile-rpc/gapi"
	"portfolio-profile-rpc/rd_portfolio_rpc"

	cache "portfolio-profile-rpc/caching"
	"portfolio-profile-rpc/util"

	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Run server grpc
	if err = runGrpcServer(config); err != nil {
		panic(err)
	}
}

func runGrpcServer(config util.Config) error {
	// Logger
	logger, cleanup, err := util.InitializeLogger(config.Log.Level)
	if err != nil {
		cleanup()
		logger.With(zap.Error(err)).Error("cannot initialize logger")
		return err
	}
	defer cleanup()

	// Database Accessor
	store, cleanupFunc, err := db.InitializeUpDB(config.Database, logger, config.Encrypt.Key)
	if err != nil {
		cleanupFunc()
		logger.Info("error InitializeUpDB")
		return err
	}
	defer cleanupFunc()

	// Mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	bOpts := &options.BSONOptions{
		NilSliceAsEmpty:        true,
		AllowTruncatingDoubles: true,
	}
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Mongo.Url).SetBSONOptions(bOpts))
	if err != nil {
		logger.Sugar().Infof("\ncanot connect mongodb: %v", err)
	}

	// Caching: in case using redis caching
	cacheMaker, err := cache.NewCachierClient(config.Cache, logger)
	if err != nil {
		logger.Info("error NewCachierClient")
		return err
	}

	// gRPC server
	server, err := gapi.NewServer(config, store, mongoClient, cacheMaker, logger, config.Encrypt.Key)
	if err != nil {
		logger.Info("cannot new server")
		return err
	}

	grpcServer := grpc.NewServer()
	rd_portfolio_rpc.RegisterRdPortfolioRpcServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPC.Address)
	if err != nil {
		logger.Info("cannot create listener")
		return err
	}

	fmt.Printf("==> start gRPC at %s ...\n", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		logger.Info("cannot start gRPC server")
		return err
	}

	return nil
}
