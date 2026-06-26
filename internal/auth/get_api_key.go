package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := strings.TrimPrefix(headers.Get("Authorization"), "ApiKey ")
	if apiKey == "" {
		return "", errors.New("Couldn't find Authorization header")
	}
	return apiKey, nil
}
