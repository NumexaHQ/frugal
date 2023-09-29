package model

import (
	"context"
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"

	commonModel "github.com/NumexaHQ/captainCache/numexa-common/model"
	postgresql_db "github.com/NumexaHQ/captainCache/numexa-common/postgresql/postgresql-db"
	"github.com/NumexaHQ/captainCache/pkg/constants"
	nxDB "github.com/NumexaHQ/captainCache/pkg/db"
)

type SettingValue struct {
	Label       string      `json:"label"`
	Value       interface{} `json:"value"`
	Description string      `json:"description"`
}

func InitializeAESSetting(ctx context.Context, db nxDB.DB) error {
	// set aes_secret in setting table, if !exists
	_, err := db.GetSetting(ctx, constants.AES_SECRET)
	if err != nil {
		key := make([]byte, 32) // 32 bytes for AES-256
		iv := make([]byte, aes.BlockSize)
		_, err = rand.Read(key)
		if err != nil {
			return err
		}

		_, err = rand.Read(iv)
		if err != nil {
			return err
		}

		aesValue := &SettingValue{
			Label:       "AES Encryption Setting",
			Description: "AES Encryption Key-IV pair",
			Value: map[string]string{
				"aes_iv":  hex.EncodeToString(iv),
				"aes_key": hex.EncodeToString(key),
			},
		}

		rawAES, err := json.Marshal(aesValue)
		if err != nil {
			return err
		}
		rawMessageAES := json.RawMessage(rawAES)

		_, err = db.CreateSetting(ctx, postgresql_db.CreateSettingParams{
			Key:     constants.AES_SECRET,
			Value:   rawMessageAES,
			Visible: false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func GetAESSettingValue(ctx context.Context, db nxDB.DB) (json.RawMessage, error) {
	setting, err := db.GetSetting(ctx, constants.AES_SECRET)
	if err != nil {
		return nil, err
	}

	var aesValue SettingValue
	err = json.Unmarshal(setting.Value, &aesValue)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(aesValue.Value)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(b), nil
}

func InitializeUsageLimitSetting(ctx context.Context, db nxDB.DB) error {
	// set usage_limit in setting table, if !exists
	_, err := db.GetSetting(ctx, constants.USAGE_LIMIT)
	if err != nil {
		l := []commonModel.PlanRequestLimit{
			{
				PlanID: "free",
				Name:   "Free",
				Limit:  50000, // 50k
			},
			{
				PlanID: "trial",
				Name:   "Trial",
				Limit:  100000, // 100k
			},
			{
				PlanID: "promo",
				Name:   "Promo",
				Limit:  50000, // 50k, because promo accounts have free provider credits
			},
			{
				PlanID: "paid",
				Name:   "Paid",
				Limit:  1000000, // 1M
			},
		}

		val := commonModel.UsageLimitSetting{
			RequestLimit: l,
		}

		usageLimitValue := &SettingValue{
			Label:       "Numexa Usage Limit Setting",
			Description: "Usage Limit",
			Value:       val,
		}

		rawUsageLimit, err := json.Marshal(usageLimitValue)
		if err != nil {
			return err
		}
		rawMessageUsageLimit := json.RawMessage(rawUsageLimit)

		_, err = db.CreateSetting(ctx, postgresql_db.CreateSettingParams{
			Key:     constants.USAGE_LIMIT,
			Value:   rawMessageUsageLimit,
			Visible: false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
