package postgres

import (
	"context"

	postgresql_db "github.com/NumexaHQ/captainCache/numexa-common/postgresql/postgresql-db"
)

func (p *Postgres) CreateSetting(ctx context.Context, setting postgresql_db.CreateSettingParams) (postgresql_db.Setting, error) {
	queries := getPostgresQueries(p.db)
	return queries.CreateSetting(ctx, setting)
}

func (p *Postgres) GetSetting(ctx context.Context, key string) (postgresql_db.Setting, error) {
	queries := getPostgresQueries(p.db)
	return queries.GetSetting(ctx, key)
}
