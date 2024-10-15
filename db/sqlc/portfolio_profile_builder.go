package db

import (
	"portfolio-profile-rpc/rd_portfolio_rpc"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

// start UpdatePortfolioProfileTxParamsBuilder
type UpdatePortfolioProfileTxParamsBuilder struct {
	updateProfileTxParams *UpdatePortfolioProfileTxParams
}

// ProfileId
func NewUpdatePortfolioProfileTxParamsBuilder(profileId string) *UpdatePortfolioProfileTxParamsBuilder {
	return &UpdatePortfolioProfileTxParamsBuilder{
		updateProfileTxParams: &UpdatePortfolioProfileTxParams{
			ProfileId: profileId,
		},
	}
}

// ProfileName
func (u *UpdatePortfolioProfileTxParamsBuilder) WithProfileName(name *wrapperspb.StringValue) *UpdatePortfolioProfileTxParamsBuilder {
	if name == nil {
		u.updateProfileTxParams.ProfileName = nil
		return u
	}

	u.updateProfileTxParams.ProfileName = &name.Value
	return u
}

// Advisors
func (u *UpdatePortfolioProfileTxParamsBuilder) WithAdvisors(advisors []*wrapperspb.StringValue) *UpdatePortfolioProfileTxParamsBuilder {
	if advisors == nil {
		u.updateProfileTxParams.Advisors = nil
		return u
	}

	for _, value := range advisors {
		u.updateProfileTxParams.Advisors = append(u.updateProfileTxParams.Advisors, &value.Value)
	}

	return u
}

// Branches
func (u *UpdatePortfolioProfileTxParamsBuilder) WithBranches(branches []*wrapperspb.StringValue) *UpdatePortfolioProfileTxParamsBuilder {
	if branches == nil {
		u.updateProfileTxParams.Branches = nil
		return u
	}

	for _, value := range branches {
		u.updateProfileTxParams.Branches = append(u.updateProfileTxParams.Branches, &value.Value)
	}

	return u
}

// Organizations
func (u *UpdatePortfolioProfileTxParamsBuilder) WithOrganizations(organizations []*wrapperspb.StringValue) *UpdatePortfolioProfileTxParamsBuilder {
	if organizations == nil {
		u.updateProfileTxParams.Organizations = nil
		return u
	}

	for _, value := range organizations {
		u.updateProfileTxParams.Organizations = append(u.updateProfileTxParams.Organizations, &value.Value)
	}

	return u
}

// Accounts
func (u *UpdatePortfolioProfileTxParamsBuilder) WithAccounts(accounts []*wrapperspb.StringValue) *UpdatePortfolioProfileTxParamsBuilder {
	if accounts == nil {
		u.updateProfileTxParams.Accounts = nil
		return u
	}

	for _, value := range accounts {
		u.updateProfileTxParams.Accounts = append(u.updateProfileTxParams.Accounts, &value.Value)
	}

	return u
}

// ExpectedReturn
func (u *UpdatePortfolioProfileTxParamsBuilder) WithExpectedReturn(expectedReturn *wrapperspb.DoubleValue) *UpdatePortfolioProfileTxParamsBuilder {
	if expectedReturn == nil {
		u.updateProfileTxParams.ExpectedReturn = nil
		return u
	}

	u.updateProfileTxParams.ExpectedReturn = &expectedReturn.Value
	return u
}

// IsNewBuyPoint
func (u *UpdatePortfolioProfileTxParamsBuilder) WithIsNewBuyPoint(isNewBuyPoint *wrapperspb.BoolValue) *UpdatePortfolioProfileTxParamsBuilder {
	if isNewBuyPoint == nil {
		u.updateProfileTxParams.IsNewBuyPoint = nil
		return u
	}

	u.updateProfileTxParams.IsNewBuyPoint = &isNewBuyPoint.Value
	return u
}

// Assets
func (u *UpdatePortfolioProfileTxParamsBuilder) WithAssets(assets []*rd_portfolio_rpc.PortfolioAsset) *UpdatePortfolioProfileTxParamsBuilder {
	if assets == nil {
		u.updateProfileTxParams.Assets = nil
		return u
	}

	for _, value := range assets {
		u.updateProfileTxParams.Assets = append(u.updateProfileTxParams.Assets, &ProfileAsset{
			TickerName: value.TickerName.GetValue(),
			Allocation: value.Allocation.GetValue(),
			Price:      value.Price.GetValue(),
		})
	}

	return u
}

// Privacy
func (u *UpdatePortfolioProfileTxParamsBuilder) WithPrivacy(privacy *wrapperspb.StringValue) *UpdatePortfolioProfileTxParamsBuilder {
	if privacy == nil {
		u.updateProfileTxParams.Privacy = nil
		return u
	}

	u.updateProfileTxParams.Privacy = &privacy.Value
	return u
}

// Build
func (u *UpdatePortfolioProfileTxParamsBuilder) Build() *UpdatePortfolioProfileTxParams {
	return u.updateProfileTxParams
}

// end UpdatePortfolioProfileTxParamsBuilder
