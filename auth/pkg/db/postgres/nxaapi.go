package postgres

import (
	"context"

	postgresql_db "github.com/NumexaHQ/captainCache/numexa-common/postgresql/postgresql-db"
)

func (p *Postgres) CreateNXAKeyProperty(ctx context.Context, nxaKeyProperty postgresql_db.CreateNXAKeyPropertyParams) (postgresql_db.NxaApiKeyProperty, error) {
	queries := getPostgresQueries(p.db)
	return queries.CreateNXAKeyProperty(ctx, nxaKeyProperty)
}
