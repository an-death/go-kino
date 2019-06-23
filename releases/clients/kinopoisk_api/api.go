package kinopoisk

import (
	"encoding/json"
	"errors"
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

type KinopoiskFilmDetail interface {
	FilmDetail(filmID int) (FilmDetail, error)
}

type KinopoiskAPI interface {
	KinopoiskReleaser
	KinopoiskFilmDetail
}

func NewKinopoiskAPI(client clients.APIClient) KinopoiskAPI {
	return &kinopoiskAPI{
		KinopoiskReleaser:   &kinopoiskReleaser{client},
		KinopoiskFilmDetail: &kinopoiskFilmDetail{client},
	}
}

type kinopoiskAPI struct {
	KinopoiskReleaser
	KinopoiskFilmDetail
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
	var movies []ReleaseItem

	innerF = func(uri url.URL) error {
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
		//log.Printf("%s\n\n", buf)
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

const (
	KINOPOISK_BASE_URL_FILMDETAIL = "https://ma.kinopoisk.ru/ios/5.0.0/"
	KINOPOISK_API_FILMDETAIL      = "getKPFilmDetailView"
	POSTER_URL                    = "https://st.kp.yandex.net/images/{}{}width=360"
)

type kinopoiskFilmDetail struct {
	api clients.APIClient
}

func (k *kinopoiskFilmDetail) FilmDetail(filmID int) (FilmDetail, error) {
	var filmDetail FilmDetail
	uri := k.prepareUri(filmID)
	resp, err := k.api.Request("GET", KINOPOISK_BASE_URL_FILMDETAIL, uri.String())
	if err != nil {
		return filmDetail, err
	}
	defer resp.Body.Close()
	var rc filmDetailResponse
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return filmDetail, err
	}
	//log.Printf("%s\n\n", buf)
	err = json.Unmarshal(buf, &rc)
	if err != nil {
		return filmDetail, err
	}

	if !rc.IsSuccess() {
		return filmDetail, errors.New(rc.Message)
	}
	filmDetail = rc.Data

	return filmDetail, nil
}

func (k *kinopoiskFilmDetail) prepareUri(filmID int) url.URL {
	query := url.Values{
		"still_limit": []string{"9"},
		"filmID":      []string{strconv.Itoa(filmID)},
	}
	uri := url.URL{
		Path:     KINOPOISK_API_FILMDETAIL,
		RawQuery: query.Encode(),
	}
	return uri
}
