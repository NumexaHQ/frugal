package postgres

import (
	"context"
	"database/sql"
	"time"

	postgresql_db "github.com/NumexaHQ/captainCache/numexa-common/postgresql/postgresql-db"
)

func (p *Postgres) UpdateUserLastLogin(ctx context.Context, user postgresql_db.User) error {
	queries := getPostgresQueries(p.db)

	_, err := queries.UpdateUserLastLogin(ctx, postgresql_db.UpdateUserLastLoginParams{
		ID:        user.ID,
		LastLogin: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return err
	}

	return nil
}
