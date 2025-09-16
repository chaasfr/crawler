package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		fmt.Println("usage: ./crawler URL maxConcurrency maxPages")
		os.Exit(1)
	}
	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	baseUrl, err := url.Parse(args[0])
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	fmt.Println("starting crawl of: " + baseUrl.String())
	
	maxConcurrency, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("error parsing to maxConcurrency to int: %s", err.Error())
	}
	maxPages, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatalf("error parsing to maxPages to int: %s", err.Error())
	}

	cfg, err := configure(baseUrl.String(), maxConcurrency, maxPages)
	if err != nil {
		log.Fatalf("Error - configure: %v", err)
	}

	cfg.Wg.Add(1)
	cfg.crawlPage(baseUrl.String())
	cfg.Wg.Wait()

	fmt.Println(" ==== DONE =====")
	
	reportFileName := "report.csv"

	err = writeCSVReport(cfg.Pages, reportFileName)
	if err != nil {
		log.Fatalf("error writing report: %s", err.Error())
	}
	fmt.Printf("report saved under %s", reportFileName)
}
