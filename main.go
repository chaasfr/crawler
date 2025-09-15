package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	baseUrl, err := url.Parse(args[0])
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	fmt.Println("starting crawl of: " + baseUrl.String())
	
	maxConcurrency := 5

	cfg, err := configure(baseUrl.String(), maxConcurrency)
	if err != nil {
		fmt.Printf("Error - configure: %v", err)
		return
	}

	cfg.Wg.Add(1)
	cfg.crawlPage(baseUrl.String())
	cfg.Wg.Wait()

	fmt.Println(" ==== DONE =====")
	for key, _ := range cfg.Pages {
		fmt.Printf("crawled %s \n", key)
	}
}
