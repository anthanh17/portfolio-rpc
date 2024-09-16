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
