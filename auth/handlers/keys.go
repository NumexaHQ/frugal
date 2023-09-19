package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/NumexaHQ/captainCache/model"
	"github.com/NumexaHQ/captainCache/pkg/providerkeys"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (h *Handler) AddProviderKeys(c *fiber.Ctx) error {
	var reqBody model.ProviderKeys
	if err := c.BodyParser(&reqBody); err != nil {
		logrus.WithError(err).Error("Error parsing request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	userId := c.Locals("user_id").(float64)
	orgId := c.Locals("organization_id").(float64)

	rByte, err := json.Marshal(reqBody)
	if err != nil {
		logrus.WithError(err).Error("Error marshalling request body")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	keyProvider, err := providerkeys.GetProvider(reqBody.Provider, rByte, false)
	if err != nil {
		logrus.WithError(err).Error("Error getting provider")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid provider",
		})
	}

	if keyProvider.KeyExists(c.Context(), h.DB, reqBody.Name) {
		logrus.WithError(err).Error("Key already exists")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Key already exists with this name",
		})
	}

	// generate uuid for key
	keyuuid, err := generateUUID()
	if err != nil {
		logrus.WithError(err).Error("Error generating uuid")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	err = keyProvider.PushKeysToDB(c.Context(), h.DB, reqBody.Name, keyuuid, int32(userId), reqBody.ProjectId, int32(orgId))
	if err != nil {
		logrus.WithError(err).Error("Error pushing keys to db")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Keys added successfully",
	})
}

func (h *Handler) GetProviderKeys(c *fiber.Ctx) error {
	// userId := c.Locals("user_id").(float64)
	// orgId := c.Locals("organization_id").(float64)
	projectId := c.Params("project_id")
	projectIdInt, err := strconv.Atoi(projectId)
	if err != nil {
		logrus.WithError(err).Error("Error converting project id to int")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid project id",
		})
	}

	keys, err := h.DB.GetProviderKeysByProjectId(c.Context(), int32(projectIdInt))
	if err != nil {
		logrus.WithError(err).Error("Error getting provider keys")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	resp := []model.ProviderKeys{}

	for _, key := range keys {
		secrets, err := h.DB.GetProviderSecretByProviderId(c.Context(), key.ID)
		if err != nil {
			if err.Error() == "no rows in result set" {
				logrus.WithError(err).Error("Error getting provider secret")
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Internal server error",
				})
			}
			logrus.WithError(err).Error("Error getting provider secret")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		secretsMap := make(map[string]string)
		for _, sv := range secrets {
			secretsMap[sv.Type] = sv.Key
		}

		// here kp.Keys is encrypted keys
		kp := model.ProviderKeys{
			Name:      key.Name,
			Provider:  key.Provider,
			Keys:      secretsMap,
			ProjectId: key.ProjectID,
		}

		kpB, err := json.Marshal(kp)
		if err != nil {
			logrus.WithError(err).Error("Error marshalling provider keys")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		provider, err := providerkeys.GetProvider(key.Provider, kpB, true)
		if err != nil {
			logrus.WithError(err).Error("Error getting provider")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		decryptedKeys, err := provider.GetDecryptedKeys(c.Context(), h.DB)
		if err != nil {
			logrus.WithError(err).Error("Error getting decrypted keys")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		// updating kp.Keys with decrypted keys
		kp.Keys = decryptedKeys

		resp = append(resp, kp)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func generateUUID() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		logrus.WithError(err).Error("Error generating uuid")
		return "", err
	}
	return uuid.String(), nil
}
