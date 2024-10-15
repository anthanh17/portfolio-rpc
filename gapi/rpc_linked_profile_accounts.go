package gapi

import (
	"context"
	"fmt"
	"math"
	db "portfolio-profile-rpc/db/sqlc"
	rd_portfolio_rpc "portfolio-profile-rpc/rd_portfolio_profile_rpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) LinkedProfileToListAccounts(ctx context.Context, in *rd_portfolio_rpc.LinkedProfileToListAccountsRequest) (*rd_portfolio_rpc.LinkedProfileToListAccountsResponse, error) {
	fmt.Printf("\n---> Linked profile to list accounts request: { %v }\n", in)

	// check profile exits in db
	exits, err := s.store.CheckIdExitsPortfolioProfile(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot CheckIdExitsPortfolioProfile: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to CheckIdExitsPortfolioProfile: %s", err)
	}

	if !exits {
		s.logger.Sugar().Infof("\ndata not found - profile id: %s\n", err)
		return nil, status.Errorf(codes.NotFound, "data not found - profile id: %s", in.ProfileId)
	}

	// table: hrn_profile_account - relationship profile linked to account
	for _, account := range in.AccountIds {
		arg := db.CreateLinkedProfileToAccountParams{
			ProfileID: in.ProfileId,
			AccountID: account,
		}
		_, err = s.store.CreateLinkedProfileToAccount(ctx, arg)
		if err != nil {
			s.logger.Sugar().Infof("\ncannot CreateLinkedProfileToAccount: %v\n", err)
			return nil, status.Errorf(codes.AlreadyExists, "failed to CreateLinkedProfileToAccount: %s", err)
		}
	}

	fmt.Printf("\n==> Created linked profile to list accounts: %s\n", in.ProfileId)
	return &rd_portfolio_rpc.LinkedProfileToListAccountsResponse{
		Status: true,
	}, nil
}

func (s *Server) GetListLinkedProfileAccounts(ctx context.Context, in *rd_portfolio_rpc.GetListLinkedProfileAccountsRequest) (*rd_portfolio_rpc.GetListLinkedProfileAccountsResponse, error) {
	fmt.Printf("\n---> Get list linked profile accounts request: { %v }\n", in)

	// check profile exits in db
	exits, err := s.store.CheckIdExitsPortfolioProfile(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot CheckIdExitsPortfolioProfile: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to CheckIdExitsPortfolioProfile: %s", err)
	}

	if !exits {
		s.logger.Sugar().Infof("\ndata not found - profile id: %s\n", err)
		return nil, status.Errorf(codes.NotFound, "data not found - profile id: %s", in.ProfileId)
	}

	// table: hrn_profile_account - relationship profile linked to account
	totalAccountsCh := make(chan int64)
	go func() {
		total, _ := s.store.CountAccountsLinkedProfileByProfileId(ctx, in.ProfileId)

		// errGet <- err
		totalAccountsCh <- total
		close(totalAccountsCh)
	}()

	total := <-totalAccountsCh

	account_idCh := make(chan string)
	account_ids := []string{}
	go func() {
		arg := db.GetListAccountsLinkedProfileByProfileIdParams{
			ProfileID: in.ProfileId,
			Limit:     int32(in.Size),
			Offset:    int32(in.Page) - 1,
		}
		account_ids, _ := s.store.GetListAccountsLinkedProfileByProfileId(ctx, arg)
		// errCh <- err

		for _, item := range account_ids {
			account_idCh <- item
		}
		close(account_idCh)
	}()

	for account_id := range account_idCh {
		account_ids = append(account_ids, account_id)
	}

	var data []*rd_portfolio_rpc.LinkedAccounts
	// TODO: Chart: future equix api
	for _, account_id := range account_ids {
		data = append(data, &rd_portfolio_rpc.LinkedAccounts{
			AccountId:   account_id,
			AccountName: "future equix api",
		})
	}

	// calc totalPage
	totalPage := int(math.Ceil(float64(total) / float64(in.Size)))

	fmt.Printf("\n==> Created linked profile to list accounts: %s\n", in.ProfileId)
	return &rd_portfolio_rpc.GetListLinkedProfileAccountsResponse{
		Total:       uint64(total),
		CurrentPage: uint64(in.Page),
		TotalPages:  uint64(totalPage),
		Data:        data,
	}, nil
}

func (s *Server) UnLinkedProfileToListAccounts(ctx context.Context, in *rd_portfolio_rpc.UnLinkedProfileToListAccountsRequest) (*rd_portfolio_rpc.UnLinkedProfileToListAccountsResponse, error) {
	fmt.Printf("\n---> Linked profile to list accounts request: { %v }\n", in)

	// check profile exits in db
	exits, err := s.store.CheckIdExitsPortfolioProfile(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot CheckIdExitsPortfolioProfile: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to CheckIdExitsPortfolioProfile: %s", err)
	}

	if !exits {
		s.logger.Sugar().Infof("\ndata not found - profile id: %s\n", err)
		return nil, status.Errorf(codes.NotFound, "data not found - profile id: %s", in.ProfileId)
	}

	// TODO: checl profile not have linked with account

	// table: hrn_profile_account - relationship profile linked to account
	// remove records related profile linked to account
	for _, account := range in.AccountIds {
		arg := db.DeleteLinkedProfileAccountParams{
			ProfileID: in.ProfileId,
			AccountID: account,
		}
		err = s.store.DeleteLinkedProfileAccount(ctx, arg)
		if err != nil {
			s.logger.Sugar().Infof("\ncannot DeleteLinkedProfileAccount: %v\n", err)
			return nil, status.Errorf(codes.Internal, "failed to DeleteLinkedProfileAccount: %s", err)
		}
	}

	fmt.Printf("\n==> UnLinked profile %s to list accounts\n", in.ProfileId)
	return &rd_portfolio_rpc.UnLinkedProfileToListAccountsResponse{
		Status: true,
	}, nil
}
