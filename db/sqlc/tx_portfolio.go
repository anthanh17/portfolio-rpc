package db

import (
	"context"
	"portfolio-profile-rpc/util"
	"time"

	"github.com/google/uuid"
)

const (
	PrivacyPublic    = "PUBLIC"
	PrivacyPrivate   = "PRIVATE"
	PrivacyProtected = "PROTECTED"
)

type ProfileAsset struct {
	TickerName string
	Allocation float64
	Price      float64
}

// CREATE
type CreatePortfolioProfileTxParams struct {
	ProfileName    string
	Privacy        string
	AuthorId       string
	Advisors       []string
	Branches       []string
	Organizations  []string
	Accounts       []string
	ExpectedReturn float64
	IsNewBuyPoint  bool
	Assets         []*ProfileAsset
}

type CreatePortfolioProfileTxResult struct {
	ProfileId string
}

func (store *SQLStore) CreatePortfolioProfileTx(ctx context.Context, arg CreatePortfolioProfileTxParams) (CreatePortfolioProfileTxResult, error) {
	var result CreatePortfolioProfileTxResult
	profileId := uuid.New().String()

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// table: portfolio_profiles
		var argPortfolioProfiles CreatePortfolioProfileParams
		if arg.Privacy == PrivacyPrivate || arg.Privacy == PrivacyProtected {

			// id, err := util.EncryptData([]byte(profileId), []byte(store.secretKeyEncryption))
			// if err != nil {
			// 	store.logger.Sugar().Infof("\n error CreatePortfolioProfileTx - EncryptData - id: %v", err)
			// 	return err
			// }

			name, err := util.EncryptData([]byte(arg.ProfileName), []byte(store.secretKeyEncryption))
			if err != nil {
				store.logger.Sugar().Infof("\n error CreatePortfolioProfileTx - EncryptData - name: %v", err)
				return err
			}

			var advisors []string
			for _, value := range arg.Advisors {
				advisor, err := util.EncryptData([]byte(value), []byte(store.secretKeyEncryption))
				if err != nil {
					store.logger.Sugar().Infof("\n error CreatePortfolioProfileTx - EncryptData - advisors: %v", err)
					return err
				}
				advisors = append(advisors, advisor)
			}

			var branches []string
			for _, value := range arg.Branches {
				branch, err := util.EncryptData([]byte(value), []byte(store.secretKeyEncryption))
				if err != nil {
					store.logger.Sugar().Infof("\n error CreatePortfolioProfileTx - EncryptData - branches: %v", err)
					return err
				}
				branches = append(branches, branch)
			}

			var organizations []string
			for _, value := range arg.Organizations {
				organization, err := util.EncryptData([]byte(value), []byte(store.secretKeyEncryption))
				if err != nil {
					store.logger.Sugar().Infof("\n error CreatePortfolioProfileTx - EncryptData - organizations: %v", err)
					return err
				}
				organizations = append(organizations, organization)
			}

			var accounts []string
			for _, value := range arg.Accounts {
				account, err := util.EncryptData([]byte(value), []byte(store.secretKeyEncryption))
				if err != nil {
					store.logger.Sugar().Infof("\n error CreatePortfolioProfileTx - EncryptData - accounts: %v", err)
					return err
				}
				accounts = append(accounts, account)
			}

			argPortfolioProfiles = CreatePortfolioProfileParams{
				ID:             profileId,
				Name:           name,
				Privacy:        arg.Privacy,
				AuthorID:       arg.AuthorId,
				Advisors:       advisors,
				Branches:       branches,
				Organizations:  organizations,
				Accounts:       accounts,
				ExpectedReturn: arg.ExpectedReturn,
				IsNewBuyPoint:  arg.IsNewBuyPoint,
			}

		} else {
			argPortfolioProfiles = CreatePortfolioProfileParams{
				ID:             profileId,
				Name:           arg.ProfileName,
				Privacy:        arg.Privacy,
				AuthorID:       arg.AuthorId,
				Advisors:       arg.Advisors,
				Branches:       arg.Branches,
				Organizations:  arg.Organizations,
				Accounts:       arg.Accounts,
				ExpectedReturn: arg.ExpectedReturn,
				IsNewBuyPoint:  arg.IsNewBuyPoint,
			}
		}

		_, err = q.CreatePortfolioProfile(ctx, argPortfolioProfiles)
		if err != nil {
			store.logger.Sugar().Infof("\n error CreatePortfolioProfileTx - CreatePortfolioProfile: %v", err)
			return err
		}

		// table: assets
		if len(arg.Assets) > 0 {
			for _, asset := range arg.Assets {
				argAssest := CreateAssetParams{
					PortfolioProfileID: profileId,
					TickerName:         asset.TickerName,
					Price:              asset.Price,
					Allocation:         asset.Allocation,
				}

				_, err := q.CreateAsset(ctx, argAssest)
				if err != nil {
					store.logger.Sugar().Infof("\n error CreatePortfolioProfileTx - CreateAsset: %v", err)
					return err
				}
			}
		}

		return err
	})

	result.ProfileId = profileId
	return result, err
}

