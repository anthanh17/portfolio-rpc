-- name: GetListAssetsByPortfolioId :many
SELECT * FROM hamonix_business.assets
WHERE portfolio_id = $1
ORDER BY id;
