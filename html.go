package main

import (
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	Url             string
	H1              []string
	FirstParagraph string
	OutgoingLinks  []string
	ImageUrls      []string
}

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

func extractPageData(html string, pageURL *url.URL) PageData{
	parsedUrls, err := getURLsFromHTML(html, pageURL)
	if err != nil {
		log.Fatalf("error getting URL: %s", err.Error())
	}

	h1, err := getH1FromHTML(html)
	if err != nil {
		log.Fatalf("error getting H1: %s", err.Error())
	}

	firstParagraph, err := getFirstParagraphFromHTML(html)
	if err != nil {
		log.Fatalf("error getting first paragraph: %s", err.Error())
	}

	parsedImg, err := getImagesFromHTML(html, pageURL)
	if err != nil {
		log.Fatalf("error getting images: %s", err.Error())
	}

	return PageData{
		Url: pageURL.String(),
		H1: h1,
		FirstParagraph: firstParagraph,
		OutgoingLinks: parsedUrls,
		ImageUrls: parsedImg,
	}
	
}