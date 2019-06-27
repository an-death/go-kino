package releases

import (
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/an-death/go-kino/providers/torrents"

	"github.com/an-death/go-kino/providers/kinopoisk"
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
			tors, err := torrents.GetTorrents(r.Id)
			if err != nil {
				log.Printf("Cannot get torrents for %v err - %s \n", r.Id, err)
			}
			tors = torrents.UniqueByQualitySeeds(tors)
			result.Torrents = make([]Torrent, 0, len(tors))
			for _, tor := range tors {
				result.Torrents = append(result.Torrents, Torrent{Link: tor.Torrent, Type: tor.Quality})
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
