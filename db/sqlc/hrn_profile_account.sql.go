// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: hrn_profile_account.sql

package db

import (
	"context"
)

const countAccountsLinkedProfileByProfileId = `-- name: CountAccountsLinkedProfileByProfileId :one
SELECT COUNT(account_id)
FROM harmonix_business.hrn_profile_account
WHERE profile_id = $1
`

func (q *Queries) CountAccountsLinkedProfileByProfileId(ctx context.Context, profileID string) (int64, error) {
	row := q.db.QueryRow(ctx, countAccountsLinkedProfileByProfileId, profileID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createLinkedProfileToAccount = `-- name: CreateLinkedProfileToAccount :one
INSERT INTO harmonix_business.hrn_profile_account (
  profile_id,
  account_id
) VALUES (
  $1, $2
) RETURNING id, profile_id, account_id, updated_at
`

type CreateLinkedProfileToAccountParams struct {
	ProfileID string `json:"profile_id"`
	AccountID string `json:"account_id"`
}

func (q *Queries) CreateLinkedProfileToAccount(ctx context.Context, arg CreateLinkedProfileToAccountParams) (HarmonixBusinessHrnProfileAccount, error) {
	row := q.db.QueryRow(ctx, createLinkedProfileToAccount, arg.ProfileID, arg.AccountID)
	var i HarmonixBusinessHrnProfileAccount
	err := row.Scan(
		&i.ID,
		&i.ProfileID,
		&i.AccountID,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAllLinkedProfileAccountByProfileId = `-- name: DeleteAllLinkedProfileAccountByProfileId :exec
DELETE FROM harmonix_business.hrn_profile_account
WHERE profile_id = $1
`

func (q *Queries) DeleteAllLinkedProfileAccountByProfileId(ctx context.Context, profileID string) error {
	_, err := q.db.Exec(ctx, deleteAllLinkedProfileAccountByProfileId, profileID)
	return err
}

const deleteLinkedProfileAccount = `-- name: DeleteLinkedProfileAccount :exec
DELETE FROM harmonix_business.hrn_profile_account
WHERE profile_id = $1 and account_id = $2
`

type DeleteLinkedProfileAccountParams struct {
	ProfileID string `json:"profile_id"`
	AccountID string `json:"account_id"`
}

func (q *Queries) DeleteLinkedProfileAccount(ctx context.Context, arg DeleteLinkedProfileAccountParams) error {
	_, err := q.db.Exec(ctx, deleteLinkedProfileAccount, arg.ProfileID, arg.AccountID)
	return err
}

const getListAccountsLinkedProfileByProfileId = `-- name: GetListAccountsLinkedProfileByProfileId :many
SELECT account_id FROM harmonix_business.hrn_profile_account
WHERE profile_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type GetListAccountsLinkedProfileByProfileIdParams struct {
	ProfileID string `json:"profile_id"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}

func (q *Queries) GetListAccountsLinkedProfileByProfileId(ctx context.Context, arg GetListAccountsLinkedProfileByProfileIdParams) ([]string, error) {
	rows, err := q.db.Query(ctx, getListAccountsLinkedProfileByProfileId, arg.ProfileID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var account_id string
		if err := rows.Scan(&account_id); err != nil {
			return nil, err
		}
		items = append(items, account_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
