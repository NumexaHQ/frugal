package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	costcalculation "github.com/NumexaHQ/captainCache/numexa-common/cost-calculation"
	"github.com/NumexaHQ/monger/model"
	vibeModel "github.com/NumexaHQ/vibe/model"
	"github.com/gofiber/fiber/v2"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const keyLength = 32

func (h *Handler) Pong(c *fiber.Ctx) error {
	return c.SendString("pong works")
}

func (h *Handler) GetRequestByUserID(c *fiber.Ctx) error {
	userID := c.Params("userID")
	to, from, err := GetTimeFilter(c.Query("to"), c.Query("from"))
	if err != nil {
		return err
	}

	userIDT, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return err
	}

	var result []model.ProxyRequest
	var res []vibeModel.AllRequestsTableResponse
	if to != 0 && from != 0 {
		_ = h.ChConfig.DB.Table("proxy_requests").Where("user_id = ? AND request_timestamp BETWEEN ? AND ?", int32(userIDT), from, to).Scan(&result)
	} else {
		_ = h.ChConfig.DB.Table("proxy_requests").Where("user_id = ?", int32(userIDT)).Scan(&result)
	}

	for _, v := range result {
		var reqBody map[string]interface{}
		var msgs []interface{}
		err := json.Unmarshal([]byte(v.RequestBody), &reqBody)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			// todo: skip it for now
			continue
		}

		var prompt string
		var cost float64
		var llmModel string
		var statusCode int
		var latency int64
		// todo: repeating code
		var presResults []model.ProxyResponse
		_ = h.ChConfig.DB.Table("proxy_responses").Where("request_id = ?", v.RequestID).Scan(&presResults)

		if len(presResults) != 0 {
			statusCode = int(presResults[len(presResults)-1].ResponseStatusCode)
		}

		msgs, ok := reqBody["messages"].([]interface{})

		if ok {
			// last element of the array is the prompt
			if len(msgs) > 0 {
				prompt = msgs[len(msgs)-1].(map[string]interface{})["content"].(string)
			}

			// var usage map[string]float64
			var responseBody map[string]interface{}

			// fmt.Println("presResults", presResults)
			if len(presResults) != 0 {
				// todo: taking last element, could take head as well!
				presResult := presResults[len(presResults)-1]
				if presResult.ResponseStatusCode == 200 {
					llmModel = reqBody["model"].(string)
					latency = presResult.ResponseTimestamp.Sub(presResult.InitiatedTimestamp).Milliseconds()
					err := json.Unmarshal([]byte(presResult.ResponseBody), &responseBody)
					if err != nil {
						fmt.Println("Error parsing JSON:", err)
					} else {
						usage := responseBody["usage"].(map[string]interface{})
						cost = costcalculation.CalculateOpenAICost(int(usage["prompt_tokens"].(float64)), int(usage["completion_tokens"].(float64)), responseBody["model"].(string))
					}
				}
			}
		}

		res = append(res, vibeModel.AllRequestsTableResponse{
			ID:             v.RequestID,
			ProjectID:      v.ProjectID,
			InitiatedAt:    v.RequestTimestamp,
			Model:          llmModel,
			Prompt:         prompt,
			StatusCode:     statusCode,
			Cost:           cost,
			Latency:        latency,
			Provider:       v.Provider,
			IsCached:       v.IsCached,
			IsCacheHit:     v.IsCacheHit,
			CustomMetaData: v.CustomMetaData,
		})
	}

	return c.JSON(res)
}

