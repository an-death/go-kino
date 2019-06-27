package torrents

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type seeds int
type leaches int
type Torrent struct {
	Title   string
	Quality string
	Date    time.Time
	Size    string
	Url     string
	Torrent string
	Seeds   seeds
	Leaches leaches
}

const (
	RUTOR_SEARCH_URL = "http://rutor.info/search/0/0/010/0/film%20"
)

// LINK: https://regex101.com/r/fGdUBo/3
var qualityRE = regexp.MustCompile(`(?:\d\) )(\w+ \d{3,4}p|\w+)`)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

func GetTorrents(do Doer, filmID int) ([]Torrent, error) {
	req := newRequest(filmID)
	resp, err := do.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("no torrents. response code %v", resp.StatusCode)
	}
	torrents, err := parseTorrents(resp.Body)
	if err != nil {
		return nil, err
	}

	return torrents, nil
}

func newRequest(filmID int) *http.Request {
	r, err := http.NewRequest("GET", RUTOR_SEARCH_URL+strconv.Itoa(filmID), &bytes.Buffer{})
	if err != nil {
		panic(err)
	}

	return r
}

func parseTorrents(r io.Reader) ([]Torrent, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	var torrents = make([]Torrent, 0, 20)
	doc.Find("#index").Find("tr:not(.backgr)").Each(func(i int, tr *goquery.Selection) {
		torrents = append(torrents, parse(tr))
	})
	return torrents, nil
}

func parse(tr *goquery.Selection) Torrent {
	tds := tr.Find("td")
	a := tds.Find("a")
	link, _ := a.Attr("href")
	seed, leaches := extractSeedLeaches(tds.Last().Text())
	name := strings.TrimSpace(a.Text())

	return Torrent{
		Title:   name,
		Torrent: link,
		Quality: extractQuality(name),
		Seeds:   seed,
		Leaches: leaches,
	}
}

func extractQuality(s string) string {
	matches := qualityRE.FindStringSubmatch(s)
	if matches == nil || len(matches) < 2 {
		return s
	}
	return matches[1]
}

func extractSeedLeaches(s string) (seeds, leaches) {
	var sl = []rune(s)
	return seeds(int(sl[1] - '0')), leaches(int(sl[4] - '0'))
}

func UniqueByQualitySeeds(ts []Torrent) []Torrent {
	var out = make([]Torrent, 0, len(ts))
	var set = make(map[string]int)
	for _, t := range ts {
		it, ok := set[t.Quality]
		if !ok {
			it = len(out)
			out = append(out, t)
			set[t.Quality] = it
			continue
		}

		if t.Seeds > out[it].Seeds {
			out[it] = t
		}
	}
	return out
}
