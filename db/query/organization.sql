-- name: GetEQOrganizationByID :one
SELECT * FROM hamonix_business.eq_organizations
WHERE id = $1 LIMIT 1;
