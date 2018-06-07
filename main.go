package main

import (
	"flag"
	"os"

	series "github.com/moolen/series-dl/series"

	log "github.com/sirupsen/logrus"
)

func main() {

	// name of series
	seriesName := flag.String("series", "", "")

	// number of seasons (defaults to all)
	seasonStart := flag.Int("season-start", 1, "")
	seasonEnd := flag.Int("season-end", 99, "")

	// search string
	search := flag.String("search", "", "")

	concurrency := flag.Int("concurrency", 4, "")

	flag.Parse()

	log.SetLevel(log.InfoLevel)

	if *search != "" {
		searchResults, err := series.Search(*search)
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range searchResults {
			log.Info(v)
		}
		os.Exit(0)
	}

	if *seriesName == "" {
		log.Fatalf("missing series")
	}

	_, err := series.Fetch(*seriesName, *seasonStart, *seasonEnd, *concurrency)
	if err != nil {
		log.Fatal(err)
	}
}
