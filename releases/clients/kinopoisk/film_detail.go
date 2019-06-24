package kinopoisk

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"strconv"

	"github.com/an-death/go-kino/releases/clients"
)

const (
	KINOPOISK_BASE_URL_FILMDETAIL = "https://ma.kinopoisk.ru/ios/5.0.0/"
	KINOPOISK_API_FILMDETAIL      = "getKPFilmDetailView"
)

type KinopoiskFilmDetail interface {
	FilmDetail(filmID int) (FilmDetail, error)
}

func NewFilmDetail(do clients.Doer) KinopoiskFilmDetail {
	client := NewAPIClient(KINOPOISK_BASE_URL_FILMDETAIL, do)
	return &kinopoiskFilmDetail{client}
}

type kinopoiskFilmDetail struct {
	api clients.APIClient
}

func (k *kinopoiskFilmDetail) FilmDetail(filmID int) (FilmDetail, error) {
	var filmDetail FilmDetail
	uri := k.prepareUri(filmID)
	resp, err := k.api.Request("GET", uri.String())
	if err != nil {
		return filmDetail, err
	}
	defer resp.Body.Close()
	var rc filmDetailResponse
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return filmDetail, err
	}
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
