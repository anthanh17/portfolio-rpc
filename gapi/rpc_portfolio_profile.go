package gapi

import (
	"context"
	"fmt"
	"math"
	db "portfolio-profile-rpc/db/sqlc"
	"portfolio-profile-rpc/rd_portfolio_rpc"

	"slices"

	"github.com/jackc/pgx/v5/pgtype"
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

	// TODO: Hanled erro goroutine
	// tale: u_portfolio -> list portfolio_id by user_id
	portfolioIdCh := make(chan string)
	go func() {
		argUPortfolio := db.GetUPortfolioByUserIdParams{
			UserID: in.UserId,
			Limit:  int32(in.Size),
			Offset: int32(in.Page),
		}
		portfolioIds, _ := s.store.GetUPortfolioByUserId(ctx, argUPortfolio)
		// errGet <- err

		for _, item := range portfolioIds {
			portfolioIdCh <- item.String
		}

		close(portfolioIdCh)
	}()

	// table: p_categories -> get list profile by category
	profileIds := []string{}
	if len(in.CategoryId) > 0 {
		profileCategoriesCh := make(chan string)
		go func() {
			portfolioIds, _ := s.store.GetListProfileIdByCategoryId(ctx, pgtype.Text{
				String: in.CategoryId,
				Valid:  true,
			})
			// errGet <- err

			for _, item := range portfolioIds {
				profileCategoriesCh <- item
			}

			close(profileCategoriesCh)
		}()

		for profileId := range profileCategoriesCh {
			profileIds = append(profileIds, profileId)
		}
	}

	total := <-totalUPortfolioCh
	portfolioIds := []string{}
	for portfolioId := range portfolioIdCh {
		if len(in.CategoryId) > 0 {
			if slices.Contains(profileIds, portfolioId) {
				portfolioIds = append(portfolioIds, portfolioId)
			}
		} else {
			portfolioIds = append(portfolioIds, portfolioId)
		}
	}

	var data []*rd_portfolio_rpc.TProfile
	for _, id := range portfolioIds {
		portfolioInfo, _ := s.store.GetProfilesByPortfolioId(ctx, id)
		// TODO: chart, Total return
		data = append(data, &rd_portfolio_rpc.TProfile{
			Id:        portfolioInfo.ID,
			Name:      portfolioInfo.Name,
			Privacy:   portfolioInfo.Privacy,
			Author:    portfolioInfo.AuthorID,
			CreatedAt: uint64(portfolioInfo.CreatedAt.Unix()),
			UpdatedAt: uint64(portfolioInfo.UpdatedAt.Unix()),
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

func (s *Server) GetDetailProfile(ctx context.Context, in *rd_portfolio_rpc.GetDetailProfileRequest) (*rd_portfolio_rpc.GetDetailProfileResponse, error) {
	portfolioInfo, err := s.store.GetProfilesByPortfolioId(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot DeletePortfolioTx: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to delete portfolio: %s", err)
	}

	res := &rd_portfolio_rpc.GetDetailProfileResponse{
		Id:      portfolioInfo.ID,
		Name:    portfolioInfo.Name,
		Privacy: portfolioInfo.Privacy,
		Author: &rd_portfolio_rpc.ObjInfo{
			Id:   portfolioInfo.AuthorID,
			Name: "todo future",
		},
		// TODO: NumberLinkedAccounts
		CreatedAt: uint64(portfolioInfo.CreatedAt.Unix()),
		UpdatedAt: uint64(portfolioInfo.UpdatedAt.Unix()),
	}

	// table: assets
	assets, err := s.store.GetListAssetsByPortfolioId(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot GetListAssetsByPortfolioId: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to GetListAssetsByPortfolioId: %s", err)
	}

	for _, value := range assets {
		res.Assets = append(res.Assets, &rd_portfolio_rpc.AssetInfo{
			TickerId:    uint64(value.TickerID),
			Name:        "todo name future",
			Description: "todo description future",
			Allocation:  value.Allocation,
		})
	}

	// table: p_categories => list categories by portfolio id
	pCategories, err := s.store.GetListCategoryPCategoryByPortfolioId(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot GetListCategoryPCategoryByPortfolioId: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to GetListCategoryPCategoryByPortfolioId: %s", err)
	}

	for _, value := range pCategories {
		res.Category = append(res.Category, &rd_portfolio_rpc.ObjInfo{
			Id:   value.String,
			Name: "todo future",
		})
	}

	// table: p_advisors => list advisor by portfolio id
	pAdvisors, err := s.store.GetListAdvisorPAdvisorsByPortfolioId(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot GetListAdvisorPAdvisorsByPortfolioId: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to GetListAdvisorPAdvisorsByPortfolioId: %s", err)
	}

	for _, value := range pAdvisors {
		res.Advisor = append(res.Advisor, &rd_portfolio_rpc.ObjInfo{
			Id:   value.String,
			Name: "todo future",
		})
	}

	// table: p_branches => list branch by portfolio id
	pBranches, err := s.store.GetListBranchPBranchByPortfolioId(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot GetListBranchPBranchByPortfolioId: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to GetListBranchPBranchByPortfolioId: %s", err)
	}

	for _, value := range pBranches {
		res.Branch = append(res.Branch, &rd_portfolio_rpc.ObjInfo{
			Id:   value.String,
			Name: "todo future",
		})
	}

	// table: p_organizations => list organizations by portfolio id
	pOrganizations, err := s.store.GetListOrganizationPOrganizationByPortfolioId(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot GetListOrganizationPOrganizationByPortfolioId: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to GetListOrganizationPOrganizationByPortfolioId: %s", err)
	}

	for _, value := range pOrganizations {
		res.Organization = append(res.Organization, &rd_portfolio_rpc.ObjInfo{
			Id:   value.String,
			Name: "todo future",
		})
	}

	fmt.Printf("\n==> Get detail profileId: %s", in.ProfileId)
	return res, nil
}
