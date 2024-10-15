package gapi

import (
	"context"
	"errors"
	"fmt"
	"math"
	cache "portfolio-profile-rpc/caching"
	db "portfolio-profile-rpc/db/sqlc"
	rd_portfolio_rpc "portfolio-profile-rpc/rd_portfolio_profile_rpc"
	"portfolio-profile-rpc/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreatePortfolioProfile(ctx context.Context, in *rd_portfolio_rpc.CreatePortfolioProfileRequest) (*rd_portfolio_rpc.CreatePortfolioProfileResponse, error) {
	fmt.Printf("\n---> Create portfolio profile request: { %v }\n", in)

	if in.AuthorId == "" {
		s.logger.Sugar().Infof("\nAuthorId empty\n")
		return nil, status.Errorf(codes.Internal, "failed to create portfolio profile: AuthorId empty")
	}

	// convert assests
	assesConvert := make([]*db.ProfileAsset, len(in.Assets))
	for i, asset := range in.Assets {
		assesConvert[i] = &db.ProfileAsset{
			TickerName: asset.TickerName.GetValue(),
			Allocation: asset.Allocation.GetValue(),
			Price:      asset.Price.GetValue(),
		}
	}

	arg := db.CreatePortfolioProfileTxParams{
		ProfileName:    in.Name,
		Privacy:        in.Privacy,
		AuthorId:       in.AuthorId,
		Advisors:       in.Advisors,
		Branches:       in.Branches,
		Organizations:  in.Organizations,
		Accounts:       in.Accounts,
		ExpectedReturn: in.ExpectedReturn,
		IsNewBuyPoint:  in.IsNewBuyPoint,
		Assets:         assesConvert,
	}

	// add transaction - create a new portfolio profile
	txResult, err := s.store.CreatePortfolioProfileTx(ctx, arg)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot CreatePortfolioProfile - CreatePortfolioProfileTx: %v\n", err)
		if db.ErrorCode(err) == db.UniqueViolation {
			s.logger.Sugar().Infof("\ncannot CreatePortfolioProfile - CreatePortfolioProfileTx: %v\n", err)
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to create portfolio profile: %s", err)
	}

	if len(in.Hashtags) > 0 {
		// 2. Hashtag Profile: profile_hashtags
		_, err = s.hashtagCache.AddProfileIdHashtagsCacheElements(ctx, cache.ProfileIdHashtags{
			ProfileId: txResult.ProfileId,
			Hashtags:  in.Hashtags,
		})
		if err != nil {
			s.logger.Sugar().Infof("\ncannot CreatePortfolioProfile - AddProfileIdHashtagsCacheElements: %v\n", err)
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		for _, hashtag := range in.Hashtags {
			// 1. Add Leaderboard
			_, err := s.hashtagCache.AddHashtagLeaderboardCacheElements(
				ctx,
				[]cache.HashtagLeaderboard{
					{
						Score:  0, // TODO: Calculate score
						Member: hashtag,
					},
				})
			if err != nil {
				s.logger.Sugar().Infof("\ncannot CreatePortfolioProfile - AddHashtagLeaderboardCacheElements: %v\n", err)
				return nil, status.Errorf(codes.Internal, err.Error())
			}

			// 2. Hashtag Profile: hashtag_profiles
			_, err = s.hashtagCache.AddHashtagProfileIdsCacheElements(ctx, cache.HashtagProfileIds{
				Hashtag:    hashtag,
				ProfileIds: []string{txResult.ProfileId},
			})
			if err != nil {
				s.logger.Sugar().Infof("\ncannot CreatePortfolioProfile - AddHashtagProfileCacheElements: %v\n", err)
				return nil, status.Errorf(codes.Internal, err.Error())
			}
		}
	}

	fmt.Printf("\n==> Created profileId: %s\n", txResult.ProfileId)
	return &rd_portfolio_rpc.CreatePortfolioProfileResponse{
		ProfileId: txResult.ProfileId,
	}, nil
}

func (s *Server) UpdatePortfolioProfile(ctx context.Context, in *rd_portfolio_rpc.UpdatePortfolioProfileRequest) (*rd_portfolio_rpc.UpdatePortfolioProfileResponse, error) {
	fmt.Printf("\n---> Update portfolio profile request: { %v }\n", in)

	// ProfileId
	arg := db.NewUpdatePortfolioProfileTxParamsBuilder(in.ProfileId)

	// ProfileName
	if in.Name != nil {
		arg.WithProfileName(in.Name)
	}

	// Advisors
	if in.Advisors != nil {
		arg.WithAdvisors(in.Advisors)
	}

	// Branches
	if in.Branches != nil {
		arg.WithBranches(in.Branches)
	}

	// Organizations
	if in.Organizations != nil {
		arg.WithOrganizations(in.Organizations)
	}

	// Accounts
	if in.Accounts != nil {
		arg.WithAccounts(in.Accounts)
	}

	// ExpectedReturn
	if in.ExpectedReturn != nil {
		arg.WithExpectedReturn(in.ExpectedReturn)
	}

	// IsNewBuyPoint
	if in.IsNewBuyPoint != nil {
		arg.WithIsNewBuyPoint(in.IsNewBuyPoint)
	}

	// Assets
	if in.Assets != nil {
		arg.WithAssets(in.Assets)
	}

	// Privacy
	if in.Privacy != nil {
		arg.WithPrivacy(in.Privacy)
	}

	argParms := *arg.Build()

	// add transaction - update a portfolio profile
	txResult, err := s.store.UpdatePortfolioProfileTx(ctx, argParms)
	if err != nil {
		s.logger.Sugar().Infof("\n cannot UpdatePortfolioProfile - UpdatePortfolioProfileTx: %v\n", err)
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "profileId not exists: %s", err)
		}

		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}

		return nil, status.Errorf(codes.Internal, "failed to update profile: %s", err)
	}

	// Hashtag Profile
	if in.Hashtags != nil {
		hashtags, err := s.hashtagCache.GetProfileIdHashtagsCacheElements(ctx, in.ProfileId)
		if err != nil {
			s.logger.Sugar().Infof("\ncannot UpdatePortfolioProfile - GetProfileIdHashtagsCacheElements: %v\n", err)
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		// delete all hashtag in profile
		ok, err := s.hashtagCache.DeleteProfileIdsHashtagProfileCache(ctx, hashtags, []string{in.ProfileId})
		if err != nil {
			s.logger.Sugar().Infof("\ncannot UpdatePortfolioProfile - DeleteProfileIdsHashtagProfileCache: %v\n", err)
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		if !ok {
			s.logger.Sugar().Infof("\ncannot UpdatePortfolioProfile - delete error\n")
			return nil, status.Errorf(codes.Internal, "failed to delete hashtag profile")
		}

		// Add elements to set hashtag profile
		if in.Hashtags[0].Value != "" {
			// --> Hashtag Profile
			// 1. Hashtag Profile: profile_hashtags
			var cnvHashtags []string
			for _, value := range in.Hashtags {
				cnvHashtags = append(cnvHashtags, value.Value)
			}

			_, err = s.hashtagCache.AddProfileIdHashtagsCacheElements(ctx, cache.ProfileIdHashtags{
				ProfileId: txResult.ProfileId,
				Hashtags:  cnvHashtags,
			})
			if err != nil {
				s.logger.Sugar().Infof("\ncannot CreatePortfolioProfile - AddProfileIdHashtagsCacheElements: %v\n", err)
				return nil, status.Errorf(codes.Internal, err.Error())
			}

			// 2. Hashtag Profile: hashtag_profiles
			for _, hashtag := range in.Hashtags {
				_, err = s.hashtagCache.AddHashtagProfileIdsCacheElements(ctx, cache.HashtagProfileIds{
					Hashtag:    hashtag.Value,
					ProfileIds: []string{txResult.ProfileId},
				})
				if err != nil {
					s.logger.Sugar().Infof("\ncannot UpdatePortfolioProfile - AddHashtagProfileIdsCacheElements: %v\n", err)
					return nil, status.Errorf(codes.Internal, err.Error())
				}
			}
		}
	}

	fmt.Printf("\n==> Updated profileId: %s\n", txResult.ProfileId)
	return &rd_portfolio_rpc.UpdatePortfolioProfileResponse{
		Status: true,
	}, nil
}

func (s *Server) DeletePortfolioProfile(ctx context.Context, in *rd_portfolio_rpc.DeletePortfolioProfileRequest) (*rd_portfolio_rpc.DeletePortfolioProfileResponse, error) {
	fmt.Printf("\n---> Delete portfolio profile request: { %v }\n", in)

	// check profile exits in db
	exits, err := s.store.CheckIdExitsPortfolioProfile(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("\nDeletePortfolioProfile - cannot CheckIdExitsPortfolioProfile: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to CheckIdExitsPortfolioProfile: %s", err)
	}

	if !exits {
		s.logger.Sugar().Infof("\n DeletePortfolioProfile - data not found - profile id: %s\n", err)
		return nil, status.Errorf(codes.NotFound, "data not found - profile id: %s", in.ProfileId)
	}

	arg := db.DeletePortfolioProfileTxParams{
		ProfileId: in.ProfileId,
	}

	// add transaction - delete a profile
	txResult, err := s.store.DeletePortfolioProfileTx(ctx, arg)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot DeletePortfolioProfile - DeletePortfolioProfileTx: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to delete profile: %s", err)
	}

	// delete all hashtag in profile
	hashtags, err := s.hashtagCache.GetProfileIdHashtagsCacheElements(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot DeletePortfolioProfile - GetProfileIdHashtagsCacheElements: %v\n", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		_, err = s.hashtagCache.DeleteProfileIdsHashtagProfileCache(ctx, hashtags, []string{in.ProfileId})
		if err != nil {
			s.logger.Sugar().Infof("\ncannot DeletePortfolioProfile - DeleteProfileIdsHashtagProfileCache: %v\n", err)
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}

	// Delete hashtag profile
	_, err = s.hashtagCache.DeleteProfileHashtagsCacheKey(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot DeletePortfolioProfile - DeleteProfileHashtagsCacheKey: %v\n", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	fmt.Printf("\n==> Deleted profileId: %s\n", txResult.ProfileId)
	return &rd_portfolio_rpc.DeletePortfolioProfileResponse{
		Status: true,
	}, nil
}

func (s *Server) GetProfileByUserID(ctx context.Context, in *rd_portfolio_rpc.GetProfileByUserIDRequest) (*rd_portfolio_rpc.GetProfileByUserIDResponse, error) {
	fmt.Printf("\n---> Get profile by userId request: { %v }\n", in)

	// get total user_profile by user id
	totalUProfileCh := make(chan int64)
	go func() {
		total, _ := s.store.CountProfilesInUserPortfolioProfile(ctx, in.UserId)

		// errGet <- err
		totalUProfileCh <- total
		close(totalUProfileCh)
	}()

	// TODO: hanled erro goroutine
	// tale: u_portfolio -> list profile_id by user_id
	profileIdsCh := make(chan string)
	go func() {
		arg := db.GetListProfileIdByUserIdParams{
			AuthorID: in.UserId,
			Limit:    int32(in.Size),
			Offset:   (int32(in.Page) - 1) * int32(in.Size),
		}
		profileIds, _ := s.store.GetListProfileIdByUserId(ctx, arg)
		// errGet <- err

		for _, item := range profileIds {
			profileIdsCh <- item
		}

		close(profileIdsCh)
	}()

	total := <-totalUProfileCh
	portfolioIds := []string{}
	for portfolioId := range profileIdsCh {
		portfolioIds = append(portfolioIds, portfolioId)
	}

	var data []*rd_portfolio_rpc.TProfile
	for _, id := range portfolioIds {
		profileInfo, err := s.store.GetProfileInfoById(ctx, id)
		if err != nil {
			s.logger.Sugar().Infof("\ncannot GetProfileByUserID - GetProfileInfoById: %v\n", err)
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		// Get redis list hashtags by profile id
		hashtags, err := s.hashtagCache.GetProfileIdHashtagsCacheElements(ctx, id)
		if err != nil {
			s.logger.Sugar().Infof("\ncannot GetProfileByUserID - GetProfileIdHashtagsCacheElements: %v\n", err)
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		// TODO: chart, Total return
		data = append(data, &rd_portfolio_rpc.TProfile{
			Id:             profileInfo.ID,
			Name:           profileInfo.Name,
			Privacy:        profileInfo.Privacy,
			Author:         profileInfo.AuthorID,
			ExpectedReturn: profileInfo.ExpectedReturn,
			IsNewBuyPoint:  profileInfo.IsNewBuyPoint,
			ProfitAndLoss:  util.RandomFloat(-100.0000, 100.0000), // TODO: Calc ProfitAndLoss
			Stars:          uint64(util.RandomInt(0, 100)),
			Follows:        uint64(util.RandomInt(0, 100)),
			Copies:         uint64(util.RandomInt(0, 100)),
			Hashtags:       hashtags,
			CreatedAt:      uint64(profileInfo.CreatedAt.Unix()),
			UpdatedAt:      uint64(profileInfo.UpdatedAt.Unix()),
		})
	}

	// calc totalPage
	totalPage := int(math.Ceil(float64(total) / float64(in.Size)))

	fmt.Printf("\n==> Get profile by userId: %s\n", in.UserId)
	return &rd_portfolio_rpc.GetProfileByUserIDResponse{
		Data:        data,
		Total:       uint64(total),
		CurrentPage: uint64(in.Page),
		TotalPages:  uint64(totalPage),
	}, nil
}

func (s *Server) GetDetailProfile(ctx context.Context, in *rd_portfolio_rpc.GetDetailProfileRequest) (*rd_portfolio_rpc.GetDetailProfileResponse, error) {
	fmt.Printf("\n---> Get detail profile request: { %v }\n", in)

	profileInfo, err := s.store.GetProfileInfoById(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot GetProfileInfoById: %v\n", err)
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "profile id not exists: %s", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to delete profile: %s", err)
	}

	// Get redis list hashtags by profile id
	hashtags, err := s.hashtagCache.GetProfileIdHashtagsCacheElements(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot GetProfileByUserID - GetProfileIdHashtagsCacheElements: %v\n", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := &rd_portfolio_rpc.GetDetailProfileResponse{
		Id:      profileInfo.ID,
		Name:    profileInfo.Name,
		Privacy: profileInfo.Privacy,
		Author: &rd_portfolio_rpc.ObjInfo{
			Id:   profileInfo.AuthorID,
			Name: "todo future",
		},
		ProfitAndLoss:  util.RandomFloat(-100.0000, 100.0000), // TODO: Calc ProfitAndLoss
		ExpectedReturn: profileInfo.ExpectedReturn,
		IsNewBuyPoint:  profileInfo.IsNewBuyPoint,
		Hashtags:       hashtags,
		Stars:          uint64(util.RandomInt(0, 100)),
		Follows:        uint64(util.RandomInt(0, 100)),
		Copies:         uint64(util.RandomInt(0, 100)),
		CreatedAt:      uint64(profileInfo.CreatedAt.Unix()),
		UpdatedAt:      uint64(profileInfo.UpdatedAt.Unix()),
	}

	// table: hrn_profile_account - relationship profile linked to account
	totalAccountsCh := make(chan int64)
	go func() {
		total, _ := s.store.CountAccountsLinkedProfileByProfileId(ctx, in.ProfileId)

		// errGet <- err
		totalAccountsCh <- total
		close(totalAccountsCh)
	}()

	numberLinkedAccounts := <-totalAccountsCh
	res.NumberLinkedAccounts = uint64(numberLinkedAccounts)

	// table: assets
	assets, err := s.store.GetListAssetsByPortfolioId(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot GetListAssetsByPortfolioId: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to GetListAssetsByPortfolioId: %s", err)
	}

	for _, value := range assets {
		res.Assets = append(res.Assets, &rd_portfolio_rpc.AssetInfo{
			TickerName:  value.TickerName,
			Description: "todo description future",
			Allocation:  value.Allocation,
		})
	}

	// get advisors | branches | organizations
	advisorsBranchesOrganizationsAccounts, err := s.store.GetListAdvisorsBranchesOrganizationsByProfileId(ctx, in.ProfileId)
	if err != nil {
		s.logger.Sugar().Infof("\ncannot advisorsBranchesOrganizations: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to advisorsBranchesOrganizations: %s", err)
	}

	// Decrypted
	if profileInfo.Privacy == db.PrivacyPrivate || profileInfo.Privacy == db.PrivacyProtected {
		for _, item := range advisorsBranchesOrganizationsAccounts {
			// Advisors
			for _, value := range item.Advisors {
				advisorDecrypted, err := util.DecryptData(value, []byte(s.secretKeyEncryption))
				if err != nil {
					s.logger.Sugar().Infof("\nGetDetailProfilet - error advisorDecrypted: %v\n", err)
					return nil, status.Errorf(codes.Internal, "failed to advisorsBranchesOrganizations: %s", err)
				}
				res.Advisors = append(res.Advisors, string(advisorDecrypted))
			}

			// Branches
			for _, value := range item.Branches {
				branchDecrypted, err := util.DecryptData(value, []byte(s.secretKeyEncryption))
				if err != nil {
					s.logger.Sugar().Infof("\nGetDetailProfilet - error branchDecrypted: %v\n", err)
					return nil, status.Errorf(codes.Internal, "failed to advisorsBranchesOrganizations: %s", err)
				}
				res.Branches = append(res.Branches, string(branchDecrypted))
			}

			// Organizations
			for _, value := range item.Organizations {
				organizationDecrypted, err := util.DecryptData(value, []byte(s.secretKeyEncryption))
				if err != nil {
					s.logger.Sugar().Infof("\nGetDetailProfilet - error organizationDecrypted: %v\n", err)
					return nil, status.Errorf(codes.Internal, "failed to advisorsBranchesOrganizations: %s", err)
				}
				res.Organizations = append(res.Organizations, string(organizationDecrypted))
			}

			// Accounts
			for _, value := range item.Accounts {
				accountDecrypted, err := util.DecryptData(value, []byte(s.secretKeyEncryption))
				if err != nil {
					s.logger.Sugar().Infof("\nGetDetailProfilet - error accountDecrypted: %v\n", err)
					return nil, status.Errorf(codes.Internal, "failed to advisorsBranchesOrganizations: %s", err)
				}
				res.Accounts = append(res.Accounts, string(accountDecrypted))
			}
		}
	} else {
		for _, item := range advisorsBranchesOrganizationsAccounts {
			res.Advisors = append(res.Advisors, item.Advisors...)
			res.Branches = append(res.Branches, item.Branches...)
			res.Organizations = append(res.Organizations, item.Organizations...)
			res.Accounts = append(res.Accounts, item.Accounts...)
		}
	}

	fmt.Printf("\n==> Get detail profileId: %s\n", in.ProfileId)
	return res, nil
}
