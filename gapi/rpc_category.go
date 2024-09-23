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

func (s *Server) GetCategoryByUserID(ctx context.Context, in *rd_portfolio_rpc.GetCategoryByUserIDRequest) (*rd_portfolio_rpc.GetCategoryByUserIDResponse, error) {
	// TODO: check lai, err add goroutine
	// from table: u_catagories -> list categories
	arg := db.GetUCategoryByUserIdParams{
		UserID: in.UserId,
		Limit:  int32(in.Size),
		Offset: int32(in.Page),
	}

	uCategories, err := s.store.GetUCategoryByUserId(ctx, arg)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot GetUCategoryByUserId: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to get user-category: %s", err)
	}

	// -------------   Start goroutines --------
	// Get total categories by user id
	errCategoriesCh := make(chan error)
	totalCategoriesCh := make(chan int64)
	go func() {
		total, err := s.store.CountCategoriesByUserID(ctx, in.UserId)
		errCategoriesCh <- err

		totalCategoriesCh <- total
		close(totalCategoriesCh)
		close(errCategoriesCh)
	}()

	// from list categories -> table: p_categories -> list profile
	dataCh := make(chan *rd_portfolio_rpc.CategoryData)
	go func() {
		for _, value := range uCategories {
			categoryInfo, _ := s.store.GetCategoryInfo(ctx, value.CategoryID.String)

			count, _ := s.store.CountProfilesInCategory(ctx, pgtype.Text{
				String: value.CategoryID.String,
				Valid:  true,
			})

			// errCh <- err
			dataCh <- &rd_portfolio_rpc.CategoryData{
				Id:            value.CategoryID.String,
				Name:          categoryInfo.Name,
				NumberProfile: uint64(count),
				CreatedAt:     uint64(categoryInfo.CreatedAt.Unix()),
				UpdatedAt:     uint64(categoryInfo.UpdatedAt.Unix()),
			}
		}
		close(dataCh)
	}()

	var data []*rd_portfolio_rpc.CategoryData
	for d := range dataCh {
		data = append(data, d)
	}

	if <-errCategoriesCh != nil {
		s.logger.Sugar().Infof("\ncannot CountCategoriesByUserID: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to CountCategoriesByUserID: %s", err)
	}

	total := <-totalCategoriesCh
	totalPage := int(math.Ceil(float64(total) / float64(in.Size)))

	fmt.Printf("\n==> Get list category by user id: %s", in.UserId)
	return &rd_portfolio_rpc.GetCategoryByUserIDResponse{
		Data:        data,
		Total:       uint64(total),
		CurrentPage: uint64(in.Page),
		TotalPage:   uint64(totalPage),
	}, nil
}

func (s *Server) GetDetailCategogy(ctx context.Context, in *rd_portfolio_rpc.GetDetailCategogyRequest) (*rd_portfolio_rpc.GetDetailCategogyResponse, error) {
	// -------------   Start goroutines --------
	// errGet := make(chan error, 4)

	// Get total p_categories by category id
	totalPCategoriesCh := make(chan int64)
	go func() {
		total, _ := s.store.CountPCategoryByCategoryId(ctx, pgtype.Text{
			String: in.CategogyId,
			Valid:  true,
		})

		// errGet <- err
		totalPCategoriesCh <- total
		close(totalPCategoriesCh)
	}()

	// tale: portfolio_categories -> Get category info from CategogyId
	categoryInfoCh := make(chan db.HamonixBusinessPortfolioCategory)
	go func() {
		categoryInfo, _ := s.store.GetCategoryInfo(ctx, in.CategogyId)
		// errGet <- err
		fmt.Println("categoryInfo:", categoryInfo)
		categoryInfoCh <- categoryInfo
		close(categoryInfoCh)
	}()

	// table: p_categories -> get list portfolio_id from CategogyId
	portfolioIDCh := make(chan string)
	portfolioIDs := []string{}
	go func() {
		arg := db.GetPCategoryByCategoryIdPagingParams{
			CategoryID: pgtype.Text{
				String: in.CategogyId,
				Valid:  true,
			},
			Limit:  int32(in.Size),
			Offset: int32(in.Page),
		}
		portfolioIDs, _ := s.store.GetPCategoryByCategoryIdPaging(ctx, arg)
		// errCh <- err

		for _, item := range portfolioIDs {
			portfolioIDCh <- item
		}
		close(portfolioIDCh)
	}()

	for portfolioID := range portfolioIDCh {
		portfolioIDs = append(portfolioIDs, portfolioID)
	}

	categoryInfo := <-categoryInfoCh
	total := <-totalPCategoriesCh

	// table: profile
	var profiles []*rd_portfolio_rpc.TCProfile
	if len(portfolioIDs) > 0 {
		for _, portfolioID := range portfolioIDs {
			profile, err := s.store.GetProfilesByPortfolioId(ctx, portfolioID)
			if err != nil {
				s.logger.Sugar().Infof("\ncannot GetProfilesByPortfolioId: %v\n", err)
				continue
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

	// Calc totalPage
	totalPage := int(math.Ceil(float64(total) / float64(in.Size)))

	fmt.Printf("\n==> Get detail category_id: %s", in.CategogyId)
	return &rd_portfolio_rpc.GetDetailCategogyResponse{
		Id:          categoryInfo.ID,
		Name:        categoryInfo.Name,
		Profiles:    profiles,
		Total:       uint64(total),
		CurrentPage: uint64(in.Page),
		TotalPage:   uint64(totalPage),
	}, nil
}
