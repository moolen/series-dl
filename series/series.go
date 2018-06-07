package app

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/PuerkitoBio/goquery"
)

var (
	baseURL  = "https://www.watchepisodeseries.com"
	intRegex = regexp.MustCompile("[^0-9]+")
)

// Series is a container for series related information
// it's the top-level type
type Series struct {
	Name    string
	Slug    string
	Seasons []*Season
}

func Fetch(seriesString string, seasonStart, seasonEnd, concurrency int) (*Series, error) {
	series := &Series{
		Slug: seriesString,
	}
	res, err := http.Get(fmt.Sprintf("%s/%s", baseURL, seriesString))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	// series name
	series.Name = parseTitle(doc.Find(".tnContent .main-title").Text())
	// find seasons
	doc.Find(".el-season-buttons a").Each(func(i int, s *goquery.Selection) {
		seasonName := s.Text()
		if seasonName == "Special Episodes" {
			seasonName = "Season 0"
		}
		seasonNumber, err := parseInt(seasonName)
		if err != nil {
			log.Warnf("could not parse season number: %s", seasonName)
			return
		}
		season := &Season{
			Name:   seasonName,
			Number: seasonNumber,
		}
		log.Infof("found season: %s", season.Name)
		series.Seasons = append(series.Seasons, season)
	})
	// find episodes
	doc.Find(".episode-list .el-item").Each(func(i int, s *goquery.Selection) {
		episodeURLString, _ := s.Find("a").Attr("href")
		seasonName := s.Find(".season").Text()
		episodeString := s.Find(".episode").Text()
		name := s.Find(".name").Text()

		season := series.GetSeason(seasonName)
		if season == nil {
			log.Warnf("season %s not assignable", seasonName)
			return
		}
		episodeURL, err := url.Parse(episodeURLString)
		if err != nil {
			log.Warnf("cannot parse episode URL: %s", episodeURLString)
			return
		}
		episodeNum, err := parseInt(episodeString)
		if err != nil {
			log.Warnf("could not parse episode num: %s", episodeString)
			return
		}
		episode := &Episode{
			Name:   name,
			Number: episodeNum,
			URL:    episodeURL,
			Season: season,
		}
		log.Infof("found episode: %s", episode.String())
		season.Episodes = append(season.Episodes, episode)
	})

	for _, seasons := range series.Seasons {
		sort.Sort(seasons)
	}
	sort.Sort(series)

	tasks := make(chan *Episode, 1000)
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			for episode := range tasks {
				err = episode.Download()
				if err != nil {
					log.Warnf("error downloading %s", episode.String())
				}
			}
			wg.Done()
		}()
	}

	for _, season := range series.Seasons {
		if season.Number < seasonStart || season.Number > seasonEnd {
			continue
		}
		for _, episode := range season.Episodes {
			err := episode.FetchLinks()
			if err != nil {
				log.Warnf("error fetching links for %s: %s", episode.String(), err)
				continue
			}
			tasks <- episode
		}
	}

	close(tasks)
	wg.Wait()

	return series, nil
}

func (s Series) GetSeason(name string) *Season {
	for _, s := range s.Seasons {
		if s.Name == name {
			return s
		}
	}
	return nil
}

func (s Series) Len() int {
	return len(s.Seasons)
}
func (s Series) Swap(i, j int) {
	s.Seasons[i], s.Seasons[j] = s.Seasons[j], s.Seasons[i]
}
func (s Series) Less(i, j int) bool {
	return s.Seasons[i].Number < s.Seasons[j].Number
}

func parseTitle(str string) string {
	return strings.TrimRight(str, " Episodes")
}

func parseInt(str string) (int, error) {
	processedString := intRegex.ReplaceAllString(str, "")
	return strconv.Atoi(processedString)
}
