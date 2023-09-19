package openai

import (
	"context"
	"encoding/json"
	"time"

	"github.com/NumexaHQ/captainCache/model"
	"github.com/NumexaHQ/captainCache/numexa-common/encryption"
	postgresql_db "github.com/NumexaHQ/captainCache/numexa-common/postgresql/postgresql-db"
	"github.com/NumexaHQ/captainCache/pkg/constants"
	nxDB "github.com/NumexaHQ/captainCache/pkg/db"
	"github.com/sirupsen/logrus"
)

// todo: call validate on the keys
// Be careful passing this struct around, it contains sensitive, unencrypted information
func New(b []byte, isEncrypted bool) (*ProviderOpenAI, error) {
	var payload Payload
	err := json.Unmarshal(b, &payload)
	if err != nil {
		return nil, err
	}
	return &ProviderOpenAI{
		Payload:   payload,
		encrypted: isEncrypted,
	}, nil
}

func (o *ProviderOpenAI) EncryptKeys(ctx context.Context, db nxDB.DB) error {
	// get the aes key from the setting table
	aesValue, err := model.GetAESSettingValue(ctx, db)
	if err != nil {
		return err
	}

	aes := encryption.AES{}
	err = json.Unmarshal(aesValue, &aes)
	if err != nil {
		return err
	}

	// encrypt the keys
	// and set the encrypted flag to true
	o.Payload.Keys.OpenAIOrg, err = aes.Encrypt(o.Payload.Keys.OpenAIOrg)
	if err != nil {
		return err
	}
	o.Payload.Keys.OpenAIKey, err = aes.Encrypt(o.Payload.Keys.OpenAIKey)
	if err != nil {
		return err
	}
	o.encrypted = true

	return nil

}

func (o *ProviderOpenAI) IsEncrypted() bool {
	return o.encrypted
}

func (o *ProviderOpenAI) GetEncryptedKeys(ctx context.Context, db nxDB.DB) (map[string]string, error) {
	if !o.encrypted {
		err := o.EncryptKeys(ctx, db)
		if err != nil {
			logrus.WithError(err).Error("error encrypting keys")
			return nil, err
		}
	}
	return map[string]string{
		constants.KEY_OPENAI_ORG: o.Payload.Keys.OpenAIOrg,
		constants.KEY_OPENAI_KEY: o.Payload.Keys.OpenAIKey,
	}, nil

}

func (o *ProviderOpenAI) GetDecryptedKeys(ctx context.Context, db nxDB.DB) (map[string]string, error) {
	if o.encrypted {
		err := o.DecryptKeys(ctx, db, o.Payload.Name)
		if err != nil {
			logrus.WithError(err).Error("error decrypting keys")
			return nil, err
		}
	}
	return map[string]string{
		constants.KEY_OPENAI_ORG: o.Payload.Keys.OpenAIOrg,
		constants.KEY_OPENAI_KEY: o.Payload.Keys.OpenAIKey,
	}, nil
}

func (o *ProviderOpenAI) DecryptKeys(ctx context.Context, db nxDB.DB, name string) error {
	// get the aes key from the setting table
	// get the aes key from the setting table
	aesValue, err := model.GetAESSettingValue(ctx, db)
	if err != nil {
		return err
	}

	aes := encryption.AES{}
	err = json.Unmarshal(aesValue, &aes)
	if err != nil {
		return err
	}

	// decrypt the keys
	// and set the encrypted flag to false
	o.Payload.Keys.OpenAIOrg, err = aes.Decrypt(o.Payload.Keys.OpenAIOrg)
	if err != nil {
		return err
	}
	o.Payload.Keys.OpenAIKey, err = aes.Decrypt(o.Payload.Keys.OpenAIKey)
	if err != nil {
		return err
	}

	o.encrypted = false

	return nil
}

func (o *ProviderOpenAI) GetKeys() map[string]string {
	return map[string]string{
		constants.KEY_OPENAI_ORG: o.Payload.Keys.OpenAIOrg,
		constants.KEY_OPENAI_KEY: o.Payload.Keys.OpenAIKey,
	}
}

func (o *ProviderOpenAI) PushKeysToDB(ctx context.Context, db nxDB.DB, name, keyuuid string, userId, projectId, orgId int32) error {
	// create a entry in the provider_keys table
	// with the provider name, user id, project id, and the keys
	// return error if any
	pk, err := db.AddProviderKeys(ctx, postgresql_db.CreateProviderKeyParams{
		Name:           name,
		KeyUuid:        keyuuid,
		Provider:       o.GetProviderName(),
		CreatorID:      userId,
		ProjectID:      projectId,
		OrganizationID: orgId,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		return err
	}

	// create entry in the provider_secrets
	// incase of openai, we have two keys
	// openai_org and openai_key, so we need to loop through the keys
	// and create a entry for each key
	keys, err := o.GetEncryptedKeys(ctx, db)
	if err != nil {
		return err
	}
	for k, v := range keys {
		_, err := db.AddProviderSecrets(ctx, postgresql_db.CreateProviderSecretParams{
			ProviderKeyID: pk.ID,
			Key:           v,
			Type:          k,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *ProviderOpenAI) GetProviderName() string {
	return "openai"
}

func (o *ProviderOpenAI) KeyExists(ctx context.Context, db nxDB.DB, name string) bool {
	_, err := db.GetProviderKeyByName(ctx, name)
	if err != nil {
		logrus.WithError(err).Error("error getting provider key by name")
		return false
	}
	return true
}

// example request body

/*
{
    "name": "openai-prod-key",
	"provider": "openai",
	"keys": {
		"openai_org": "org-5asdf6asdfghasf7as6as8d",
		"openai_key": "key-287asfdvuas66hasd6767agsd7aa"
	}

}

*/
