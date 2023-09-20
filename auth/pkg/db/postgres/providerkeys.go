package postgres

import (
	"context"
	"database/sql"
	"time"

	postgresql_db "github.com/NumexaHQ/captainCache/numexa-common/postgresql/postgresql-db"
	"github.com/sirupsen/logrus"
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

func (p *Postgres) CheckProviderAndNXAKeyPropertyFromNXAKey(ctx context.Context, nxaKey string, providerName string) (bool, postgresql_db.ProviderKey, postgresql_db.NxaApiKeyProperty, []postgresql_db.ProviderSecret, error) {
	notValid := false
	var keyProperty postgresql_db.NxaApiKeyProperty
	queries := getPostgresQueries(p.db)

	nxaKeyRow, err := queries.GetAPIkeyByApiKey(ctx, nxaKey)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("api key not found: %s", nxaKey)
			return notValid, postgresql_db.ProviderKey{}, postgresql_db.NxaApiKeyProperty{}, []postgresql_db.ProviderSecret{}, nil
		}
		return notValid, postgresql_db.ProviderKey{}, postgresql_db.NxaApiKeyProperty{}, []postgresql_db.ProviderSecret{}, err
	}

	if nxaKeyRow.ExpiresAt.Before(time.Now()) {
		return notValid, postgresql_db.ProviderKey{}, postgresql_db.NxaApiKeyProperty{}, []postgresql_db.ProviderSecret{}, nil
	}

	if nxaKeyRow.Revoked {
		return notValid, postgresql_db.ProviderKey{}, postgresql_db.NxaApiKeyProperty{}, []postgresql_db.ProviderSecret{}, nil
	}

	if nxaKeyRow.Disabled {
		return notValid, postgresql_db.ProviderKey{}, postgresql_db.NxaApiKeyProperty{}, []postgresql_db.ProviderSecret{}, nil
	}

	if nxaKeyRow.NxaApiKeyPropertyID.Valid {
		// api key might not have a provider key associated with it
		// and can still be valid
		keyProperty, err = queries.GetNXAKeyPropertyByID(ctx, nxaKeyRow.NxaApiKeyPropertyID.Int32)
		if err != nil {
			if err == sql.ErrNoRows {
				return true, postgresql_db.ProviderKey{}, postgresql_db.NxaApiKeyProperty{}, []postgresql_db.ProviderSecret{}, nil
			}
			return notValid, postgresql_db.ProviderKey{}, postgresql_db.NxaApiKeyProperty{}, []postgresql_db.ProviderSecret{}, err
		}
	}

	// get provider key
	providerKey, err := queries.GetProviderKeyByID(ctx, nxaKeyRow.ProviderKeyID.Int32)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("provider key not found for api key: %s", nxaKey)
			return notValid, postgresql_db.ProviderKey{}, postgresql_db.NxaApiKeyProperty{}, []postgresql_db.ProviderSecret{}, nil
		} else {
			return notValid, postgresql_db.ProviderKey{}, postgresql_db.NxaApiKeyProperty{}, []postgresql_db.ProviderSecret{}, err
		}
	}

	// get provider secrets
	providerSecret, err := queries.GetProviderSecretsByProviderKeyID(ctx, nxaKeyRow.ProviderKeyID.Int32)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("provider secrets not found for provider key: %s", providerKey.Name)
			return notValid, postgresql_db.ProviderKey{}, postgresql_db.NxaApiKeyProperty{}, []postgresql_db.ProviderSecret{}, nil
		} else {
			return notValid, postgresql_db.ProviderKey{}, postgresql_db.NxaApiKeyProperty{}, []postgresql_db.ProviderSecret{}, err
		}
	}

	return true, providerKey, keyProperty, providerSecret, nil
}
