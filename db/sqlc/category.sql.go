// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: category.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countCategoriesByUserID = `-- name: CountCategoriesByUserID :one
SELECT COUNT(id)
FROM hamonix_business.u_categories
WHERE user_id = $1
`

func (q *Queries) CountCategoriesByUserID(ctx context.Context, userID string) (int64, error) {
	row := q.db.QueryRow(ctx, countCategoriesByUserID, userID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countPCategoryByCategoryId = `-- name: CountPCategoryByCategoryId :one
SELECT COUNT(id)
FROM hamonix_business.p_categories
WHERE category_id = $1
`

func (q *Queries) CountPCategoryByCategoryId(ctx context.Context, categoryID pgtype.Text) (int64, error) {
	row := q.db.QueryRow(ctx, countPCategoryByCategoryId, categoryID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countProfilesInCategory = `-- name: CountProfilesInCategory :one
SELECT COUNT(portfolio_id)
FROM hamonix_business.p_categories
WHERE category_id = $1
`

func (q *Queries) CountProfilesInCategory(ctx context.Context, categoryID pgtype.Text) (int64, error) {
	row := q.db.QueryRow(ctx, countProfilesInCategory, categoryID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createPortfolioCategory = `-- name: CreatePortfolioCategory :one
INSERT INTO hamonix_business.portfolio_categories (
    id,
  name,
  description
) VALUES (
  $1, $2, $3
) RETURNING id, name, description, created_at, updated_at
`

type CreatePortfolioCategoryParams struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
}

func (q *Queries) CreatePortfolioCategory(ctx context.Context, arg CreatePortfolioCategoryParams) (HamonixBusinessPortfolioCategory, error) {
	row := q.db.QueryRow(ctx, createPortfolioCategory, arg.ID, arg.Name, arg.Description)
	var i HamonixBusinessPortfolioCategory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createUCategory = `-- name: CreateUCategory :one
INSERT INTO hamonix_business.u_categories (
  category_id,
  user_id
) VALUES (
  $1, $2
) RETURNING id, category_id, user_id
`

type CreateUCategoryParams struct {
	CategoryID pgtype.Text `json:"category_id"`
	UserID     string      `json:"user_id"`
}

func (q *Queries) CreateUCategory(ctx context.Context, arg CreateUCategoryParams) (HamonixBusinessUCategory, error) {
	row := q.db.QueryRow(ctx, createUCategory, arg.CategoryID, arg.UserID)
	var i HamonixBusinessUCategory
	err := row.Scan(&i.ID, &i.CategoryID, &i.UserID)
	return i, err
}

const deletePortfolioCategory = `-- name: DeletePortfolioCategory :exec
DELETE FROM hamonix_business.portfolio_categories
WHERE id = $1
`

func (q *Queries) DeletePortfolioCategory(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deletePortfolioCategory, id)
	return err
}

const getCategoryInfo = `-- name: GetCategoryInfo :one
SELECT id, name, description, created_at, updated_at FROM hamonix_business.portfolio_categories
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetCategoryInfo(ctx context.Context, id string) (HamonixBusinessPortfolioCategory, error) {
	row := q.db.QueryRow(ctx, getCategoryInfo, id)
	var i HamonixBusinessPortfolioCategory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPCategoryByCategoryId = `-- name: GetPCategoryByCategoryId :many
SELECT id, portfolio_id, category_id FROM hamonix_business.p_categories
WHERE category_id = $1
`

func (q *Queries) GetPCategoryByCategoryId(ctx context.Context, categoryID pgtype.Text) ([]HamonixBusinessPCategory, error) {
	rows, err := q.db.Query(ctx, getPCategoryByCategoryId, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []HamonixBusinessPCategory{}
	for rows.Next() {
		var i HamonixBusinessPCategory
		if err := rows.Scan(&i.ID, &i.PortfolioID, &i.CategoryID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPCategoryByCategoryIdPaging = `-- name: GetPCategoryByCategoryIdPaging :many
SELECT portfolio_id FROM hamonix_business.p_categories
WHERE category_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type GetPCategoryByCategoryIdPagingParams struct {
	CategoryID pgtype.Text `json:"category_id"`
	Limit      int32       `json:"limit"`
	Offset     int32       `json:"offset"`
}

func (q *Queries) GetPCategoryByCategoryIdPaging(ctx context.Context, arg GetPCategoryByCategoryIdPagingParams) ([]string, error) {
	rows, err := q.db.Query(ctx, getPCategoryByCategoryIdPaging, arg.CategoryID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var portfolio_id string
		if err := rows.Scan(&portfolio_id); err != nil {
			return nil, err
		}
		items = append(items, portfolio_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPortfolioCategoryById = `-- name: GetPortfolioCategoryById :one
SELECT id, name, description, created_at, updated_at FROM hamonix_business.portfolio_categories
WHERE id = $1
`

func (q *Queries) GetPortfolioCategoryById(ctx context.Context, id string) (HamonixBusinessPortfolioCategory, error) {
	row := q.db.QueryRow(ctx, getPortfolioCategoryById, id)
	var i HamonixBusinessPortfolioCategory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUCategoryByUserId = `-- name: GetUCategoryByUserId :many
SELECT id, category_id, user_id FROM hamonix_business.u_categories
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type GetUCategoryByUserIdParams struct {
	UserID string `json:"user_id"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) GetUCategoryByUserId(ctx context.Context, arg GetUCategoryByUserIdParams) ([]HamonixBusinessUCategory, error) {
	rows, err := q.db.Query(ctx, getUCategoryByUserId, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []HamonixBusinessUCategory{}
	for rows.Next() {
		var i HamonixBusinessUCategory
		if err := rows.Scan(&i.ID, &i.CategoryID, &i.UserID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePortfolioCategory = `-- name: UpdatePortfolioCategory :one
UPDATE hamonix_business.portfolio_categories
SET name = $1, description = $2
WHERE id = $3 RETURNING id, name, description, created_at, updated_at
`

type UpdatePortfolioCategoryParams struct {
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
	ID          string      `json:"id"`
}

func (q *Queries) UpdatePortfolioCategory(ctx context.Context, arg UpdatePortfolioCategoryParams) (HamonixBusinessPortfolioCategory, error) {
	row := q.db.QueryRow(ctx, updatePortfolioCategory, arg.Name, arg.Description, arg.ID)
	var i HamonixBusinessPortfolioCategory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