// UPDATE
type UpdatePortfolioProfileTxParams struct {
	ProfileId      string
	ProfileName    *string
	Privacy        *string
	Advisors       []*string
	Branches       []*string
	Organizations  []*string
	Accounts       []*string
	ExpectedReturn *float64
	IsNewBuyPoint  *bool
	Assets         []*ProfileAsset
	Hashtags       []*string
}

type UpdatePortfolioProfileTxResult struct {
	ProfileId string
}

func (store *SQLStore) UpdatePortfolioProfileTx(ctx context.Context, arg UpdatePortfolioProfileTxParams) (UpdatePortfolioProfileTxResult, error) {
	var result UpdatePortfolioProfileTxResult

	/*
		- error channel hold err:
			- assets
	*/
	errGet := make(chan error, 2)

	// Assest
	assetIdsCh := make(chan []int64)
	var assetIds []int64

	if arg.Assets != nil {
		go func() {
			// get all ticker
			assetIds, err := store.GetListAssetIdsByPortfolioId(ctx, arg.ProfileId)
			errGet <- err

			// push data to channel
			assetIdsCh <- assetIds
		}()
	}

	// assets
	if arg.Assets != nil {
		assetIds = <-assetIdsCh
	}

	close(errGet)
	// collect and handle errors
	for err := range errGet {
		if err != nil {
			store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx: %v", err)
			return result, err
		}
	}

	// --------------- Start transaction --------------------
	err := store.execTx(ctx, func(q *Queries) error {
		// table: portfolios profile
		privacy, err := q.GetPrivacyProfileById(ctx, arg.ProfileId)
		if err != nil {
			store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - GetPrivacyProfileById: %v", err)
			return err
		}

		// check update name portfolio
		if arg.ProfileName != nil {
			var argPortfolio UpdateNamePortfolioProfileParams
			if privacy == PrivacyPrivate || privacy == PrivacyProtected {
				name, err := util.EncryptData([]byte(*arg.ProfileName), []byte(store.secretKeyEncryption))
				if err != nil {
					store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - EncryptData - name: %v", err)
					return err
				}

				argPortfolio = UpdateNamePortfolioProfileParams{
					ID:        arg.ProfileId,
					Name:      name,
					UpdatedAt: time.Now(),
				}
			} else {
				argPortfolio = UpdateNamePortfolioProfileParams{
					ID:        arg.ProfileId,
					Name:      *arg.ProfileName,
					UpdatedAt: time.Now(),
				}
			}

			_, err := q.UpdateNamePortfolioProfile(ctx, argPortfolio)
			if err != nil {
				store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - UpdateNamePortfolioProfile: %v", err)
				return err
			}
		}

		// check update privacy portfolio
		if arg.Privacy != nil {
			argPortfolio := UpdatePrivacyPortfolioProfileParams{
				ID:        arg.ProfileId,
				Privacy:   *arg.Privacy,
				UpdatedAt: time.Now(),
			}

			_, err := q.UpdatePrivacyPortfolioProfile(ctx, argPortfolio)
			if err != nil {
				store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - UpdatePrivacyPortfolioProfile: %v", err)
				return err
			}
		}

		// update advisors
		if arg.Advisors != nil {
			// case: user pass update empty array
			if *arg.Advisors[0] == "" {
				arg := UpdateAdvisorsPortfolioProfileParams{
					ID:       arg.ProfileId,
					Advisors: []string{},
				}
				_, err := q.UpdateAdvisorsPortfolioProfile(ctx, arg)
				if err != nil {
					store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - UpdateAdvisorsPortfolioProfile: %v", err)
					return err
				}
			} else {
				var advisors []string
				if privacy == PrivacyPrivate || privacy == PrivacyProtected {
					for _, value := range arg.Advisors {
						advisor, err := util.EncryptData([]byte(*value), []byte(store.secretKeyEncryption))
						if err != nil {
							store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - EncryptData - advisors: %v", err)
							return err
						}
						advisors = append(advisors, advisor)
					}
				} else {
					for _, value := range arg.Advisors {
						advisors = append(advisors, *value)
					}
				}

				arg := UpdateAdvisorsPortfolioProfileParams{
					ID:       arg.ProfileId,
					Advisors: advisors,
				}
				_, err := q.UpdateAdvisorsPortfolioProfile(ctx, arg)
				if err != nil {
					store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - UpdateAdvisorsPortfolioProfile: %v", err)
					return err
				}
			}
		}

		// update branches
		if arg.Branches != nil {
			// case: user pass update empty array
			if *arg.Branches[0] == "" {
				arg := UpdateBranchesPortfolioProfileParams{
					ID:       arg.ProfileId,
					Branches: []string{},
				}
				_, err := q.UpdateBranchesPortfolioProfile(ctx, arg)
				if err != nil {
					store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - UpdateBranchesPortfolioProfile: %v", err)
					return err
				}
			} else {
				var branches []string
				if privacy == PrivacyPrivate || privacy == PrivacyProtected {
					for _, value := range arg.Branches {
						branch, err := util.EncryptData([]byte(*value), []byte(store.secretKeyEncryption))
						if err != nil {
							store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - EncryptData - branches: %v", err)
							return err
						}
						branches = append(branches, branch)
					}
				} else {
					for _, value := range arg.Branches {
						branches = append(branches, *value)
					}
				}

				arg := UpdateBranchesPortfolioProfileParams{
					ID:       arg.ProfileId,
					Branches: branches,
				}
				_, err := q.UpdateBranchesPortfolioProfile(ctx, arg)
				if err != nil {
					store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - UpdateBranchesPortfolioProfile: %v", err)
					return err
				}
			}
		}

		// update organizations
		if arg.Organizations != nil {
			// case: user pass update empty array
			if *arg.Organizations[0] == "" {
				arg := UpdateOrganizationsPortfolioProfileParams{
					ID:            arg.ProfileId,
					Organizations: []string{},
				}
				_, err := q.UpdateOrganizationsPortfolioProfile(ctx, arg)
				if err != nil {
					store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - UpdateOrganizationsPortfolioProfile: %v", err)
					return err
				}
			} else {
				var organizations []string
				if privacy == PrivacyPrivate || privacy == PrivacyProtected {
					for _, value := range arg.Organizations {
						organization, err := util.EncryptData([]byte(*value), []byte(store.secretKeyEncryption))
						if err != nil {
							store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - EncryptData - organizations: %v", err)
							return err
						}
						organizations = append(organizations, organization)
					}
				} else {
					for _, value := range arg.Organizations {
						organizations = append(organizations, *value)
					}
				}

				arg := UpdateOrganizationsPortfolioProfileParams{
					ID:            arg.ProfileId,
					Organizations: organizations,
				}
				_, err := q.UpdateOrganizationsPortfolioProfile(ctx, arg)
				if err != nil {
					store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - UpdateOrganizationsPortfolioProfile: %v", err)
					return err
				}
			}
		}

		// update accounts
		if arg.Accounts != nil {
			// case: user pass update empty array
			if *arg.Accounts[0] == "" {
				arg := UpdateAccountsPortfolioProfileParams{
					ID:       arg.ProfileId,
					Accounts: []string{},
				}
				_, err := q.UpdateAccountsPortfolioProfile(ctx, arg)
				if err != nil {
					store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - UpdateAccountsPortfolioProfile: %v", err)
					return err
				}
			} else {
				var accounts []string
				if privacy == PrivacyPrivate || privacy == PrivacyProtected {
					for _, value := range arg.Accounts {
						account, err := util.EncryptData([]byte(*value), []byte(store.secretKeyEncryption))
						if err != nil {
							store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - EncryptData - accounts: %v", err)
							return err
						}
						accounts = append(accounts, account)
					}
				} else {
					for _, value := range arg.Accounts {
						accounts = append(accounts, *value)
					}
				}

				arg := UpdateAccountsPortfolioProfileParams{
					ID:       arg.ProfileId,
					Accounts: accounts,
				}
				_, err := q.UpdateAccountsPortfolioProfile(ctx, arg)
				if err != nil {
					store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - UpdateAccountsPortfolioProfile: %v", err)
					return err
				}
			}
		}

		// update expected_return
		if arg.ExpectedReturn != nil {
			arg := UpdateExpectedReturnPortfolioProfileParams{
				ID:             arg.ProfileId,
				ExpectedReturn: *arg.ExpectedReturn,
			}
			_, err := q.UpdateExpectedReturnPortfolioProfile(ctx, arg)
			if err != nil {
				store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - UpdateExpectedReturnPortfolioProfile: %v", err)
				return err
			}
		}

		// update is_new_buy_point
		if arg.IsNewBuyPoint != nil {
			arg := UpdateIsNewBuyPointPortfolioProfileParams{
				ID:            arg.ProfileId,
				IsNewBuyPoint: *arg.IsNewBuyPoint,
			}
			_, err := q.UpdateIsNewBuyPointPortfolioProfile(ctx, arg)
			if err != nil {
				store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - UpdateIsNewBuyPointPortfolioProfile: %v", err)
				return err
			}
		}

		// table: assets
		if arg.Assets != nil {
			// delete assets current
			if len(assetIds) > 0 {
				for _, id := range assetIds {
					err := q.DeleteListAssetsById(ctx, id)
					if err != nil {
						store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - DeleteListAssetsById: %v", err)
						return err
					}
				}
			}

			// check case update empty array
			isAddUpdate := true
			if arg.Assets[0].TickerName == "0" && arg.Assets[0].Price == 0.0 && arg.Assets[0].Allocation == 0.0 {
				isAddUpdate = false
			}

			if isAddUpdate {
				// add assets new
				for _, asset := range arg.Assets {
					argAssest := CreateAssetParams{
						PortfolioProfileID: arg.ProfileId,
						TickerName:         asset.TickerName,
						Price:              asset.Price,
						Allocation:         asset.Allocation,
					}

					_, err := q.CreateAsset(ctx, argAssest)
					if err != nil {
						store.logger.Sugar().Infof("\n error UpdatePortfolioProfileTx - CreateAsset: %v", err)
						return err
					}
				}
			}
		}

		return nil
	})

	result.ProfileId = arg.ProfileId
	return result, err
}

