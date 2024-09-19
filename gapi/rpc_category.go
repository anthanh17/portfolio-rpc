package gapi

import (
	"context"
	"fmt"
	"math"
	db "portfolio-profile-rpc/db/sqlc"
	"portfolio-profile-rpc/rd_portfolio_rpc"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func paginate[T any](data []T, page, pageSize int) []T {
	// Check page and pageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// Calc offset
	startIndex := (page - 1) * pageSize
	endIndex := min(startIndex+pageSize, len(data))

	// Return results
	return data[startIndex:endIndex]
}

func (s *Server) CreateCategory(ctx context.Context, in *rd_portfolio_rpc.CreateCategoryRequest) (*rd_portfolio_rpc.CreateCategoryResponse, error) {
	arg := db.CreatePortfolioCategoryTxParams{
		Name:       in.Name,
		ProfileIds: in.ProfileIds,
		UserID:     in.UserId,
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
		Name:       in.Name,
		ProfileIds: in.ProfileIds,
	}

	// Add transaction - update a category
	txResult, err := s.store.UpdatePortfolioCategoryTx(ctx, arg)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot UpdatePortfolioCategoryTx: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to update category: %s", err)
	}

	fmt.Printf("\n==> Updated categoryID: %s", txResult.CategoryID)
	return &rd_portfolio_rpc.UpdateCategoryResponse{
		Name:       in.Name,
		ProfileIds: in.ProfileIds,
	}, nil
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

func (s *Server) GetCategoryByUserID(ctx context.Context, in *rd_portfolio_rpc.GetCategoryByUserIDRequest) (*rd_portfolio_rpc.GetCategoryByUserIDResponse, error) {
	var data []*rd_portfolio_rpc.CategoryData

	// from table: u_catagories -> list categories
	uCategories, err := s.store.GetUCategoryByUserId(ctx, in.UserId)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot GetUCategoryByUserId: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to get user-category: %s", err)
	}

	// from list categories -> table: p_categories -> list profile
	for _, value := range uCategories {
		categoryInfo, err := s.store.GetCategoryInfo(ctx, value.CategoryID.String)
		if err != nil {
			s.logger.Sugar().Infof("\ncannot GetCategoryInfo: %v\n", err)
			return nil, status.Errorf(codes.Internal, "failed to get cateory info: %s", err)
		}

		count, err := s.store.CountProfilesInCategory(ctx, pgtype.Text{
			String: value.CategoryID.String,
			Valid:  true,
		})
		if err != nil {
			s.logger.Sugar().Infof("\ncannot CountProfilesInCategory: %v\n", err)
			return nil, status.Errorf(codes.Internal, "failed to count profile in category: %s", err)
		}

		data = append(data, &rd_portfolio_rpc.CategoryData{
			Id:            value.CategoryID.String,
			Name:          categoryInfo.Name,
			NumberProfile: uint64(count),
			CreatedAt:     uint64(categoryInfo.CreatedAt.Unix()),
			UpdatedAt:     uint64(categoryInfo.UpdatedAt.Unix()),
		})
	}

	pagingResults := paginate(data, int(in.Page), int(in.Size))
	// Calc totalPage
	total := len(data)
	totalPage := int(math.Ceil(float64(total) / float64(in.Size)))

	fmt.Printf("\n==> Get list category by user id: %s", in.UserId)
	return &rd_portfolio_rpc.GetCategoryByUserIDResponse{
		Data:        pagingResults,
		Total:       uint64(total),
		CurrentPage: uint64(in.Page),
		TotalPage:   uint64(totalPage),
	}, nil
}

// Remove portfolio profile in category api
func (s *Server) RemovePortfolioProfileInCategory(ctx context.Context, in *rd_portfolio_rpc.RemovePortfolioProfileInCategoryRequest) (*rd_portfolio_rpc.RemovePortfolioProfileInCategoryResponse, error) {
	arg := db.RemovePortfolioProfileInCategoryTxParams{
		CategoryID:   in.CategogyId,
		PortfolioIDs: in.ProfileIds,
	}

	// Add transaction - Remove portfolio profile in category
	txResult, err := s.store.RemovePortfolioProfileInCategoryTx(ctx, arg)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot RemovePortfolioProfileInCategoryTx: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to remove portfolio profile in category: %s", err)
	}

	fmt.Printf("\n==> Remove portfolio profile in category_id: %s", in.CategogyId)
	return &rd_portfolio_rpc.RemovePortfolioProfileInCategoryResponse{
		Status: txResult.Status,
	}, nil
}

func (s *Server) GetDetailCategogy(ctx context.Context, in *rd_portfolio_rpc.GetDetailCategogyRequest) (*rd_portfolio_rpc.GetDetailCategogyResponse, error) {
	var profiles []*rd_portfolio_rpc.TCProfile
	// tale: portfolio_categories -> Get category info
	categoryInfo, err := s.store.GetCategoryInfo(ctx, in.CategogyId)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot GetCategoryInfo: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to get cateory info: %s", err)
	}

	// table: p_categories -> get list portfolio_id
	pCategories, err := s.store.GetPCategoryByCategoryId(ctx, pgtype.Text{
		String: in.CategogyId,
		Valid:  true,
	})
	if err != nil {
		s.logger.Sugar().Infof("\ncannot GetPCategoryByCategoryId: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to GetPCategoryByCategoryId: %s", err)
	}

	if len(pCategories) > 0 {
		for _, value := range pCategories {
			profile, err := s.store.GetProfilesByPortfolioId(ctx, value.PortfolioID)
			if err != nil {
				s.logger.Sugar().Infof("\ncannot GetProfilesByPortfolioId: %v\n", err)
				return nil, status.Errorf(codes.Internal, "failed to GetProfilesByPortfolioId: %s", err)
			}

			// TODO: Charts, TotalReturn
			profiles = append(profiles, &rd_portfolio_rpc.TCProfile{
				Id:        profile.ID,
				Name:      profile.Name,
				Privacy:   profile.Privacy,
				AuthorId:  profile.AuthorID,
				CreatedAt: uint64(profile.CreatedAt.Unix()),
				UpdatedAt: uint64(profile.UpdatedAt.Unix()),
			})

		}
	}

	pagingResults := paginate(profiles, int(in.Page), int(in.Size))
	// Calc totalPage
	total := len(profiles)
	totalPage := int(math.Ceil(float64(total) / float64(in.Size)))

	fmt.Printf("\n==> Get detail category_id: %s", in.CategogyId)
	return &rd_portfolio_rpc.GetDetailCategogyResponse{
		Id:          categoryInfo.ID,
		Name:        categoryInfo.Name,
		Profiles:    pagingResults,
		Total:       uint64(total),
		CurrentPage: uint64(in.Page),
		TotalPage:   uint64(totalPage),
	}, nil
}
