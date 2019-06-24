package releases

import (
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	kinopoisk "github.com/an-death/go-kino/releases/clients/kinopoisk"
)

func NewKinopoiskProvider() ReleaseProvider {
	client := http.Client{}
	return &kinopoiskProvider{
		KinopoiskReleaser:   kinopoisk.NewReleases(&client),
		KinopoiskFilmDetail: kinopoisk.NewFilmDetail(&client),
	}
}

type kinopoiskProvider struct {
	kinopoisk.KinopoiskReleaser
	kinopoisk.KinopoiskFilmDetail
}

func (p *kinopoiskProvider) GetReleases(from, to time.Time) []Release {
	releases, err := p.KinopoiskReleaser.GetReleases(from, to)
	if err != nil {
		log.Printf("%s\n", err)
		return nil
	}
	return p.fillReleases(releases)
}

func (p *kinopoiskProvider) fillReleases(movies []kinopoisk.ReleaseItem) []Release {
	var stack = newReleasesStack(len(movies))
	var group sync.WaitGroup

	defer func(now time.Time) {
		log.Printf("FilmDetail done %s", time.Now().Sub(now))
	}(time.Now())

	group.Add(len(movies))
	for i, releaseItem := range movies {
		go func(id int, r kinopoisk.ReleaseItem) {
			defer group.Done()
			result, err := p.getReleaseInfo(r)
			if err != nil {
				log.Printf("release item Id:%v, %v detail request failed with %s\n", id, r, err)
				return
			}
			stack.Add(result)
		}(i, releaseItem)
	}
	group.Wait()
	sort.Sort(stack.releases)
	return stack.releases
}

func (p *kinopoiskProvider) getReleaseInfo(item kinopoisk.ReleaseItem) (Release, error) {
	info, err := p.KinopoiskFilmDetail.FilmDetail(item.Id)
	if err != nil {
		return Release{}, err
	}
	release := Release{
		OriginName: info.NameEn,
		NameRu:     info.NameRu,
		InfoTable: map[string]string{
			"Год":               info.Year,
			"Страна":            info.Country,
			"Режисёр":           info.Creators.Directors.String(),
			"Актёры":            info.Creators.Actors.String(),
			"Жанр":              info.Genre,
			"Возраст":           info.RaitingAgeLimit,
			"Продолжительность": info.FilmLength,
			"Рeйтинг КиноПоиск": info.RatingData.Rating,
			"Рейтинг IMDb":      info.RatingData.RatingIMDb,
			"Описание":          info.Description,
		},
		PosterUrl: item.Poster.Url,
		Rating:    info.Rating(),
		Date:      item.ContextData.ReleaseDate.AsDate(),
		WebURL:    info.WebURL,
	}
	return release, nil
}

type releasesStack struct {
	releases Releases
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