func (h *Handler) GetResponseByRequestID(c *fiber.Ctx) error {
	requestID := c.Params("requestID")
	// convert requestID to uint64
	// requestIDT, err := strconv.ParseInt(requestID, 10, 64)
	// if err != nil {
	// 	return err
	// }
	var result []model.ProxyResponse
	_ = h.ChConfig.DB.Table("proxy_responses").Where("request_id = ?", requestID).Scan(&result)
	// Create a variable to hold the parsed JSON data
	var jsonData map[string]interface{}

	// Parse the JSON string into the jsonData variablexw

	if len(result) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "response body is empty",
		})
	} else {
		err := json.Unmarshal([]byte(result[0].ResponseBody), &jsonData)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return err
		}
	}

	return c.Status(fiber.StatusOK).JSON(jsonData)
}

func (h *Handler) GetTotalRequestsCountbyProjectID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(float64)
	projectID := c.Params("projectID")

	to, from, err := GetTimeFilter(c.Query("to"), c.Query("from"))
	if err != nil {
		return err
	}

	projectIDT, err := strconv.ParseInt(projectID, 10, 64)
	if err != nil {
		return err
	}

	var result []model.ProxyRequest
	if to != 0 && from != 0 {
		_ = h.ChConfig.DB.Table("proxy_requests").Where("user_id = ? AND project_id = ? AND request_timestamp BETWEEN ? AND ?", userID, int32(projectIDT), from, to).Scan(&result)
	} else {
		_ = h.ChConfig.DB.Table("proxy_requests").Where("user_id = ? AND project_id = ?", userID, int32(projectIDT)).Scan(&result)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total_requests": len(result),
	})
}

func (h *Handler) ComputeAvgResponseLatencyByProjectID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(float64)
	projectID := c.Params("projectID")

	to, from, err := GetTimeFilter(c.Query("to"), c.Query("from"))
	if err != nil {
		return err
	}

	projectIDT, err := strconv.ParseInt(projectID, 10, 64)
	if err != nil {
		return err
	}

	var responses []model.ProxyResponse
	if to != 0 && from != 0 {
		_ = h.ChConfig.DB.Table("proxy_responses").Where("user_id = ? AND project_id = ? AND response_timestamp BETWEEN ? AND ?", userID, int32(projectIDT), from, to).Scan(&responses)
	} else {
		_ = h.ChConfig.DB.Table("proxy_responses").Where("user_id = ? AND project_id = ? ", userID, int32(projectIDT)).Scan(&responses)
	}

	var totalLatency int64
	var totalOkResponses int64

	for _, response := range responses {
		// if response code was not 200, skip
		if response.ResponseStatusCode == 200 {
			diff := response.ResponseTimestamp.Sub(response.InitiatedTimestamp)
			// diff to milliseconds
			totalLatency += diff.Milliseconds()
			totalOkResponses++
		}
	}

	var avgLatency int64

	if len(responses) == 0 {
		// Handle the case when the responses slice is empty to avoid division by zero.
		// For example, you can set avgLatency to 0 or return an error message.
		avgLatency = 0
	} else {
		avgLatency = totalLatency / totalOkResponses
	}

	// Create a map to include both the response and the average latency
	responseData := map[string]interface{}{
		"avg_latency":     avgLatency,
		"total_responses": len(responses), // Converting to seconds for easy representation.
	}

	return c.Status(fiber.StatusOK).JSON(responseData)
}

func (h *Handler) ComputeLatencyByRequestId(c *fiber.Ctx) error {
	requestID := c.Params("requestID")

	var response model.ProxyResponse
	_ = h.ChConfig.DB.Table("proxy_responses").Where("request_id = ?", requestID).Scan(&response)

	// Calculate the latency
	latency := response.ResponseTimestamp.Sub(response.InitiatedTimestamp).Milliseconds()

	// Create a map to include both the response and the average latency
	responseData := map[string]interface{}{
		"latency": latency, // Converting to seconds for easy representation.
	}

	return c.Status(fiber.StatusOK).JSON(responseData)
}

