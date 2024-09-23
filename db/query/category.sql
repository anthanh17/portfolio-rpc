-- name: CreatePortfolioCategory :one
INSERT INTO hamonix_business.portfolio_categories (
    id,
  name,
  description
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetPortfolioCategoryById :one
SELECT * FROM hamonix_business.portfolio_categories
WHERE id = $1;

-- name: UpdatePortfolioCategory :one
UPDATE hamonix_business.portfolio_categories
SET name = $1, description = $2
WHERE id = $3 RETURNING *;

-- name: DeletePortfolioCategory :exec
DELETE FROM hamonix_business.portfolio_categories
WHERE id = $1;

-- name: GetPCategoryByCategoryId :many
SELECT * FROM hamonix_business.p_categories
WHERE category_id = $1;

-- name: GetListProfileIdByCategoryId :many
SELECT portfolio_id FROM hamonix_business.p_categories
WHERE category_id = $1;

-- name: CountPCategoryByCategoryId :one
SELECT COUNT(id)
FROM hamonix_business.p_categories
WHERE category_id = $1;

-- name: GetPCategoryByCategoryIdPaging :many
SELECT portfolio_id FROM hamonix_business.p_categories
WHERE category_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;;

-- name: CreateUCategory :one
INSERT INTO hamonix_business.u_categories (
  category_id,
  user_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetUCategoryByUserId :many
SELECT * FROM hamonix_business.u_categories
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: CountCategoriesByUserID :one
SELECT COUNT(id)
FROM hamonix_business.u_categories
WHERE user_id = $1;

-- name: CountProfilesInCategory :one
SELECT COUNT(portfolio_id)
FROM hamonix_business.p_categories
WHERE category_id = $1;

-- name: GetCategoryInfo :one
SELECT * FROM hamonix_business.portfolio_categories
WHERE id = $1 LIMIT 1;
