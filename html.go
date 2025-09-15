package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	Url            string
	H1             string
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
		if !ok {
			return
		}
		href = strings.TrimSpace(href)
		if href == "" {
			return
		}

		u, err := url.Parse(href)
		if err != nil {
			fmt.Printf("couldn't parse href %q: %v\n", href, err)
			return
		}

		resolved := baseURL.ResolveReference(u)
		result = append(result, resolved.String())
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
		src, ok := s.Attr("src")
		if !ok || strings.TrimSpace(src) == "" {
			return
		}

		u, err := url.Parse(src)
		if err != nil {
			fmt.Printf("couldn't parse src %q: %v\n", src, err)
			return
		}

		absolute := baseURL.ResolveReference(u)
		result = append(result, absolute.String())
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

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		fmt.Printf("error creating req to %s: %s\n", rawURL, err.Error())
		return "", err
	}
	req.Header.Add("User-Agent", "BootCrawler/1.0")

	res, err := http.DefaultClient.Do(req)
	if res.StatusCode > 399 {
		fmt.Printf("req to %s returned an error: %v - %s\n", rawURL, res.StatusCode, res.Status)
		return "", fmt.Errorf("req error:  %v - %s", res.StatusCode, res.Status)
	}
	contentType := res.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		fmt.Printf("content from %s is not html\n", rawURL)
		return "", fmt.Errorf("content from %s is not html", rawURL)
	}
	if err != nil {
		fmt.Printf("error executing req to %s - %s\n", rawURL, err.Error())
		return "", err
	}

	resBody, err := io.ReadAll(res.Body)
 	if err != nil {
		fmt.Printf("error parsing req to %s: %s\n", rawURL, err.Error())
		return "", err
	}

	return string(resBody), nil
}