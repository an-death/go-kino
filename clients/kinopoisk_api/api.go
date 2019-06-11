package kinopoisk

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strconv"
	"time"

	"github.com/an-death/go-kino/clients"
)

const (
	_releaseMonthFormat         = "01.2006" // eq "%m.%Y"
	KINOPOISK_API_RELEASES_PATH = "/k/v1/films/releases/digital"
	KINOPOISK_BASE_URL          = "https://ma.kinopoisk.ru"
)

type KinopoiskAPI interface {
	GetReleases(date time.Time) ([]MovieItem, error)
	FilmDetail(filmID int) (MovieItem, error)
}

type kinopoiskAPI struct {
	clients.APIClient
}

func (api *kinopoiskAPI) GetReleases(date time.Time) ([]MovieItem, error) {
	return api.getReleases(date, 0)
}

func (api *kinopoiskAPI) getReleases(date time.Time, offset int) ([]MovieItem, error) {
	uri, err := api.prepareRealeasesUri(date, offset)
	if err != nil {
		return nil, err
	}
	var innerF func([]MovieItem, url.URL) error
	var movies []MovieItem

	innerF = func(movies []MovieItem, uri url.URL) error {
		resp, err := api.APIClient.Request("GET", KINOPOISK_BASE_URL, uri.String())
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		var rc responseContainer
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(buf, &rc)
		if err != nil || !rc.IsSuccess {
			return err
		}

		if movies == nil {
			movies = rc.Data.Items
		} else {
			movies = append(movies, rc.Data.Items...)
		}

		if rc.Data.Stats.Offset == 0 {
			return nil
		}

		newUri, err := api.prepareRealeasesUri(date, rc.Data.Stats.Offset)
		if err != nil {
			return err
		}
		return innerF(movies, *newUri)
	}
	return movies, innerF(movies, *uri)
}

func (api *kinopoiskAPI) prepareRealeasesUri(date time.Time, offset int) (*url.URL, error) {
	uri, err := url.Parse(KINOPOISK_API_RELEASES_PATH)
	if err != nil {
		return uri, err
	}
	q := uri.Query()
	q.Add("digitalReleaseMonth", date.Format(_releaseMonthFormat))
	q.Add("limit", "1000")
	q.Add("offset", strconv.Itoa(offset))
	uri.RawQuery = q.Encode()
	return uri, nil

}
