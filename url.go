package main

import (
	"net/url"
	"strings"
)

func normalizeURL(input string) (string, error) {
	urlReceived, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	cleanUrl := urlReceived.Host + strings.TrimSuffix(urlReceived.Path, "/")
	return cleanUrl, nil
}
