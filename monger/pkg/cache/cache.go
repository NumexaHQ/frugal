package cache

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	commonConstants "github.com/NumexaHQ/captainCache/numexa-common/constants"
	nxAuthDB "github.com/NumexaHQ/captainCache/pkg/db"
	"github.com/NumexaHQ/monger/model"
	nxClickhouse "github.com/NumexaHQ/monger/pkg/db/clickhouse"
	"github.com/sirupsen/logrus"
)

type Cache struct {
	Enabled        bool     `json:"enabled"`
	OriginalPrompt string   `json:"original_prompt"`
	GPTResponse    GPTCache `json:"gpt_response"`
}

type GPTCache struct {
	Prompt string `json:"prompt"`
	Answer string `json:"answer"`
}

func New(prompt string, enabled bool) Cache {
	if enabled {
		gptCache, err := fetchFromGPTCache(prompt)
		if err != nil {
			logrus.Errorf("Error fetching from GPT cache: %v", err)
		}

		return Cache{
			Enabled:        enabled,
			OriginalPrompt: prompt,
			GPTResponse:    gptCache,
		}
	} else {
		return Cache{
			Enabled:        enabled,
			OriginalPrompt: prompt,
		}
	}
}

func (c *Cache) GetOriginalPrompt() string {
	return c.OriginalPrompt
}

func (c *Cache) GetCachedPrompt() string {
	return c.GPTResponse.Prompt
}

func (c *Cache) GetCachedAnswer() string {
	return c.GPTResponse.Answer
}

func (c *Cache) SetCachedAnswer(answer string) error {
	c.GPTResponse.Answer = answer
	return c.storeInCache()
}

// private method. todo: make this public, if needed
func (c *Cache) storeInCache() error {
	// Marshal the cache data to JSON
	cacheJSON, err := json.Marshal(c.GPTResponse)
	if err != nil {
		return err
	}

	// Create a new request to store data in the cache server
	cacheURL := commonConstants.GPTCACHE_URL + "/put"
	cacheReq, err := http.NewRequest(http.MethodPost, cacheURL, bytes.NewBuffer(cacheJSON))
	if err != nil {
		logrus.Errorf("Error creating cache request: %v", err)
		return err
	}

	// Set headers for the cache request
	cacheReq.Header.Set("Accept", "application/json")
	cacheReq.Header.Set("Content-Type", "application/json")
	// Add any additional headers you need here

	// Perform the cache request to store the response
	_, err = http.DefaultClient.Do(cacheReq)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) CacheExists() bool {
	return c.GPTResponse.Prompt != ""
}

func fetchFromGPTCache(prompt string) (GPTCache, error) {
	// Create a map to represent the request data
	requestData := map[string]interface{}{
		"prompt": prompt,
	}

	// Marshal the request data to JSON
	requestJSON, err := json.Marshal(requestData)
	if err != nil {
		return GPTCache{}, err
	}

	// Create a new request to fetch data from the cache server
	cacheURL := commonConstants.GPTCACHE_URL + "/get"
	cacheReq, err := http.NewRequest(http.MethodPost, cacheURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		return GPTCache{}, err
	}

	// Set headers for the cache request
	cacheReq.Header.Set("Content-Type", "application/json")

	// Perform the cache request
	cacheResp, err := http.DefaultClient.Do(cacheReq)
	if err != nil {
		return GPTCache{}, err
	}

	return parseGPTCache(cacheResp)
}

func parseGPTCache(c *http.Response) (GPTCache, error) {
	var cr GPTCache
	err := json.NewDecoder(c.Body).Decode(&cr)
	if err != nil {
		logrus.Errorf("Error decoding cache response: %v", err)
		return GPTCache{}, err
	}
	return cr, nil
}

func (c *Cache) IngestCachedRequest(r *http.Request, rt time.Time, authDB nxAuthDB.DB, url, apiKey string, chConfig nxClickhouse.ClickhouseConfig) {
	pr, err := model.ProxyRequestBuilderForHTTPRequest(r, rt, authDB, url, apiKey)
	if err != nil {
		logrus.Errorf("Error building proxy request: %v", err)
	}

	pr.IsCacheHit = true
	go func() {
		chConfig.ReqC <- &pr
	}()
}
