-- name: CreateAsset :one
INSERT INTO harmonix_business.assets (
  portfolio_profile_id,
  ticker_name,
  price,
  allocation
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetListAssetsByPortfolioId :many
SELECT * FROM harmonix_business.assets
WHERE portfolio_profile_id = $1
ORDER BY id;

-- name: GetListAssetIdsByPortfolioId :many
SELECT id FROM harmonix_business.assets
WHERE portfolio_profile_id = $1;

-- name: GetAssetsByProfileId :many
SELECT * FROM harmonix_business.assets
WHERE portfolio_profile_id = $1;

-- name: DeleteListAssetsById :exec
DELETE FROM harmonix_business.assets
WHERE id = $1;

-- name: DeleteAsset :exec
DELETE FROM harmonix_business.assets
WHERE portfolio_profile_id = $1 AND ticker_name = $2;
