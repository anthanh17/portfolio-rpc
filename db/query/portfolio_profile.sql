-- name: CreatePortfolioProfile :one
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
) RETURNING *;

-- name: GetListProfileIdByUserId :many
SELECT id FROM harmonix_business.portfolio_profiles
WHERE author_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: GetListAdvisorsBranchesOrganizationsByProfileId :many
SELECT advisors, branches, organizations, accounts
FROM harmonix_business.portfolio_profiles
WHERE id = $1 LIMIT 1;

-- name: GetProfileInfoById :one
SELECT * FROM harmonix_business.portfolio_profiles
WHERE id = $1 LIMIT 1;

-- name: GetPrivacyProfileById :one
SELECT privacy FROM harmonix_business.portfolio_profiles
WHERE id = $1 LIMIT 1;

-- name: UpdateNamePortfolioProfile :one
UPDATE harmonix_business.portfolio_profiles
SET
  name = $2,
  updated_at = $3
WHERE id = $1
RETURNING *;

-- name: UpdatePrivacyPortfolioProfile :one
UPDATE harmonix_business.portfolio_profiles
SET
  privacy = $2,
  updated_at = $3
WHERE id = $1
RETURNING *;

-- name: UpdateAdvisorsPortfolioProfile :one
UPDATE harmonix_business.portfolio_profiles
SET
  advisors = $2
WHERE id = $1
RETURNING *;

-- name: UpdateBranchesPortfolioProfile :one
UPDATE harmonix_business.portfolio_profiles
SET
  branches = $2
WHERE id = $1
RETURNING *;

-- name: UpdateOrganizationsPortfolioProfile :one
UPDATE harmonix_business.portfolio_profiles
SET
  organizations = $2
WHERE id = $1
RETURNING *;

-- name: UpdateAccountsPortfolioProfile :one
UPDATE harmonix_business.portfolio_profiles
SET
  accounts = $2
WHERE id = $1
RETURNING *;

-- name: UpdateExpectedReturnPortfolioProfile :one
UPDATE harmonix_business.portfolio_profiles
SET
  expected_return = $2
WHERE id = $1
RETURNING *;

-- name: UpdateIsNewBuyPointPortfolioProfile :one
UPDATE harmonix_business.portfolio_profiles
SET
  is_new_buy_point = $2
WHERE id = $1
RETURNING *;

-- name: DeletePortfolioProfile :exec
DELETE FROM harmonix_business.portfolio_profiles
WHERE id = $1;

-- name: CheckIdExitsPortfolioProfile :one
SELECT EXISTS (
  SELECT 1
  FROM harmonix_business.portfolio_profiles
  WHERE id = $1
);

-- name: CountProfilesInUserPortfolioProfile :one
SELECT COUNT(id)
FROM harmonix_business.portfolio_profiles
WHERE author_id = $1;
