package main

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func normalizeURL(input string) (string, error) {
	urlReceived, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	cleanUrl := urlReceived.Host + strings.TrimSuffix(urlReceived.Path, "/") 
	return cleanUrl, nil
}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	htmlTree, err:= html.Parse(reader)
	if err != nil {
		return nil, err
	}
	var result []string
	for n := range htmlTree.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					result = append(result, a.Val)
					break
				}
			}
		}
	}
	return result, nil
}