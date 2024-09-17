-- name: CreatePortfolio :one
INSERT INTO hamonix_business.portfolios (
  id,
  name,
  privacy,
  author_id
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: CreateAsset :one
INSERT INTO hamonix_business.assets (
  portfolio_id,
  ticker_id,
  price,
  allocation
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: CreatePCategory :one
INSERT INTO hamonix_business.p_categories (
  portfolio_id,
  category_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: CreatePBranch :one
INSERT INTO hamonix_business.p_branches (
  portfolio_id,
  branch_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: CreatePAdvisor :one
INSERT INTO hamonix_business.p_advisors (
  portfolio_id,
  advisor_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: CreatePOrganization :one
INSERT INTO hamonix_business.p_organizations (
  portfolio_id,
  organization_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: UpdatePortfolio :one
UPDATE hamonix_business.portfolios
SET
  name = $2,
  privacy = $3
WHERE id = $1
RETURNING *;

-- name: UpdateAsset :one
UPDATE hamonix_business.assets
SET
  price = $3,
  allocation = $4
WHERE portfolio_id = $1 AND ticker_id = $2
RETURNING *;

-- name: UpdatePCategory :one
UPDATE hamonix_business.p_categories
SET
  category_id = $3
WHERE portfolio_id = $1 AND category_id = $2
RETURNING *;

-- name: UpdatePBranch :one
UPDATE hamonix_business.p_branches
SET
  branch_id = $3
WHERE portfolio_id = $1 AND branch_id = $2
RETURNING *;

-- name: UpdatePAdvisor :one
UPDATE hamonix_business.p_advisors
SET
  advisor_id = $3
WHERE portfolio_id = $1 AND advisor_id = $2
RETURNING *;

-- name: UpdatePOrganization :one
UPDATE hamonix_business.p_organizations
SET
  organization_id = $3
WHERE portfolio_id = $1 AND organization_id = $2
RETURNING *;

-- name: DeletePortfolio :exec
DELETE FROM hamonix_business.portfolios
WHERE id = $1;

-- name: DeleteAsset :exec
DELETE FROM hamonix_business.assets
WHERE portfolio_id = $1 AND ticker_id = $2;

-- name: DeletePCategory :exec
DELETE FROM hamonix_business.p_categories
WHERE portfolio_id = $1 AND category_id = $2;

-- name: DeletePBranch :exec
DELETE FROM hamonix_business.p_branches
WHERE portfolio_id = $1 AND branch_id = $2;

-- name: DeletePAdvisor :exec
DELETE FROM hamonix_business.p_advisors
WHERE portfolio_id = $1 AND advisor_id = $2;

-- name: DeletePOrganization :exec
DELETE FROM hamonix_business.p_organizations
WHERE portfolio_id = $1 AND organization_id = $2;

-- name: GetAssetsByPortfolioId :many
SELECT * FROM hamonix_business.assets
WHERE portfolio_id = $1;

-- name: GetPCategoryByPortfolioId :many
SELECT * FROM hamonix_business.p_categories
WHERE portfolio_id = $1;

-- name: GetPBranchByPortfolioId :many
SELECT * FROM hamonix_business.p_branches
WHERE portfolio_id = $1;

-- name: GetPOrganizationByPortfolioId :many
SELECT * FROM hamonix_business.p_organizations
WHERE portfolio_id = $1;

-- name: GetPAdvisorByPortfolioId :many
SELECT * FROM hamonix_business.p_advisors
WHERE portfolio_id = $1;

-- name: GetProfilesByPortfolioId :one
SELECT * FROM hamonix_business.portfolios
WHERE id = $1 LIMIT 1;
