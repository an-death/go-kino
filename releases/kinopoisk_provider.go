package releases

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	kinopoisk "github.com/an-death/go-kino/releases/clients/kinopoisk_api"
)

func NewKinopoiskProvider() ReleaseProvider {
	api_client := kinopoisk.NewAPIClient(&http.Client{})
	client := kinopoisk.NewKinopoiskAPI(api_client)
	return &kinopoiskProvider{client}
}

type kinopoiskProvider struct {
	client kinopoisk.KinopoiskAPI
}

func (p *kinopoiskProvider) GetReleases(from, to time.Time) []Release {
	var movies []kinopoisk.MovieItem
	log.Printf("Kinopoisk Start from %s to %s /n", from, to)
	for date := from; date.Before(to) || date.Month() == to.Month(); date = date.AddDate(0, 1, 0) {
		newM, err := p.client.GetReleases(date)
		log.Printf("Movies recived %v", len(newM))
		if err != nil {
			log.Println(err)
			continue
		}
		movies = append(movies, newM...)
		log.Printf("Movies collected %v", len(movies))
	}
	if movies == nil {
		return nil
	}

	return p.fromMoviesToReleases(movies)
}

func (p kinopoiskProvider) fromMoviesToReleases(movies []kinopoisk.MovieItem) []Release {
	var releases = make([]Release, len(movies), len(movies))
	for i, movie := range movies {
		release := releases[i]
		release.OriginName = movie.OriginTitle
		release.Raiting = movie.Rating.Value
		release.NameRu = movie.Title
		release.PosterUrl = movie.Poster.Url
		release.InfoTable = []Info{
			{Key: "Описание", Val: "NotImplemented"},
			{Key: "Год", Val: strconv.Itoa(movie.Year)},
			{Key: "Страна", Val: fmt.Sprintf("%v", movie.Countries)},
			{Key: "Жанр", Val: fmt.Sprintf("%v", movie.Countries)},
			{Key: "Описание", Val: ""},
			{Key: "Описание", Val: ""},
			{Key: "Описание", Val: ""},
			{Key: "Описание", Val: ""},
			{Key: "Описание", Val: ""},
			{Key: "Описание", Val: ""},
			{Key: "Описание", Val: ""},
			{Key: "Описание", Val: ""},
			{Key: "Описание", Val: ""},
			{Key: "Описание", Val: ""},
		}
		releases[i] = release
	}
	log.Printf("%v", releases)
	return releases
}
