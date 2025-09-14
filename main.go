package main

import (
	"fmt"
	"log"
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
	baseUrl := args[0]
	fmt.Println("starting crawl of: " + baseUrl)

	html, err := getHTML(baseUrl)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	fmt.Println(html)
}
