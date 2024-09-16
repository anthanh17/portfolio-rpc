-- name: GetEQBranchByID :one
SELECT * FROM hamonix_business.eq_branchs
WHERE id = $1 LIMIT 1;
