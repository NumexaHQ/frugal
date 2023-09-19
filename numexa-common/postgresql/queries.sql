
-- name: CreateOrganization :one
INSERT INTO organizations (name, created_at, updated_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateUser :one
INSERT INTO users (name, organization_id, email, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CreateProject :one
INSERT INTO projects (name, organization_id, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateProjectUser :one
INSERT INTO project_users (project_id, user_id, role_id)
VALUES ($1, $2, $3)
RETURNING *;


-- name: GetOrganization :one
SELECT * FROM organizations WHERE id = $1;

-- name: GetOrganizationByName :one
SELECT * FROM organizations WHERE name = $1;

-- name: GetUsers :many
SELECT * FROM users WHERE organization_id = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUsersByProjectId :many
SELECT * FROM users WHERE id IN (SELECT user_id FROM project_users WHERE project_id = $1);

-- name: GetProjects :many
SELECT * FROM projects WHERE organization_id = $1;

-- name: GetProject :one
SELECT * FROM projects WHERE id = $1;

-- name: GetProjectByName :one
SELECT * FROM projects WHERE name = $1;

-- name: GetProjectUsers :many
SELECT * FROM project_users WHERE project_id = $1;

-- name: GetProjectUser :one
SELECT * FROM project_users WHERE project_id = $1 AND user_id = $2;

-- name: GetTokenByProjectId :one
SELECT * FROM nxa_api_key WHERE project_id = $1 AND api_key = $2;

-- name: CreateApiKey :one
INSERT INTO nxa_api_key (name, api_key, user_id, project_id, created_at, updated_at, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetAPIkeyByUserId :one
SELECT * FROM nxa_api_key WHERE user_id = $1;

-- name: GetAPIkeyByApiKey :one
SELECT * FROM nxa_api_key WHERE api_key = $1;

-- name: GetAPIKeyByNameAndProjectId :one
SELECT * FROM nxa_api_key WHERE name = $1 AND project_id = $2;

-- name: GetProjectUserByProjectIDAndUserID :one
SELECT * FROM project_users WHERE project_id = $1 AND user_id = $2;

-- name: GetAllApiKeysByUserId :many
SELECT  created_at, updated_at, expires_at, name, project_id, user_id FROM nxa_api_key WHERE user_id = $1;

-- name: UpdateUserLastLogin :one
UPDATE users SET last_login = $1, total_logins = total_logins + 1 WHERE id = $2 RETURNING *;

-- name: CreateProviderKey :one
INSERT INTO provider_keys (name, key_uuid, provider, creator_id, organization_id, project_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: CreateProviderSecret :one
INSERT INTO provider_secrets (provider_key_id, type, key, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetProviderKey :one
SELECT * FROM provider_keys WHERE id = $1;

-- name: GetProviderKeyByUUID :one
SELECT * FROM provider_keys WHERE key_uuid = $1;

-- name: GetProviderKeyByName :one
SELECT * FROM provider_keys WHERE name = $1;

-- name: GetProviderKeysByProjectID :many
SELECT * FROM provider_keys WHERE project_id = $1;

-- name: GetProviderSecret :one 
SELECT * FROM provider_secrets WHERE id = $1;

-- name: GetProviderSecretsByProviderKeyID :many
SELECT * FROM provider_secrets WHERE provider_key_id = $1;

-- name: GetProviderSecretsByProviderKeyID :many
SELECT * FROM provider_secrets WHERE provider_key_id = $1; 

-- name: GetProviderSecretByProviderKeyIDAndType :one
SELECT * FROM provider_secrets WHERE provider_key_id = $1 AND type = $2;

-- name: CreateSetting :one
INSERT INTO setting (key, value, visible, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSetting :one
SELECT * FROM setting WHERE key = $1;


