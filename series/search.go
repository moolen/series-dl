package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

type searchResult struct {
	Status int                  `json:"status"`
	Series []searchResultSeries `json:"series"`
}

type searchResultSeries struct {
	HasImage     int    `json:"has_image"`
	OriginalName string `json:"original_name"`
	SeoName      string `json:"seo_name"`
	SeriesID     int    `json:"series_id"`
}

func Search(seriesString string) ([]string, error) {
	var results []string
	URL, err := url.Parse(fmt.Sprintf("%s/home/search", baseURL))
	if err != nil {
		return results, err
	}
	parameters := url.Values{}
	parameters.Add("q", seriesString)
	URL.RawQuery = parameters.Encode()

	log.Debugf("search url: %s", URL.String())
	res, err := http.Get(URL.String())
	if err != nil {
		return results, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return results, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	result := searchResult{}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&result)
	if err != nil {
		return results, err
	}

	if result.Status != 1 {
		return results, fmt.Errorf("search result error status: %d, %#v", result.Status, result)
	}

	for _, item := range result.Series {
		results = append(results, item.SeoName)
	}

	return results, nil
}
