// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: queries.sql

package postgresql_db

import (
	"context"
	"database/sql"
	"time"
)

const createApiKey = `-- name: CreateApiKey :one
INSERT INTO nxa_api_key (name, api_key, user_id, project_id, created_at, updated_at, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, api_key, user_id, project_id, created_at, updated_at, expires_at
`

type CreateApiKeyParams struct {
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
	UserID    int32     `json:"user_id"`
	ProjectID int32     `json:"project_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (q *Queries) CreateApiKey(ctx context.Context, arg CreateApiKeyParams) (NxaApiKey, error) {
	row := q.db.QueryRowContext(ctx, createApiKey,
		arg.Name,
		arg.ApiKey,
		arg.UserID,
		arg.ProjectID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.ExpiresAt,
	)
	var i NxaApiKey
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ApiKey,
		&i.UserID,
		&i.ProjectID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const createOrganization = `-- name: CreateOrganization :one
INSERT INTO organizations (name, created_at, updated_at)
VALUES ($1, $2, $3)
RETURNING id, name, created_at, updated_at
`

type CreateOrganizationParams struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) CreateOrganization(ctx context.Context, arg CreateOrganizationParams) (Organization, error) {
	row := q.db.QueryRowContext(ctx, createOrganization, arg.Name, arg.CreatedAt, arg.UpdatedAt)
	var i Organization
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createProject = `-- name: CreateProject :one
INSERT INTO projects (name, organization_id, description)
VALUES ($1, $2, $3)
RETURNING id, organization_id, name, description
`

type CreateProjectParams struct {
	Name           string         `json:"name"`
	OrganizationID int32          `json:"organization_id"`
	Description    sql.NullString `json:"description"`
}

func (q *Queries) CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error) {
	row := q.db.QueryRowContext(ctx, createProject, arg.Name, arg.OrganizationID, arg.Description)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.OrganizationID,
		&i.Name,
		&i.Description,
	)
	return i, err
}

const createProjectUser = `-- name: CreateProjectUser :one
INSERT INTO project_users (project_id, user_id, role_id)
VALUES ($1, $2, $3)
RETURNING id, project_id, user_id, role_id
`

type CreateProjectUserParams struct {
	ProjectID int32 `json:"project_id"`
	UserID    int32 `json:"user_id"`
	RoleID    int32 `json:"role_id"`
}

func (q *Queries) CreateProjectUser(ctx context.Context, arg CreateProjectUserParams) (ProjectUser, error) {
	row := q.db.QueryRowContext(ctx, createProjectUser, arg.ProjectID, arg.UserID, arg.RoleID)
	var i ProjectUser
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.UserID,
		&i.RoleID,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (name, organization_id, email, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, organization_id, name, email, password, created_at, updated_at
`

type CreateUserParams struct {
	Name           string    `json:"name"`
	OrganizationID int32     `json:"organization_id"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Name,
		arg.OrganizationID,
		arg.Email,
		arg.Password,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.OrganizationID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAPIKeyByNameAndProjectId = `-- name: GetAPIKeyByNameAndProjectId :one
SELECT id, name, api_key, user_id, project_id, created_at, updated_at, expires_at FROM nxa_api_key WHERE name = $1 AND project_id = $2
`

type GetAPIKeyByNameAndProjectIdParams struct {
	Name      string `json:"name"`
	ProjectID int32  `json:"project_id"`
}

func (q *Queries) GetAPIKeyByNameAndProjectId(ctx context.Context, arg GetAPIKeyByNameAndProjectIdParams) (NxaApiKey, error) {
	row := q.db.QueryRowContext(ctx, getAPIKeyByNameAndProjectId, arg.Name, arg.ProjectID)
	var i NxaApiKey
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ApiKey,
		&i.UserID,
		&i.ProjectID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const getAPIkeyByApiKey = `-- name: GetAPIkeyByApiKey :one
SELECT id, name, api_key, user_id, project_id, created_at, updated_at, expires_at FROM nxa_api_key WHERE api_key = $1
`

func (q *Queries) GetAPIkeyByApiKey(ctx context.Context, apiKey string) (NxaApiKey, error) {
	row := q.db.QueryRowContext(ctx, getAPIkeyByApiKey, apiKey)
	var i NxaApiKey
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ApiKey,
		&i.UserID,
		&i.ProjectID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const getAPIkeyByUserId = `-- name: GetAPIkeyByUserId :one
SELECT id, name, api_key, user_id, project_id, created_at, updated_at, expires_at FROM nxa_api_key WHERE user_id = $1
`

func (q *Queries) GetAPIkeyByUserId(ctx context.Context, userID int32) (NxaApiKey, error) {
	row := q.db.QueryRowContext(ctx, getAPIkeyByUserId, userID)
	var i NxaApiKey
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ApiKey,
		&i.UserID,
		&i.ProjectID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const getAllApiKeysByUserId = `-- name: GetAllApiKeysByUserId :many
SELECT  created_at, updated_at, expires_at, name, project_id, user_id FROM nxa_api_key WHERE user_id = $1
`

type GetAllApiKeysByUserIdRow struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Name      string    `json:"name"`
	ProjectID int32     `json:"project_id"`
	UserID    int32     `json:"user_id"`
}

func (q *Queries) GetAllApiKeysByUserId(ctx context.Context, userID int32) ([]GetAllApiKeysByUserIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllApiKeysByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllApiKeysByUserIdRow
	for rows.Next() {
		var i GetAllApiKeysByUserIdRow
		if err := rows.Scan(
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ExpiresAt,
			&i.Name,
			&i.ProjectID,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOrganization = `-- name: GetOrganization :one
SELECT id, name, created_at, updated_at FROM organizations WHERE id = $1
`

func (q *Queries) GetOrganization(ctx context.Context, id int32) (Organization, error) {
	row := q.db.QueryRowContext(ctx, getOrganization, id)
	var i Organization
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getOrganizationByName = `-- name: GetOrganizationByName :one
SELECT id, name, created_at, updated_at FROM organizations WHERE name = $1
`

func (q *Queries) GetOrganizationByName(ctx context.Context, name string) (Organization, error) {
	row := q.db.QueryRowContext(ctx, getOrganizationByName, name)
	var i Organization
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProject = `-- name: GetProject :one
SELECT id, organization_id, name, description FROM projects WHERE id = $1
`

func (q *Queries) GetProject(ctx context.Context, id int32) (Project, error) {
	row := q.db.QueryRowContext(ctx, getProject, id)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.OrganizationID,
		&i.Name,
		&i.Description,
	)
	return i, err
}

const getProjectByName = `-- name: GetProjectByName :one
SELECT id, organization_id, name, description FROM projects WHERE name = $1
`

func (q *Queries) GetProjectByName(ctx context.Context, name string) (Project, error) {
	row := q.db.QueryRowContext(ctx, getProjectByName, name)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.OrganizationID,
		&i.Name,
		&i.Description,
	)
	return i, err
}

const getProjectUser = `-- name: GetProjectUser :one
SELECT id, project_id, user_id, role_id FROM project_users WHERE project_id = $1 AND user_id = $2
`

type GetProjectUserParams struct {
	ProjectID int32 `json:"project_id"`
	UserID    int32 `json:"user_id"`
}

func (q *Queries) GetProjectUser(ctx context.Context, arg GetProjectUserParams) (ProjectUser, error) {
	row := q.db.QueryRowContext(ctx, getProjectUser, arg.ProjectID, arg.UserID)
	var i ProjectUser
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.UserID,
		&i.RoleID,
	)
	return i, err
}

const getProjectUserByProjectIDAndUserID = `-- name: GetProjectUserByProjectIDAndUserID :one
SELECT id, project_id, user_id, role_id FROM project_users WHERE project_id = $1 AND user_id = $2
`

type GetProjectUserByProjectIDAndUserIDParams struct {
	ProjectID int32 `json:"project_id"`
	UserID    int32 `json:"user_id"`
}

func (q *Queries) GetProjectUserByProjectIDAndUserID(ctx context.Context, arg GetProjectUserByProjectIDAndUserIDParams) (ProjectUser, error) {
	row := q.db.QueryRowContext(ctx, getProjectUserByProjectIDAndUserID, arg.ProjectID, arg.UserID)
	var i ProjectUser
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.UserID,
		&i.RoleID,
	)
	return i, err
}

const getProjectUsers = `-- name: GetProjectUsers :many
SELECT id, project_id, user_id, role_id FROM project_users WHERE project_id = $1
`

func (q *Queries) GetProjectUsers(ctx context.Context, projectID int32) ([]ProjectUser, error) {
	rows, err := q.db.QueryContext(ctx, getProjectUsers, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ProjectUser
	for rows.Next() {
		var i ProjectUser
		if err := rows.Scan(
			&i.ID,
			&i.ProjectID,
			&i.UserID,
			&i.RoleID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProjects = `-- name: GetProjects :many
SELECT id, organization_id, name, description FROM projects WHERE organization_id = $1
`

func (q *Queries) GetProjects(ctx context.Context, organizationID int32) ([]Project, error) {
	rows, err := q.db.QueryContext(ctx, getProjects, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Project
	for rows.Next() {
		var i Project
		if err := rows.Scan(
			&i.ID,
			&i.OrganizationID,
			&i.Name,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTokenByProjectId = `-- name: GetTokenByProjectId :one
SELECT id, name, api_key, user_id, project_id, created_at, updated_at, expires_at FROM nxa_api_key WHERE project_id = $1 AND api_key = $2
`

type GetTokenByProjectIdParams struct {
	ProjectID int32  `json:"project_id"`
	ApiKey    string `json:"api_key"`
}

func (q *Queries) GetTokenByProjectId(ctx context.Context, arg GetTokenByProjectIdParams) (NxaApiKey, error) {
	row := q.db.QueryRowContext(ctx, getTokenByProjectId, arg.ProjectID, arg.ApiKey)
	var i NxaApiKey
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ApiKey,
		&i.UserID,
		&i.ProjectID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, organization_id, name, email, password, created_at, updated_at FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.OrganizationID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, organization_id, name, email, password, created_at, updated_at FROM users WHERE id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.OrganizationID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, organization_id, name, email, password, created_at, updated_at FROM users WHERE organization_id = $1
`

func (q *Queries) GetUsers(ctx context.Context, organizationID int32) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsers, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.OrganizationID,
			&i.Name,
			&i.Email,
			&i.Password,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsersByProjectId = `-- name: GetUsersByProjectId :many
SELECT id, organization_id, name, email, password, created_at, updated_at FROM users WHERE id IN (SELECT user_id FROM project_users WHERE project_id = $1)
`

func (q *Queries) GetUsersByProjectId(ctx context.Context, projectID int32) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsersByProjectId, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.OrganizationID,
			&i.Name,
			&i.Email,
			&i.Password,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
