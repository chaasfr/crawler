package main

import (
	"encoding/csv"
	"os"
)


func writeCSVReport(pages map[string]PageData, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)

	headers := []string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"}
	w.Write(headers)

	for _,p := range pages{
		w.Write(p.toSlice())
	}

	return nil
}