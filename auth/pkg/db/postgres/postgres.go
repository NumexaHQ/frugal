package postgres

import (
	"context"
	"database/sql"

	"time"

	postgresql_db "github.com/NumexaHQ/captainCache/numexa-common/postgresql/postgresql-db"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// todo: use a connection pool
// todo: use a context
// todo: use a logger
// todo: use a config
// todo: seperate out the queries into a seperate file

func (p *Postgres) Init() error {
	connStr := "host=nxa-postgres port=5432 user=numexa password=numexa dbname=numexa sslmode=disable"
	// connStr = "host=localhost port=5432 user=numexa password=numexa dbname=numexa sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	p.db = db
	return nil
}

// todo: use this function to migrate the database
func migratePostgres() error {
	dbURL := "postgres://numexa:numexa@nxa-postgres:5432/numexa?sslmode=disable"
	// dbURL = "postgres://numexa:numexa@localhost:5432/numexa?sslmode=disable"

	migrationFile := "file:///usr/local/postgresql"
	m, err := migrate.New(
		migrationFile, dbURL)
	if err != nil {
		return err
	}
	return m.Up()
}

func getPostgresQueries(db *sql.DB) *postgresql_db.Queries {
	return postgresql_db.New(db)
}

func (p *Postgres) GetUserByEmail(ctx context.Context, email string) (postgresql_db.User, error) {
	queries := getPostgresQueries(p.db)

	user, err := queries.GetUserByEmail(ctx, email)
	if err != nil {
		return postgresql_db.User{}, err
	}

	return user, nil
}

func (p *Postgres) RegisterUser(ctx context.Context, user postgresql_db.User) (postgresql_db.User, error) {
	queries := getPostgresQueries(p.db)

	return queries.CreateUser(ctx, postgresql_db.CreateUserParams{
		Email:          user.Email,
		OrganizationID: user.OrganizationID,
		Name:           user.Name,
		Password:       user.Password,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
}

func (p *Postgres) GetUserById(ctx context.Context, id int32) (postgresql_db.User, error) {
	queries := getPostgresQueries(p.db)

	user, err := queries.GetUserById(ctx, id)
	if err != nil {
		return postgresql_db.User{}, err
	}

	return user, nil
}

func (p *Postgres) CreateProject(ctx context.Context, project postgresql_db.Project) (postgresql_db.Project, error) {
	queries := getPostgresQueries(p.db)

	project, err := queries.CreateProject(ctx, postgresql_db.CreateProjectParams{
		Name:           project.Name,
		OrganizationID: project.OrganizationID,
		Description:    project.Description,
	})

	return project, err
}

func (p *Postgres) CreateProjectUser(ctx context.Context, projectUser postgresql_db.ProjectUser) (postgresql_db.ProjectUser, error) {
	queries := getPostgresQueries(p.db)

	projectUser, err := queries.CreateProjectUser(ctx, postgresql_db.CreateProjectUserParams{
		ProjectID: projectUser.ProjectID,
		UserID:    projectUser.UserID,
		RoleID:    projectUser.RoleID,
	})

	return projectUser, err
}

func (p *Postgres) GetProjectUsers(ctx context.Context, projectID int32) ([]postgresql_db.ProjectUser, error) {
	queries := getPostgresQueries(p.db)

	projectUsers, err := queries.GetProjectUsers(ctx, int32(projectID))
	if err != nil {
		return nil, err
	}

	return projectUsers, nil
}

func (p *Postgres) GetUsersByProjectId(ctx context.Context, projectID int32) ([]postgresql_db.User, error) {
	queries := getPostgresQueries(p.db)

	users, err := queries.GetUsersByProjectId(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (p *Postgres) GetProject(ctx context.Context, projectID int32) (postgresql_db.Project, error) {
	queries := getPostgresQueries(p.db)

	project, err := queries.GetProject(ctx, projectID)
	if err != nil {
		return postgresql_db.Project{}, err
	}

	return project, nil
}

func (p *Postgres) CreateOrganization(ctx context.Context, organization postgresql_db.Organization) (postgresql_db.Organization, error) {
	queries := getPostgresQueries(p.db)

	organization, err := queries.CreateOrganization(ctx, postgresql_db.CreateOrganizationParams{
		Name:      organization.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	return organization, err
}

func (p *Postgres) GetOrganization(ctx context.Context, organizationID int32) (postgresql_db.Organization, error) {
	queries := getPostgresQueries(p.db)

	organization, err := queries.GetOrganization(ctx, organizationID)
	if err != nil {
		return postgresql_db.Organization{}, err
	}

	return organization, nil
}

func (p *Postgres) CreateApiKey(ctx context.Context, apiKey postgresql_db.NxaApiKey) (postgresql_db.NxaApiKey, error) {
	queries := getPostgresQueries(p.db)

	apiKey, err := queries.CreateApiKey(ctx, postgresql_db.CreateApiKeyParams{
		Name:      apiKey.Name,
		ApiKey:    apiKey.ApiKey,
		UserID:    apiKey.UserID,
		ProjectID: apiKey.ProjectID,
		ExpiresAt: apiKey.ExpiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	return apiKey, err
}

func (p *Postgres) GetAPIkeyByUserId(ctx context.Context, userID int32) ([]postgresql_db.NxaApiKey, error) {
	queries := getPostgresQueries(p.db)

	apiKey, err := queries.GetAPIkeyByUserId(ctx, userID)
	if err != nil {
		return []postgresql_db.NxaApiKey{}, err
	}

	return []postgresql_db.NxaApiKey{apiKey}, nil
}

func (p *Postgres) GetAPIkeyByApiKey(ctx context.Context, key string) (postgresql_db.NxaApiKey, error) {
	queries := getPostgresQueries(p.db)

	apiKey, err := queries.GetAPIkeyByApiKey(ctx, key)
	if err != nil {
		return postgresql_db.NxaApiKey{}, err
	}

	return apiKey, nil
}

func (p *Postgres) GetProjectUserByProjectIDAndUserID(ctx context.Context, projectID int32, userID int32) (postgresql_db.ProjectUser, error) {
	queries := getPostgresQueries(p.db)

	projectUser, err := queries.GetProjectUserByProjectIDAndUserID(ctx, postgresql_db.GetProjectUserByProjectIDAndUserIDParams{
		ProjectID: projectID,
		UserID:    userID,
	})
	if err != nil {
		return postgresql_db.ProjectUser{}, err
	}

	return projectUser, nil
}

func (p *Postgres) GetAPIKeyByNameAndProjectId(ctx context.Context, name string, projectID int32) (postgresql_db.NxaApiKey, error) {
	queries := getPostgresQueries(p.db)

	apiKey, err := queries.GetAPIKeyByNameAndProjectId(ctx, postgresql_db.GetAPIKeyByNameAndProjectIdParams{
		Name:      name,
		ProjectID: projectID,
	})
	if err != nil {
		return postgresql_db.NxaApiKey{}, err
	}

	return apiKey, nil
}

// func (p *Postgres) GetAllApiKeysByUserId(ctx context.Context, userID int32) ([]postgresql_db.NxaApiKey, error) {
// 	queries := getPostgresQueries(p.db)

// 	apiKeys, err := queries.GetAllApiKeysByUserId(ctx, userID)
// 	if err != nil {
// 		return []postgresql_db.NxaApiKey{}, err
// 	}

// 	return apiKeys, nil
// }

func (p *Postgres) GetAllApiKeysByUserId(ctx context.Context, userID int32) ([]postgresql_db.GetAllApiKeysByUserIdRow, error) {
	queries := getPostgresQueries(p.db)

	apiKeys, err := queries.GetAllApiKeysByUserId(ctx, userID)
	if err != nil {
		return []postgresql_db.GetAllApiKeysByUserIdRow{}, err
	}

	return apiKeys, nil
}

func (p *Postgres) GetProjectsByOrgId(ctx context.Context, orgID int32) ([]postgresql_db.Project, error) {
	queries := getPostgresQueries(p.db)

	projects, err := queries.GetProjects(ctx, orgID)
	if err != nil {
		return []postgresql_db.Project{}, err
	}

	return projects, nil
}
