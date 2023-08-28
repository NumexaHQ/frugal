package cache

// Define a struct for sending data to the caching server
type CacheRequest struct {
	Prompt string `json:"prompt"`
	Answer string `json:"answer"`
}

// func serveCachedResponse(w http.ResponseWriter, cachedResp *model.ProxyResponse) {
// 	// Copy all headers from the cached response to the original response
// 	for key, values := range cachedResp.Headers {
// 		for _, value := range values {
// 			w.Header().Add(key, value)
// 		}
// 	}

// 	// Set the status code for the original response to match the cached response
// 	w.WriteHeader(cachedResp.StatusCode)

// 	// Write the cached response body to the original response
// 	w.Write(cachedResp.Body)
// }

// func storeResponseInCache(cacheKey string, cacheValue *model.ProxyResponse) error {
// 	// Prepare the data for sending to the caching server
// 	cacheRequest := CacheRequest{
// 		Prompt: "Your Prompt Data",
// 		Answer: "Your Answer Data",
// 	}

// 	// Marshal the cacheRequest into JSON
// 	requestData, err := json.Marshal(cacheRequest)
// 	if err != nil {
// 		return err
// 	}

// 	// Send a POST request to the caching server to store the response
// 	resp, err := http.Post("http://localhost:8000/put", "application/json", bytes.NewBuffer(requestData))
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("Caching server returned non-OK status code: %d", resp.StatusCode)
// 	}

// 	return nil
// }
