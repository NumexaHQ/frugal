package postgres

import (
	"context"

	postgresql_db "github.com/NumexaHQ/captainCache/numexa-common/postgresql/postgresql-db"
)

func (p *Postgres) AddProviderKeys(ctx context.Context, pk postgresql_db.CreateProviderKeyParams) (postgresql_db.ProviderKey, error) {
	queries := getPostgresQueries(p.db)
	return queries.CreateProviderKey(ctx, pk)
}

func (p *Postgres) AddProviderSecrets(ctx context.Context, ps postgresql_db.CreateProviderSecretParams) (postgresql_db.ProviderSecret, error) {
	queries := getPostgresQueries(p.db)
	return queries.CreateProviderSecret(ctx, ps)
}

func (p *Postgres) GetProviderKeyByName(ctx context.Context, name string) (postgresql_db.ProviderKey, error) {
	queries := getPostgresQueries(p.db)
	return queries.GetProviderKeyByName(ctx, name)
}

func (p *Postgres) GetProviderKeysByProjectId(ctx context.Context, projectID int32) ([]postgresql_db.ProviderKey, error) {
	queries := getPostgresQueries(p.db)
	return queries.GetProviderKeysByProjectID(ctx, projectID)
}

func (p *Postgres) GetProviderSecretByProviderId(ctx context.Context, id int32) ([]postgresql_db.ProviderSecret, error) {
	queries := getPostgresQueries(p.db)
	return queries.GetProviderSecretsByProviderKeyID(ctx, id)
}

func (p *Postgres) GetProviderKeyById(ctx context.Context, id int32) (postgresql_db.ProviderKey, error) {
	queries := getPostgresQueries(p.db)
	return queries.GetProviderKeyByID(ctx, id)
}