// DELETE
type DeletePortfolioProfileTxParams struct {
	ProfileId string
}

type DeletePortfolioProfileTxResult struct {
	ProfileId string
}

func (store *SQLStore) DeletePortfolioProfileTx(ctx context.Context, arg DeletePortfolioProfileTxParams) (DeletePortfolioProfileTxResult, error) {
	var result DeletePortfolioProfileTxResult

	/*
		- assets
	*/
	errGet := make(chan error, 2)

	// get list assets by profile id
	assetsCh := make(chan []HarmonixBusinessAsset)
	go func() {
		assets, err := store.GetAssetsByProfileId(ctx, arg.ProfileId)
		errGet <- err

		// push data to channel
		assetsCh <- assets
	}()

	assets := <-assetsCh

	close(errGet)
	// Collect and handle errors
	for err := range errGet {
		if err != nil {
			store.logger.Sugar().Infof("\n error DeletePortfolioProfileTx - errGet: %v", err)
			return result, err
		}
	}

	err := store.execTx(ctx, func(q *Queries) error {
		// table: portfolio profile
		err := q.DeletePortfolioProfile(ctx, arg.ProfileId)
		if err != nil {
			store.logger.Sugar().Infof("\n error DeletePortfolioProfileTx - DeletePortfolioProfile: %v", err)
			return err
		}

		// table: assets
		if len(assets) > 0 {
			for _, asset := range assets {
				argAssest := DeleteAssetParams{
					PortfolioProfileID: asset.PortfolioProfileID,
					TickerName:         asset.TickerName,
				}

				err := q.DeleteAsset(ctx, argAssest)
				if err != nil {
					store.logger.Sugar().Infof("\n error DeletePortfolioProfileTx - DeleteAsset: %v", err)
					return err
				}
			}
		}

		// table: hrn_profile_account
		err = q.DeleteAllLinkedProfileAccountByProfileId(ctx, arg.ProfileId)
		if err != nil {
			store.logger.Sugar().Infof("\n error DeletePortfolioProfileTx - DeleteAllLinkedProfileAccountByProfileId: %v", err)
			return err
		}

		return err
	})

	result.ProfileId = arg.ProfileId
	return result, err
}
