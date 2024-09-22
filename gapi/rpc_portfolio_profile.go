package gapi

import (
	"context"
	"fmt"
	"math"
	db "portfolio-profile-rpc/db/sqlc"
	"portfolio-profile-rpc/rd_portfolio_rpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreatePortfolioProfile(ctx context.Context, in *rd_portfolio_rpc.CreatePortfolioProfileRequest) (*rd_portfolio_rpc.CreatePortfolioProfileResponse, error) {
	if in.AuthorId == "" {
		s.logger.Sugar().Infof("\nAuthorId empty\n")
		return nil, status.Errorf(codes.Internal, "failed to create portfolio: AuthorId empty")
	}

	// convert assests
	assesConvert := make([]*db.PortfolioAsset, len(in.Assets))
	for i, asset := range in.Assets {
		assesConvert[i] = &db.PortfolioAsset{
			TickerId:   asset.TickerId,
			Allocation: asset.Allocation,
			Price:      asset.Price,
		}
	}

	arg := db.CreatePortfolioTxParams{
		CategoryID:     in.CategoryId,
		PortfolioName:  in.Name,
		OrganizationId: in.OrganizationId,
		BranchId:       in.BranchId,
		AdvisorId:      in.AdvisorId,
		Assets:         assesConvert,
		Privacy:        in.Privacy,
		AuthorID:       in.AuthorId,
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
		ProfileId: txResult.PortfolioID,
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

func (s *Server) GetProfileByUserID(ctx context.Context, in *rd_portfolio_rpc.GetProfileByUserIDRequest) (*rd_portfolio_rpc.GetProfileByUserIDResponse, error) {
	// get total u_portfolio by user id
	totalUPortfolioCh := make(chan int64)
	go func() {
		total, _ := s.store.CountProfilesInUserPortfolio(ctx, in.UserId)

		// errGet <- err
		totalUPortfolioCh <- total
		close(totalUPortfolioCh)
	}()

	// tale: u_portfolio -> list portfolio_id by user_id
	portfolioIdCh := make(chan string)
	go func() {
		argUPortfolio := db.GetUPortfolioByUserIdParams{
			UserID: in.UserId,
			Limit: int32(in.Size),
			Offset: int32(in.Page),
		}
		portfolioIds, _ := s.store.GetUPortfolioByUserId(ctx, argUPortfolio)
		// errGet <- err

		for _, item := range portfolioIds {
			portfolioIdCh <- item.String
		}

		close(portfolioIdCh)
	}()

	total := <-totalUPortfolioCh

	portfolioIds := []string{}
	for portfolioId := range portfolioIdCh {
		portfolioIds = append(portfolioIds, portfolioId)
	}

	var data []*rd_portfolio_rpc.TProfile
	for _, id := range portfolioIds {
		portfolioInfo, _ := s.store.GetProfilesByPortfolioId(ctx, id)
		// TODO: chart, Total return
		data = append(data, &rd_portfolio_rpc.TProfile{
			Id: portfolioInfo.ID,
			Name: portfolioInfo.Name,
			Privacy: portfolioInfo.Privacy,
			Author: portfolioInfo.AuthorID,
			CreatedAt:     uint64(portfolioInfo.CreatedAt.Unix()),
			UpdatedAt:     uint64(portfolioInfo.UpdatedAt.Unix()),
		})
	}

	// Calc totalPage
	totalPage := int(math.Ceil(float64(total) / float64(in.Size)))

	fmt.Printf("\n==> Get Profile By UserID: %s", in.UserId)
	return &rd_portfolio_rpc.GetProfileByUserIDResponse{
		Data:        data,
		Total:       uint64(total),
		CurrentPage: uint64(in.Page),
		TotalPage:   uint64(totalPage),
	}, nil
}
