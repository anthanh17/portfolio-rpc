-- name: GetEQAdvisorByID :one
SELECT * FROM hamonix_business.eq_advisors
WHERE id = $1 LIMIT 1;