func (h *Handler) ComputeAverageTokensByProjectID(c *fiber.Ctx) error {
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

	var responses []model.ProxyResponse
	if to != 0 && from != 0 {
		_ = h.ChConfig.DB.Table("proxy_responses").Where("user_id = ? AND project_id = ? AND response_timestamp BETWEEN ? AND ?", userID, int32(projectIDT), from, to).Scan(&responses)
	} else {
		_ = h.ChConfig.DB.Table("proxy_responses").Where("user_id = ? AND project_id = ? ", userID, int32(projectIDT)).Scan(&responses)
	}
	var totalPromptTokens int
	var totalTotalTokens int
	var totalCompletionTokens int
	var totalCost float64
	count := 0
	for _, response := range responses {
		// if response code was not 200, skip
		if response.ResponseStatusCode == 200 {
			var jsonData map[string]interface{}
			err := json.Unmarshal([]byte(response.ResponseBody), &jsonData)
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

				// calculate total cost
				model := jsonData["model"].(string)
				totalCost += costcalculation.CalculateOpenAICost(int(promptTokens), int(completionTokens), model)
			}
			count++
		}
	}

	var avgPromptTokens int
	var avgTotalTokens int
	var avgCompletionTokens int

	if count == 0 {
		avgPromptTokens = 0
		avgTotalTokens = 0
		avgCompletionTokens = 0
	} else {
		avgPromptTokens = totalPromptTokens / count
		avgTotalTokens = totalTotalTokens / count
		avgCompletionTokens = totalCompletionTokens / count
	}

	// Create a map to include both the response and the average latency
	responseData := map[string]interface{}{
		"avg_prompt_tokens":     avgPromptTokens,
		"avg_total_tokens":      avgTotalTokens,
		"avg_completion_tokens": avgCompletionTokens,
		"Total_responses":       len(responses), // Converting to seconds for easy representation.
		"total_success":         count,
		"total_failure":         len(responses) - count,
		"total_cost":            fmt.Sprintf("%.4f", totalCost),
	}

	return c.Status(fiber.StatusOK).JSON(responseData)
}

// funtion to Get count for number of times unique models were used by a user for a project ID (unique models count) by reading responseBody

func (h *Handler) GetUniqueModelsCountByProjectID(c *fiber.Ctx) error {
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

	var responses []model.ProxyResponse
	if to != 0 && from != 0 {
		if err := h.ChConfig.DB.Table("proxy_responses").Where("user_id = ? AND project_id = ? AND response_timestamp BETWEEN ? AND ?", userID, int32(projectIDT), from, to).Find(&responses).Error; err != nil {
			return err
		}
	} else {
		if err := h.ChConfig.DB.Table("proxy_responses").Where("user_id = ? AND project_id = ? ", userID, int32(projectIDT)).Find(&responses).Error; err != nil {
			return err
		}
	}
	// Create a map to count the occurrences of unique model names
	uniqueModelsCount := make(map[string]int)
	for _, response := range responses {
		if response.ResponseStatusCode == 200 {
			var jsonData map[string]interface{}
			err := json.Unmarshal([]byte(response.ResponseBody), &jsonData)
			if err != nil {
				return err
			}
			if model, ok := jsonData["model"]; ok {
				// Get the model name from the response
				modelName := model.(string)
				uniqueModelsCount[modelName]++
			}
		}
	}

	// Create a response map
	var responseData []map[string]interface{}
	for modelName, count := range uniqueModelsCount {
		responseData = append(responseData, map[string]interface{}{
			"name":  modelName,
			"count": count,
		})
	}

	return c.Status(fiber.StatusOK).JSON(responseData)
}

func GetTimeFilter(t, f string) (toUnix, fromUnix int64, err error) {
	if t == "" || f == "" {
		return toUnix, fromUnix, nil
	}

	to, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return toUnix, fromUnix, err
	}
	toUnix = to.Unix()

	from, err := time.Parse(time.RFC3339, f)
	if err != nil {
		return toUnix, fromUnix, err
	}
	fromUnix = from.Unix()

	return toUnix, fromUnix, nil
}
