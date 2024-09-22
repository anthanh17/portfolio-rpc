-- name: CreateUserPortfolio :one
INSERT INTO hamonix_business.u_portfolio (
  user_id,
  portfolio_id
) VALUES (
  $1, $2
) RETURNING *;
