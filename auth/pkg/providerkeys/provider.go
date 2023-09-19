package providerkeys

import (
	"context"
	"fmt"

	"github.com/NumexaHQ/captainCache/pkg/constants"
	nxDB "github.com/NumexaHQ/captainCache/pkg/db"
	"github.com/NumexaHQ/captainCache/pkg/providerkeys/openai"
)

func GetProvider(provider string, b []byte, isEncrypted bool) (Provider, error) {
	switch provider {
	case constants.PROVIDER_OPENAI:
		return openai.New(b, isEncrypted)
	default:
		return nil, fmt.Errorf("invalid provider: %s", provider)
	}
}

type Provider interface {
	GetKeys() map[string]string
	IsEncrypted() bool
	GetEncryptedKeys(ctx context.Context, db nxDB.DB) (map[string]string, error)
	GetDecryptedKeys(ctx context.Context, db nxDB.DB) (map[string]string, error)
	PushKeysToDB(ctx context.Context, db nxDB.DB, name, keyuuid string, userId, projectId, orgId int32) error
	KeyExists(ctx context.Context, db nxDB.DB, name string) bool
}
