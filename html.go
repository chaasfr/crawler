package main

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getH1FromHTML(html string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	result := []string{}
	doc.Find("h1").Each(func (i int, s *goquery.Selection) {
		result = append(result, s.Text())
	})
	return result, nil
}

func getFirstParagraphFromHTML(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}
	main := doc.Find("main").Find("p").Text()
	if main != "" {
		return main, nil
	}
	p := doc.Find("p").Text()
	return p, nil
}


func getURLsFromHTML(html string, baseURL *url.URL) ([]string, error) {
	doc, err:= goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	result := []string{}
	doc.Find("a[href]").Each(func (i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			if strings.HasPrefix(href, "/") {
				result = append(result, baseURL.String() + href)
			} else {
				result = append(result, href)
			}
		}
	})
	return result, nil
}

func getImagesFromHTML(html string, baseURL *url.URL) ([]string, error) {
	doc, err:= goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	result := []string{}
	doc.Find("img[src]").Each(func (i int, s *goquery.Selection) {
		href, ok := s.Attr("src")
		if ok {
			if strings.HasPrefix(href, "/") {
				result = append(result, baseURL.String() + href)
			} else {
				result = append(result, href)
			}
		}
	})
	return result, nil
}