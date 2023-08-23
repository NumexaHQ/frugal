package utils

import (
	"fmt"
	"strings"

	"github.com/pkoukk/tiktoken-go"
	log "github.com/sirupsen/logrus"
)

func GetBPETokensByModel(s string, model string) (tokens []int, err error) {
	encoding := "r50k_base"
	if strings.Contains(model, "gpt-3.5-turbo") ||
		strings.Contains(model, "gpt-4") ||
		strings.Contains(model, "text-embedding-ada-002") {
		encoding = "cl100k_base"
	} else if strings.Contains(model, "text-davinci-002") ||
		strings.Contains(model, "text-davinci-003") {
		encoding = "p50k_base"
	} else {
		log.Warnf("model %s not found, using default encoding", model)
	}

	// if you don't want download dictionary at runtime, you can use offline loader
	// tiktoken.SetBpeLoader(tiktoken_loader.NewOfflineLoader())
	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		err = fmt.Errorf("getEncoding: %v", err)
		return nil, err
	}

	// encode
	token := tke.Encode(s, nil, nil)

	return token, nil
}

func GetBPETokenSizeByModel(s string, model string) (size int, err error) {
	tokenSize, err := GetBPETokensByModel(s, model)
	return len(tokenSize), err
}
