-- name: CreateUserPortfolio :one
INSERT INTO hamonix_business.u_portfolio (
  user_id,
  portfolio_id
) VALUES (
  $1, $2
) RETURNING *;


-- name: CountProfilesInUserPortfolio :one
SELECT COUNT(portfolio_id)
FROM hamonix_business.u_portfolio
WHERE user_id = $1;

-- name: GetUPortfolioByUserId :many
SELECT portfolio_id FROM hamonix_business.u_portfolio
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;
