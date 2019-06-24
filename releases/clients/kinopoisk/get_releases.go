package kinopoisk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/an-death/go-kino/releases/clients"
)

const (
	_releaseMonthFormat         = "01.2006" // eq "%m.%Y"
	KINOPOISK_API_RELEASES_PATH = "/k/v1/films/releases/digital"
	KINOPOISK_BASE_URL          = "https://ma.kinopoisk.ru"
)

type apiError string

func (e apiError) Error() string {
	return string(e)
}

const (
	NoReleasesFound apiError = "No releases found"
)

type KinopoiskReleaser interface {
	GetReleases(from, to time.Time) ([]ReleaseItem, error)
}

func NewReleases(do clients.Doer) KinopoiskReleaser {
	client := NewAPIClient(KINOPOISK_BASE_URL, do)
	return &kinopoiskReleaser{client}
}

type kinopoiskReleaser struct {
	clients.APIClient
}

func (api *kinopoiskReleaser) GetReleases(from, to time.Time) ([]ReleaseItem, error) {
	var releases = make([]ReleaseItem, 0, 150)
	log.Printf("Kinopoisk get releases from %s to %s /n", from, to)
	for date := from; date.Before(to) || date.Month() == to.Month(); date = date.AddDate(0, 1, 0) {
		newR, err := api.getReleases(date, 0)
		log.Printf("Movies recived %v", len(newR))
		if err != nil {
			log.Println(err)
			continue
		}
		releases = append(releases, newR...)
		log.Printf("Movies collected %v", len(releases))
	}
	if releases == nil {
		return nil, NoReleasesFound
	}

	return releases, nil
}

func (api *kinopoiskReleaser) getReleases(date time.Time, offset int) ([]ReleaseItem, error) {
	log.Printf("Kinopoisk Request for date %s  with offset %v", date, offset)
	uri := api.prepareRealeasesUri(date, offset)
	var innerF func(url.URL) error
	var movies = make([]ReleaseItem, 0, 150)

	innerF = func(uri url.URL) error {
		resp, err := api.APIClient.Request("GET", uri.String())
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
		if err != nil {
			return err
		}

		if !rc.IsSuccess {
			return fmt.Errorf("get releases request failed. Content %s", buf)
		}

		movies = append(movies, rc.Data.Items...)

		if rc.Data.Stats.Offset == 0 {
			return nil
		}

		newUri := api.prepareRealeasesUri(date, rc.Data.Stats.Offset)
		return innerF(newUri)
	}
	return movies, innerF(uri)
}

func (api *kinopoiskReleaser) prepareRealeasesUri(date time.Time, offset int) url.URL {
	q := url.Values{
		"digitalReleaseMonth": []string{date.Format(_releaseMonthFormat)},
		"limit":               []string{"1000"},
		"offset":              []string{strconv.Itoa(offset)},
	}
	return url.URL{
		Path:     KINOPOISK_API_RELEASES_PATH,
		RawQuery: q.Encode(),
	}

}
