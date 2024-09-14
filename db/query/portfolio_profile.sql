-- name: CreatePortfolio :one
INSERT INTO hamonix_business.portfolios (
  id,
  name,
  privacy
) VALUES (
  $1, $2, $3
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

-- name: CreateTickerPrice :one
INSERT INTO hamonix_business.ticker_prices (
  ticker_id,
  open,
  low,
  close,
  date
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: CreatePortfolioCategory :one
INSERT INTO hamonix_business.portfolio_categories (
  name,
  description
) VALUES (
  $1, $2
) RETURNING *;

-- name: CreatePCategory :one
INSERT INTO hamonix_business.p_categories (
  portfolio_id,
  category_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: CreatePBranche :one
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

-- name: CreateEqWhitelable :one
INSERT INTO hamonix_business.eq_whitelables (
  id,
  name,
  url,
  description
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: CreateEqBackoffice :one
INSERT INTO hamonix_business.eq_backoffices (
  id,
  whitelable_id,
  name,
  description
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: CreateEqOrganization :one
INSERT INTO hamonix_business.eq_organizations (
  id,
  backoffice_id,
  code,
  description
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: CreateEqBranch :one
INSERT INTO hamonix_business.eq_branchs (
  id,
  code
) VALUES (
  $1, $2
) RETURNING *;

-- name: CreateEqAdvisor :one
INSERT INTO hamonix_business.eq_advisors (
  id,
  code,
  description
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: CreateEqAccount :one
INSERT INTO hamonix_business.eq_accounts (
  id,
  advisor_id,
  code
) VALUES (
  $1, $2, $3
) RETURNING *;
