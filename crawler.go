package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

func (cfg *Config) addPageVisit(normalizedURL string) bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	_, seen := cfg.Pages[normalizedURL]
	if seen {
		return false
	}
	cfg.Pages[normalizedURL] = PageData{}
	return true
}

func (cfg *Config) setPageData(normalizedURL string, data PageData) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.Pages[normalizedURL] = data
}

func (cfg *Config) reachedMaxPage() bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	return len(cfg.Pages) >= cfg.maxPages
}

func (cfg *Config) crawlPage(rawCurrentURL string) {

	cfg.ConcurrencyControl <- struct{}{}
	defer func() {
		<-cfg.ConcurrencyControl
		cfg.Wg.Done()
	}()

	if cfg.reachedMaxPage() {
		return
	}

	if !strings.HasPrefix(rawCurrentURL, cfg.BaseUrl.String()) {
		return
	}
	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Fatalf("error normalizing url: %s\n", err)
	}

	if isFirst := cfg.addPageVisit(normalizedURL); !isFirst {
		return
	}

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}

	fmt.Println("fetched from " + rawCurrentURL)

	urlParsed, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	pageData := extractPageData(html, urlParsed)
	cfg.setPageData(normalizedURL, pageData)

	for _, nextURL := range pageData.OutgoingLinks {
		cfg.Wg.Add(1)
		go cfg.crawlPage(nextURL)
	}
}
