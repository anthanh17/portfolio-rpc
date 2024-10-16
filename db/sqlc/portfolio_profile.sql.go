// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: portfolio_profile.sql

package db

import (
	"context"
	"time"
)

const checkIdExitsPortfolioProfile = `-- name: CheckIdExitsPortfolioProfile :one
SELECT EXISTS (
  SELECT 1
  FROM harmonix_business.portfolio_profiles
  WHERE id = $1
)
`

func (q *Queries) CheckIdExitsPortfolioProfile(ctx context.Context, id string) (bool, error) {
	row := q.db.QueryRow(ctx, checkIdExitsPortfolioProfile, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const countProfilesInUserPortfolioProfile = `-- name: CountProfilesInUserPortfolioProfile :one
SELECT COUNT(id)
FROM harmonix_business.portfolio_profiles
WHERE author_id = $1
`

func (q *Queries) CountProfilesInUserPortfolioProfile(ctx context.Context, authorID string) (int64, error) {
	row := q.db.QueryRow(ctx, countProfilesInUserPortfolioProfile, authorID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createPortfolioProfile = `-- name: CreatePortfolioProfile :one
INSERT INTO harmonix_business.portfolio_profiles (
  id,
  name,
  privacy,
  author_id,
  advisors,
  branches,
  organizations,
  accounts,
  expected_return,
  is_new_buy_point
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING id, name, privacy, author_id, advisors, branches, organizations, accounts, expected_return, is_new_buy_point, created_at, updated_at
`

type CreatePortfolioProfileParams struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Privacy        string   `json:"privacy"`
	AuthorID       string   `json:"author_id"`
	Advisors       []string `json:"advisors"`
	Branches       []string `json:"branches"`
	Organizations  []string `json:"organizations"`
	Accounts       []string `json:"accounts"`
	ExpectedReturn float64  `json:"expected_return"`
	IsNewBuyPoint  bool     `json:"is_new_buy_point"`
}

func (q *Queries) CreatePortfolioProfile(ctx context.Context, arg CreatePortfolioProfileParams) (HarmonixBusinessPortfolioProfile, error) {
	row := q.db.QueryRow(ctx, createPortfolioProfile,
		arg.ID,
		arg.Name,
		arg.Privacy,
		arg.AuthorID,
		arg.Advisors,
		arg.Branches,
		arg.Organizations,
		arg.Accounts,
		arg.ExpectedReturn,
		arg.IsNewBuyPoint,
	)
	var i HarmonixBusinessPortfolioProfile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Privacy,
		&i.AuthorID,
		&i.Advisors,
		&i.Branches,
		&i.Organizations,
		&i.Accounts,
		&i.ExpectedReturn,
		&i.IsNewBuyPoint,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deletePortfolioProfile = `-- name: DeletePortfolioProfile :exec
DELETE FROM harmonix_business.portfolio_profiles
WHERE id = $1
`

func (q *Queries) DeletePortfolioProfile(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deletePortfolioProfile, id)
	return err
}

const getListAdvisorsBranchesOrganizationsByProfileId = `-- name: GetListAdvisorsBranchesOrganizationsByProfileId :many
SELECT advisors, branches, organizations, accounts
FROM harmonix_business.portfolio_profiles
WHERE id = $1 LIMIT 1
`

type GetListAdvisorsBranchesOrganizationsByProfileIdRow struct {
	Advisors      []string `json:"advisors"`
	Branches      []string `json:"branches"`
	Organizations []string `json:"organizations"`
	Accounts      []string `json:"accounts"`
}

func (q *Queries) GetListAdvisorsBranchesOrganizationsByProfileId(ctx context.Context, id string) ([]GetListAdvisorsBranchesOrganizationsByProfileIdRow, error) {
	rows, err := q.db.Query(ctx, getListAdvisorsBranchesOrganizationsByProfileId, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetListAdvisorsBranchesOrganizationsByProfileIdRow{}
	for rows.Next() {
		var i GetListAdvisorsBranchesOrganizationsByProfileIdRow
		if err := rows.Scan(
			&i.Advisors,
			&i.Branches,
			&i.Organizations,
			&i.Accounts,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getListProfileIdByUserId = `-- name: GetListProfileIdByUserId :many
SELECT id FROM harmonix_business.portfolio_profiles
WHERE author_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3
`

type GetListProfileIdByUserIdParams struct {
	AuthorID string `json:"author_id"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

func (q *Queries) GetListProfileIdByUserId(ctx context.Context, arg GetListProfileIdByUserIdParams) ([]string, error) {
	rows, err := q.db.Query(ctx, getListProfileIdByUserId, arg.AuthorID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPrivacyProfileById = `-- name: GetPrivacyProfileById :one
SELECT privacy FROM harmonix_business.portfolio_profiles
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPrivacyProfileById(ctx context.Context, id string) (string, error) {
	row := q.db.QueryRow(ctx, getPrivacyProfileById, id)
	var privacy string
	err := row.Scan(&privacy)
	return privacy, err
}

const getProfileInfoById = `-- name: GetProfileInfoById :one
SELECT id, name, privacy, author_id, advisors, branches, organizations, accounts, expected_return, is_new_buy_point, created_at, updated_at FROM harmonix_business.portfolio_profiles
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetProfileInfoById(ctx context.Context, id string) (HarmonixBusinessPortfolioProfile, error) {
	row := q.db.QueryRow(ctx, getProfileInfoById, id)
	var i HarmonixBusinessPortfolioProfile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Privacy,
		&i.AuthorID,
		&i.Advisors,
		&i.Branches,
		&i.Organizations,
		&i.Accounts,
		&i.ExpectedReturn,
		&i.IsNewBuyPoint,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateAccountsPortfolioProfile = `-- name: UpdateAccountsPortfolioProfile :one
UPDATE harmonix_business.portfolio_profiles
SET
  accounts = $2
WHERE id = $1
RETURNING id, name, privacy, author_id, advisors, branches, organizations, accounts, expected_return, is_new_buy_point, created_at, updated_at
`

type UpdateAccountsPortfolioProfileParams struct {
	ID       string   `json:"id"`
	Accounts []string `json:"accounts"`
}

func (q *Queries) UpdateAccountsPortfolioProfile(ctx context.Context, arg UpdateAccountsPortfolioProfileParams) (HarmonixBusinessPortfolioProfile, error) {
	row := q.db.QueryRow(ctx, updateAccountsPortfolioProfile, arg.ID, arg.Accounts)
	var i HarmonixBusinessPortfolioProfile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Privacy,
		&i.AuthorID,
		&i.Advisors,
		&i.Branches,
		&i.Organizations,
		&i.Accounts,
		&i.ExpectedReturn,
		&i.IsNewBuyPoint,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateAdvisorsPortfolioProfile = `-- name: UpdateAdvisorsPortfolioProfile :one
UPDATE harmonix_business.portfolio_profiles
SET
  advisors = $2
WHERE id = $1
RETURNING id, name, privacy, author_id, advisors, branches, organizations, accounts, expected_return, is_new_buy_point, created_at, updated_at
`

type UpdateAdvisorsPortfolioProfileParams struct {
	ID       string   `json:"id"`
	Advisors []string `json:"advisors"`
}

func (q *Queries) UpdateAdvisorsPortfolioProfile(ctx context.Context, arg UpdateAdvisorsPortfolioProfileParams) (HarmonixBusinessPortfolioProfile, error) {
	row := q.db.QueryRow(ctx, updateAdvisorsPortfolioProfile, arg.ID, arg.Advisors)
	var i HarmonixBusinessPortfolioProfile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Privacy,
		&i.AuthorID,
		&i.Advisors,
		&i.Branches,
		&i.Organizations,
		&i.Accounts,
		&i.ExpectedReturn,
		&i.IsNewBuyPoint,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateBranchesPortfolioProfile = `-- name: UpdateBranchesPortfolioProfile :one
UPDATE harmonix_business.portfolio_profiles
SET
  branches = $2
WHERE id = $1
RETURNING id, name, privacy, author_id, advisors, branches, organizations, accounts, expected_return, is_new_buy_point, created_at, updated_at
`

type UpdateBranchesPortfolioProfileParams struct {
	ID       string   `json:"id"`
	Branches []string `json:"branches"`
}

func (q *Queries) UpdateBranchesPortfolioProfile(ctx context.Context, arg UpdateBranchesPortfolioProfileParams) (HarmonixBusinessPortfolioProfile, error) {
	row := q.db.QueryRow(ctx, updateBranchesPortfolioProfile, arg.ID, arg.Branches)
	var i HarmonixBusinessPortfolioProfile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Privacy,
		&i.AuthorID,
		&i.Advisors,
		&i.Branches,
		&i.Organizations,
		&i.Accounts,
		&i.ExpectedReturn,
		&i.IsNewBuyPoint,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateExpectedReturnPortfolioProfile = `-- name: UpdateExpectedReturnPortfolioProfile :one
UPDATE harmonix_business.portfolio_profiles
SET
  expected_return = $2
WHERE id = $1
RETURNING id, name, privacy, author_id, advisors, branches, organizations, accounts, expected_return, is_new_buy_point, created_at, updated_at
`

type UpdateExpectedReturnPortfolioProfileParams struct {
	ID             string  `json:"id"`
	ExpectedReturn float64 `json:"expected_return"`
}

func (q *Queries) UpdateExpectedReturnPortfolioProfile(ctx context.Context, arg UpdateExpectedReturnPortfolioProfileParams) (HarmonixBusinessPortfolioProfile, error) {
	row := q.db.QueryRow(ctx, updateExpectedReturnPortfolioProfile, arg.ID, arg.ExpectedReturn)
	var i HarmonixBusinessPortfolioProfile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Privacy,
		&i.AuthorID,
		&i.Advisors,
		&i.Branches,
		&i.Organizations,
		&i.Accounts,
		&i.ExpectedReturn,
		&i.IsNewBuyPoint,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateIsNewBuyPointPortfolioProfile = `-- name: UpdateIsNewBuyPointPortfolioProfile :one
UPDATE harmonix_business.portfolio_profiles
SET
  is_new_buy_point = $2
WHERE id = $1
RETURNING id, name, privacy, author_id, advisors, branches, organizations, accounts, expected_return, is_new_buy_point, created_at, updated_at
`

type UpdateIsNewBuyPointPortfolioProfileParams struct {
	ID            string `json:"id"`
	IsNewBuyPoint bool   `json:"is_new_buy_point"`
}

func (q *Queries) UpdateIsNewBuyPointPortfolioProfile(ctx context.Context, arg UpdateIsNewBuyPointPortfolioProfileParams) (HarmonixBusinessPortfolioProfile, error) {
	row := q.db.QueryRow(ctx, updateIsNewBuyPointPortfolioProfile, arg.ID, arg.IsNewBuyPoint)
	var i HarmonixBusinessPortfolioProfile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Privacy,
		&i.AuthorID,
		&i.Advisors,
		&i.Branches,
		&i.Organizations,
		&i.Accounts,
		&i.ExpectedReturn,
		&i.IsNewBuyPoint,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateNamePortfolioProfile = `-- name: UpdateNamePortfolioProfile :one
UPDATE harmonix_business.portfolio_profiles
SET
  name = $2,
  updated_at = $3
WHERE id = $1
RETURNING id, name, privacy, author_id, advisors, branches, organizations, accounts, expected_return, is_new_buy_point, created_at, updated_at
`

type UpdateNamePortfolioProfileParams struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) UpdateNamePortfolioProfile(ctx context.Context, arg UpdateNamePortfolioProfileParams) (HarmonixBusinessPortfolioProfile, error) {
	row := q.db.QueryRow(ctx, updateNamePortfolioProfile, arg.ID, arg.Name, arg.UpdatedAt)
	var i HarmonixBusinessPortfolioProfile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Privacy,
		&i.AuthorID,
		&i.Advisors,
		&i.Branches,
		&i.Organizations,
		&i.Accounts,
		&i.ExpectedReturn,
		&i.IsNewBuyPoint,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateOrganizationsPortfolioProfile = `-- name: UpdateOrganizationsPortfolioProfile :one
UPDATE harmonix_business.portfolio_profiles
SET
  organizations = $2
WHERE id = $1
RETURNING id, name, privacy, author_id, advisors, branches, organizations, accounts, expected_return, is_new_buy_point, created_at, updated_at
`

type UpdateOrganizationsPortfolioProfileParams struct {
	ID            string   `json:"id"`
	Organizations []string `json:"organizations"`
}

func (q *Queries) UpdateOrganizationsPortfolioProfile(ctx context.Context, arg UpdateOrganizationsPortfolioProfileParams) (HarmonixBusinessPortfolioProfile, error) {
	row := q.db.QueryRow(ctx, updateOrganizationsPortfolioProfile, arg.ID, arg.Organizations)
	var i HarmonixBusinessPortfolioProfile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Privacy,
		&i.AuthorID,
		&i.Advisors,
		&i.Branches,
		&i.Organizations,
		&i.Accounts,
		&i.ExpectedReturn,
		&i.IsNewBuyPoint,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updatePrivacyPortfolioProfile = `-- name: UpdatePrivacyPortfolioProfile :one
UPDATE harmonix_business.portfolio_profiles
SET
  privacy = $2,
  updated_at = $3
WHERE id = $1
RETURNING id, name, privacy, author_id, advisors, branches, organizations, accounts, expected_return, is_new_buy_point, created_at, updated_at
`

type UpdatePrivacyPortfolioProfileParams struct {
	ID        string    `json:"id"`
	Privacy   string    `json:"privacy"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) UpdatePrivacyPortfolioProfile(ctx context.Context, arg UpdatePrivacyPortfolioProfileParams) (HarmonixBusinessPortfolioProfile, error) {
	row := q.db.QueryRow(ctx, updatePrivacyPortfolioProfile, arg.ID, arg.Privacy, arg.UpdatedAt)
	var i HarmonixBusinessPortfolioProfile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Privacy,
		&i.AuthorID,
		&i.Advisors,
		&i.Branches,
		&i.Organizations,
		&i.Accounts,
		&i.ExpectedReturn,
		&i.IsNewBuyPoint,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
