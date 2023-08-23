package metrics

import (
	"encoding/json"
	"fmt"

	"github.com/NumexaHQ/captainCache/types"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var redisClient *redis.Client

func init() {
	// Connect to Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // Enter your Redis password if applicable
		DB:       0,  // Use default Redis database
	})
}

// SaveResponse saves the given response in Redis.
func SaveResponse(resp types.ChatCompletionResponse) error {
	ctx := context.Background()

	// Initialize Redis client for DB 0
	rdb := redisClient

	// Marshal the response data to JSON
	responseJSON, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("error marshaling response JSON: %w", err)
	}

	// Store the response data in Redis Hash, where the field name is the model name
	modelKey := "response:" + resp.Model
	err = rdb.HSet(ctx, modelKey, resp.ID, responseJSON).Err()
	if err != nil {
		return fmt.Errorf("error storing response in Redis: %w", err)
	}

	return nil
}

// CalculateSums calculates the sums of prompt tokens, completion tokens, and total tokens
// for the given model from the responses stored in a different Redis database (DB 1).
func CalculateSums(modelToCalculate string) (int, int, int, error) {
	ctx := context.Background()
	rdb := redisClient

	// Get all the response IDs for the given model from the new database
	responseIDs, err := rdb.HKeys(ctx, "response:"+modelToCalculate).Result()
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error retrieving response IDs from Redis: %w", err)
	}

	var totalPromptTokens, totalCompletionTokens, totalTokens int
	for _, id := range responseIDs {
		// Get the response data for the given ID and model from the new database
		responseJSON, err := rdb.HGet(ctx, "response:"+modelToCalculate, id).Result()
		if err != nil {
			return 0, 0, 0, fmt.Errorf("error retrieving response from Redis: %w", err)
		}

		// Parse the response JSON
		var resp types.ChatCompletionResponse
		err = json.Unmarshal([]byte(responseJSON), &resp)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("error parsing response JSON: %w", err)
		}

		// Calculate sums
		totalPromptTokens += resp.Usage.PromptTokens
		totalCompletionTokens += resp.Usage.CompletionTokens
		totalTokens += resp.Usage.TotalTokens
	}

	return totalPromptTokens, totalCompletionTokens, totalTokens, nil
}

// GetAllResponses retrieves all stored responses from Redis and returns them as a map
// where the key is the response ID and the value is the Response struct.
// GetAllResponses gets all stored responses for a given model from the Redis database.
func GetAllResponses(modelToRetrieve string) ([]types.ChatCompletionResponse, error) {
	ctx := context.Background()

	rdb := redisClient

	// Get all the response IDs for the given model from the Redis database
	responseIDs, err := rdb.HKeys(ctx, "response:"+modelToRetrieve).Result()
	if err != nil {
		return nil, fmt.Errorf("error retrieving response IDs from Redis: %w", err)
	}

	var responses []types.ChatCompletionResponse
	for _, id := range responseIDs {
		// Get the response data for the given ID and model from the Redis database
		responseJSON, err := rdb.HGet(ctx, "response:"+modelToRetrieve, id).Result()
		if err != nil {
			return nil, fmt.Errorf("error retrieving response from Redis: %w", err)
		}

		// Parse the response JSON
		var resp types.ChatCompletionResponse
		err = json.Unmarshal([]byte(responseJSON), &resp)
		if err != nil {
			return nil, fmt.Errorf("error parsing response JSON: %w", err)
		}

		// Append the response to the list of responses
		responses = append(responses, resp)
	}

	return responses, nil
}
