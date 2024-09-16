package gapi

import (
	"context"
	"fmt"
	db "portfolio-profile-rpc/db/sqlc"
	"portfolio-profile-rpc/rd_portfolio_rpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateCategory(ctx context.Context, in *rd_portfolio_rpc.CreateCategoryRequest) (*rd_portfolio_rpc.CreateCategoryResponse, error) {
	arg := db.CreatePortfolioCategoryTxParams{
		Name:       in.Name,
		ProfileIds: in.ProfileIds,
	}

	// Add transaction - create a new category
	txResult, err := s.store.CreatePortfolioCategoryTx(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			s.logger.Sugar().Infof("\ncannot CreatePortfolioCategoryTx: %v\n", err)
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		s.logger.Sugar().Infof("\ncannot CreatePortfolioCategoryTx: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to create portfolio: %s", err)
	}

	fmt.Printf("\n==> Created categoryID: %s", txResult.CategoryID)
	return &rd_portfolio_rpc.CreateCategoryResponse{
		CategoryId: txResult.CategoryID,
	}, nil
}

func (s *Server) UpdateCategory(ctx context.Context, in *rd_portfolio_rpc.UpdateCategoryRequest) (*rd_portfolio_rpc.UpdateCategoryResponse, error) {
	arg := db.UpdatePortfolioCategoryTxParams{
		CategoryID: in.CategoryId,
		Name: in.Name,
		ProfileIds: in.ProfileIds,
	}

	// Add transaction - update a category
	txResult, err := s.store.UpdatePortfolioCategoryTx(ctx, arg)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot UpdatePortfolioCategoryTx: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to update category: %s", err)
	}

	fmt.Printf("\n==> Updated categoryID: %s", txResult.CategoryID)
	return &rd_portfolio_rpc.UpdateCategoryResponse{}, nil
}

func (s *Server) DeleteCategory(ctx context.Context, in *rd_portfolio_rpc.DeleteCategoryRequest) (*rd_portfolio_rpc.DeleteCategoryResponse, error) {
	arg := db.DeletePortfolioCategoryTxParams{
		CategoryID: in.Id,
	}

	// Add transaction - delete a category
	txResult, err := s.store.DeletePortfolioCategoryTx(ctx, arg)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot DeletePortfolioCategoryTx: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to delete category: %s", err)
	}

	fmt.Printf("\n==> Deleted categoryID: %s", in.Id)
	return &rd_portfolio_rpc.DeleteCategoryResponse{
		Status: txResult.Status,
	}, nil
}
