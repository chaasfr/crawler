package main

import (
	"fmt"
	"net/url"
	"sync"
)


type Config struct {
	Pages              map[string]PageData
	BaseUrl            *url.URL
	mu                 *sync.Mutex
	ConcurrencyControl chan struct{}
	Wg                 *sync.WaitGroup
}

func configure(rawBaseURL string, maxConcurrency int) (*Config, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	return &Config{
		Pages:              make(map[string]PageData),
		BaseUrl:            baseURL,
		mu:                 &sync.Mutex{},
		ConcurrencyControl: make(chan struct{}, maxConcurrency),
		Wg:                 &sync.WaitGroup{},
	}, nil
}