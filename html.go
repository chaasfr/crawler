package main

import (
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	Url             string
	H1              string
	FirstParagraph string
	OutgoingLinks  []string
	ImageUrls      []string
}

func getH1FromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	h1 := doc.Find("h1").First().Text()
	return strings.TrimSpace(h1)
}

func getFirstParagraphFromHTML(html string) string{
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	main := doc.Find("main").Find("p").Text()
	if main != "" {
		return main
	}
	p := doc.Find("p").Text()
	return p
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
	h1 := getH1FromHTML(html)
	firstParagraph := getFirstParagraphFromHTML(html)
	parsedUrls, err := getURLsFromHTML(html, pageURL)
	if err != nil {
		log.Printf("error getting URL: %s\n", err.Error())
		parsedUrls = nil
	}

	parsedImg, err := getImagesFromHTML(html, pageURL)
	if err != nil {
		log.Printf("error getting images: %s", err.Error())
		parsedImg = nil
	}

	return PageData{
		Url: pageURL.String(),
		H1: h1,
		FirstParagraph: firstParagraph,
		OutgoingLinks: parsedUrls,
		ImageUrls: parsedImg,
	}	
}