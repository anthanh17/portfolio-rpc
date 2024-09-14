package main

import (
	"fmt"
	"log"
	"net"
	db "portfolio-profile-rpc/db/sqlc"
	"portfolio-profile-rpc/gapi"
	"portfolio-profile-rpc/rd_portfolio_rpc"

	"portfolio-profile-rpc/util"

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
	store, cleanupFunc, err := db.InitializeUpDB(config.Database, logger)
	if err != nil {
		cleanupFunc()
		logger.Info("error InitializeUpDB")
		return err
	}
	defer cleanupFunc()

	// gRPC server
	server, err := gapi.NewServer(config, store, logger)
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
