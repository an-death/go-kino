package releases

import (
	"log"
	"net/http"
	"sync"
	"time"

	kinopoisk "github.com/an-death/go-kino/releases/clients/kinopoisk_api"
)

func NewKinopoiskProvider() ReleaseProvider {
	api_client := kinopoisk.NewAPIClient(&http.Client{})
	client := kinopoisk.NewKinopoiskAPI(api_client)
	return &kinopoiskProvider{client}
}

type kinopoiskProvider struct {
	kinopoisk kinopoisk.KinopoiskAPI
}

func (p *kinopoiskProvider) GetReleases(from, to time.Time) []Release {
	var releases []kinopoisk.ReleaseItem
	releases, err := p.kinopoisk.GetReleases(from, to)
	if err != nil {
		log.Printf("%s\n", err)
		return nil
	}
	return p.fillReleases(releases)
}

func (p *kinopoiskProvider) fillReleases(movies []kinopoisk.ReleaseItem) []Release {
	var stack = newReleasesStack(len(movies))
	for i, releaseItem := range movies {
		go func(id int, r kinopoisk.ReleaseItem) {
			result, err := p.getReleaseInfo(releaseItem)
			if err != nil {
				log.Printf("release item %v detail request failed with %s\n", releaseItem, err)
				return
			}
			stack.Add(result)
		}(i, releaseItem)
	}
	return stack.releases
}

func (p *kinopoiskProvider) getReleaseInfo(item kinopoisk.ReleaseItem) (Release, error) {
	return Release{}, nil
}

type releasesStack struct {
	releases []Release
	sync.Mutex
}

func newReleasesStack(size int) *releasesStack {
	return &releasesStack{releases: make([]Release, 0, size)}
}

func (s *releasesStack) Add(release Release) {
	s.Lock()
	s.releases = append(s.releases, release)
	s.Unlock()
}
