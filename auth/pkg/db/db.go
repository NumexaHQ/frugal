package db

import (
	"context"

	postgresql_db "github.com/NumexaHQ/captainCache/numexa-common/postgresql/postgresql-db"
	"github.com/NumexaHQ/captainCache/pkg/db/postgres"
)

func New(d string) DB {
	switch d {
	case "postgres":
		return &postgres.Postgres{}
	default:
		return nil
	}
}

type DB interface {
	Init() error
	RegisterUser(ctx context.Context, user postgresql_db.User) (postgresql_db.User, error)
	GetUserByEmail(ctx context.Context, email string) (postgresql_db.User, error)
	CreateProject(ctx context.Context, project postgresql_db.Project) (postgresql_db.Project, error)
	CreateProjectUser(ctx context.Context, projectUser postgresql_db.ProjectUser) (postgresql_db.ProjectUser, error)
	GetProjectUsers(ctx context.Context, projectID int32) ([]postgresql_db.ProjectUser, error)
	GetUsersByProjectId(ctx context.Context, projectID int32) ([]postgresql_db.User, error)
	GetProject(ctx context.Context, projectID int32) (postgresql_db.Project, error)
	CreateOrganization(ctx context.Context, organization postgresql_db.Organization) (postgresql_db.Organization, error)
	CreateApiKey(ctx context.Context, apiKey postgresql_db.NxaApiKey) (postgresql_db.NxaApiKey, error)
	GetAPIkeyByUserId(ctx context.Context, userID int32) ([]postgresql_db.NxaApiKey, error)
	GetUserById(ctx context.Context, id int32) (postgresql_db.User, error)
	GetAPIkeyByApiKey(ctx context.Context, apiKey string) (postgresql_db.NxaApiKey, error)
	GetProjectUserByProjectIDAndUserID(ctx context.Context, projectID int32, userID int32) (postgresql_db.ProjectUser, error)
	GetAPIKeyByNameAndProjectId(ctx context.Context, name string, projectID int32) (postgresql_db.NxaApiKey, error)
	GetAllApiKeysByUserId(ctx context.Context, userID int32) ([]postgresql_db.GetAllApiKeysByUserIdRow, error)
	GetProjectsByOrgId(ctx context.Context, orgID int32) ([]postgresql_db.Project, error)
	UpdateUserLastLogin(ctx context.Context, user postgresql_db.User) error
}
