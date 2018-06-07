package app

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

// Episode is a container for episode related information
type Episode struct {
	Name     string
	Number   int
	Metadata EpisodeMetadata
	URL      *url.URL
	Link     []*DownloadLink
	Season   *Season
}

type EpisodeMetadata struct {
	ProgressPercentage string
	Filesize           string
	Speed              string
	ETA                string
}

// DownloadLink contains the URL and other information to fetch a episode
type DownloadLink struct {
	Rank    int
	Domain  string
	Link    string
	Episode *Episode
}

func (e *Episode) String() string {
	return fmt.Sprintf("S%dE%d %s", e.Season.Number, e.Number, e.Name)
}
func (e *Episode) FetchLinks() error {
	res, err := http.Get(e.URL.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}
	// find seasons
	doc.Find(".link-list .ll-item").Each(func(i int, s *goquery.Selection) {
		rankString := s.Find(".rank").Text()
		domainString := s.Find(".domain").Text()
		linkString, _ := s.Find("a").Attr("href")

		rank, err := strconv.Atoi(rankString)
		if err != nil {
			log.Warnf("could not parse rank string: %s", rankString)
			return
		}

		downloadLink := &DownloadLink{
			Rank:    rank,
			Domain:  domainString,
			Link:    linkString,
			Episode: e,
		}
		e.Link = append(e.Link, downloadLink)
	})

	return nil
}

func (e *Episode) Download() error {
	log.Infof("downloading %s", e.String())

	if len(e.Link) == 0 {
		log.Infof("no links found, skipping")
		return nil
	}

	// todo: increase value or make it configurable?
	maxAttempts := 15
	for _, link := range e.Link {
		<-time.After(time.Second)
		maxAttempts--
		if maxAttempts == 0 {
			return fmt.Errorf("max attempts reached")
		}
		log.Infof("trying link %s", link.Link)
		res, err := http.Get(link.Link)
		if err != nil {
			log.Warnf("could not fetch link page: %s", link.Link)
			continue
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Warnf("status code error: %d %s", res.StatusCode, res.Status)
			continue
		}
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Warn(err)
			continue
		}
		watchButtonLink, _ := doc.Find("a.watch-button").Attr("href")
		log.Infof("found watch button link: %s", watchButtonLink)

		// todo: specify --output to normalize the output filename
		cmd := exec.Command("youtube-dl", "-v", "--newline", watchButtonLink)
		cmd.Env = os.Environ()
		output, err := cmd.StdoutPipe()
		if err != nil {
			log.Warnf("could not create stdout pipe")
			continue
		}
		err = cmd.Start()
		if err != nil {
			log.Warn("could not start youtube-dl")
			continue
		}

		// we want to parse the download progress metadata that is printed to stdout
		scanner := bufio.NewScanner(output)
		go func() {
			for scanner.Scan() {
				meta := EpisodeMetadata{}
				line := scanner.Text()
				fmt.Sscanf(line, "[download] %s of %s at %s ETA %s", &meta.ProgressPercentage, &meta.Filesize, &meta.Speed, &meta.ETA)
				if meta.Filesize != "" {
					e.Metadata = meta
					log.Infof("S%dE%d [%s] [%s]", e.Season.Number, e.Number, meta.ProgressPercentage, meta.Speed)
				}
			}
		}()
		err = cmd.Wait()

		if err != nil {
			log.Warnf("youtube-dl failed for link %s", link.Link)
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					log.Warnf("Exit Status: %d", status.ExitStatus())
				}
			}
			continue
		}
		log.Infof("done downloading %s", e.String())
		break
	}

	return nil
}
