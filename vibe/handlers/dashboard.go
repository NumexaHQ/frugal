package handlers

import (
	"encoding/json"
	"strconv"

	costcalculation "github.com/NumexaHQ/captainCache/numexa-common/cost-calculation"
	"github.com/NumexaHQ/monger/model"
	vibeModel "github.com/NumexaHQ/vibe/model"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetUserRequestsStatsByProjectID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(float64)
	projectID := c.Params("projectID")

	projectIDT, err := strconv.ParseInt(projectID, 10, 64)
	if err != nil {
		return err
	}

	to, from, err := GetTimeFilter(c.Query("to"), c.Query("from"))
	if err != nil {
		return err
	}

	tempStats := make(map[int32]vibeModel.UserUsageStats)

	// get all requests for this project
	var results []model.ProxyResponse
	if to != 0 && from != 0 {
		h.ChConfig.DB.Table("proxy_responses").Where("user_id = ? AND project_id = ? AND response_timestamp BETWEEN ? AND ?", userID, int32(projectIDT), from, to).Scan(&results)
	} else {
		h.ChConfig.DB.Table("proxy_responses").Where("user_id = ? AND project_id = ?", userID, int32(projectIDT)).Scan(&results)
	}

	var totalPromptTokens int
	var totalTotalTokens int
	var totalCompletionTokens int
	for _, result := range results {
		if result.ResponseStatusCode == 200 {
			var jsonData map[string]interface{}
			err := json.Unmarshal([]byte(result.ResponseBody), &jsonData)
			if err != nil {
				return err
			}
			if usage, ok := jsonData["usage"]; ok {
				// get the prompt tokens from the response
				promptTokens := usage.(map[string]interface{})["prompt_tokens"].(float64)
				totalPromptTokens += int(promptTokens)

				// get the total tokens from the response
				totalTokens := usage.(map[string]interface{})["total_tokens"].(float64)
				totalTotalTokens += int(totalTokens)

				// get the completion tokens from the response
				completionTokens := usage.(map[string]interface{})["completion_tokens"].(float64)
				totalCompletionTokens += int(completionTokens)

				// calculate cost
				model := jsonData["model"].(string)
				cost := costcalculation.CalculateOpenAICost(int(promptTokens), int(completionTokens), model)
				if _, ok := tempStats[result.UserID]; !ok {
					tempStats[result.UserID] = vibeModel.UserUsageStats{}
				}
				tempStats[result.UserID] = vibeModel.UserUsageStats{
					TotalCost:    tempStats[result.UserID].TotalCost + cost,
					TotalRequest: tempStats[result.UserID].TotalRequest,
				}
			}
			tempStats[result.UserID] = vibeModel.UserUsageStats{
				TotalCost:    tempStats[result.UserID].TotalCost,
				TotalRequest: tempStats[result.UserID].TotalRequest + 1,
			}
		}
	}

	// get all users for this project
	projectUsers, err := h.AuthDB.GetUsersByProjectId(c.Context(), int32(projectIDT))
	if err != nil {
		logrus.Errorf("Error getting project users: %v", err)
		c.Status(500).JSON(fiber.Map{
			"error": "Error getting project users",
		})
	}

	// populate the response
	response := make(vibeModel.GetUserRequestsStatsByProjectID)
	for _, projectUser := range projectUsers {
		if _, ok := tempStats[projectUser.ID]; ok {
			response[projectUser.Email] = tempStats[projectUser.ID]
		} else {
			response[projectUser.Email] = vibeModel.UserUsageStats{}
		}
	}

	var responseToSend []vibeModel.GetUserRequestsStatsByProjectIDResponse

	for email, stats := range response {
		responseToSend = append(responseToSend, vibeModel.GetUserRequestsStatsByProjectIDResponse{
			Email:        email,
			TotalRequest: stats.TotalRequest,
			TotalCost:    stats.TotalCost,
		})
	}

	return c.JSON(responseToSend)
}
