-- name: CreateLinkedProfileToAccount :one
INSERT INTO harmonix_business.hrn_profile_account (
  profile_id,
  account_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetListAccountsLinkedProfileByProfileId :many
SELECT account_id FROM harmonix_business.hrn_profile_account
WHERE profile_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: CountAccountsLinkedProfileByProfileId :one
SELECT COUNT(account_id)
FROM harmonix_business.hrn_profile_account
WHERE profile_id = $1;

-- name: DeleteLinkedProfileAccount :exec
DELETE FROM harmonix_business.hrn_profile_account
WHERE profile_id = $1 and account_id = $2;

-- name: DeleteAllLinkedProfileAccountByProfileId :exec
DELETE FROM harmonix_business.hrn_profile_account
WHERE profile_id = $1;
