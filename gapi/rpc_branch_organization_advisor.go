package gapi

import (
	"context"
	"fmt"
	"portfolio-profile-rpc/rd_portfolio_rpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetBranchByID(ctx context.Context, in *rd_portfolio_rpc.GetBranchByIDRequest) (*rd_portfolio_rpc.GetBranchByIDResponse, error) {
	result, err := s.store.GetEQBranchByID(ctx, in.Id)
	if err != nil {
		s.logger.Sugar().Infof("\n error GetEQBranchByID: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to GetEQBranchByID: %s", err)
	}

	fmt.Printf("\n==> GetEQBranchByID: %s", in.Id)
	return &rd_portfolio_rpc.GetBranchByIDResponse{
		Id:          result.ID,
		Code:        result.Code,
		Description: result.Description.String,
	}, nil
}

func (s *Server) GetOrganizationByID(ctx context.Context, in *rd_portfolio_rpc.GetOrganizationByIDRequest) (*rd_portfolio_rpc.GetOrganizationByIDResponse, error) {
	result, err := s.store.GetEQOrganizationByID(ctx, in.Id)
	if err != nil {
		s.logger.Sugar().Infof("\n error GetEQOrganizationByID: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to GetEQOrganizationByID: %s", err)
	}

	fmt.Printf("\n==> GetEQOrganizationByID: %s", in.Id)
	return &rd_portfolio_rpc.GetOrganizationByIDResponse{
		Id:           result.ID,
		Code:         result.Code,
		BackofficeId: result.BackofficeID.String,
		Description:  result.Description.String,
	}, nil
}

func (s *Server) GetAdvisorByID(ctx context.Context, in *rd_portfolio_rpc.GetAdvisorByIDRequest) (*rd_portfolio_rpc.GetAdvisorByIDResponse, error) {
	result, err := s.store.GetEQAdvisorByID(ctx, in.Id)
	if err != nil {
		s.logger.Sugar().Infof("\n error GetEQAdvisorByID: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to GetEQAdvisorByID: %s", err)
	}

	fmt.Printf("\n==> GetEQAdvisorByID: %s", in.Id)
	return &rd_portfolio_rpc.GetAdvisorByIDResponse{
		Id:          result.ID,
		Code:        result.Code.String,
		Description: result.Description.String,
	}, nil
}
