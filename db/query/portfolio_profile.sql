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

-- name: UpdateTickerPrice :one
UPDATE hamonix_business.ticker_prices
SET
  open = $2,
  low = $3,
  close = $4,
  date = $5
WHERE ticker_id = $1
RETURNING *;

-- name: UpdatePortfolioCategory :one
UPDATE hamonix_business.portfolio_categories
SET
  name = $2,
  description = $3
WHERE id = $1
RETURNING *;

-- name: UpdatePCategory :one
UPDATE hamonix_business.p_categories
SET
  category_id = $2
WHERE portfolio_id = $1
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

-- name: UpdateEqWhitelable :one
UPDATE hamonix_business.eq_whitelables
SET
  name = $2,
  url = $3,
  description = $4
WHERE id = $1
RETURNING *;

-- name: UpdateEqBackoffice :one
UPDATE hamonix_business.eq_backoffices
SET
  name = $2,
  description = $3
WHERE whitelable_id = $1
RETURNING *;

-- name: UpdateEqOrganization :one
UPDATE hamonix_business.eq_organizations
SET
  code = $2,
  description = $3
WHERE backoffice_id = $1
RETURNING *;

-- name: UpdateEqBranch :one
UPDATE hamonix_business.eq_branchs
SET
  code = $2
WHERE id = $1
RETURNING *;

-- name: UpdateEqAdvisor :one
UPDATE hamonix_business.eq_advisors
SET
  code = $2
WHERE id = $1
RETURNING *;

-- name: UpdateEqAccount :one
UPDATE hamonix_business.eq_accounts
SET
  code = $2
WHERE advisor_id = $1
RETURNING *;
