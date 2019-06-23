package kinopoisk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type responseContainer struct {
	IsSuccess bool         `json:"success"`
	Data      responseData `json:"data"`
}
type responseData struct {
	Items []ReleaseItem    `json:"items"`
	Stats responseDataStat `json:"stats"`
}
type responseDataStat struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ReleaseItem struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Poster struct {
		Url string `json:"url"`
	} `json:"poster"`
	ContextData struct {
		ReleaseDate releaseDate `json:"releaseDate"`
	} `json:"contextData"`
}

type releaseDate time.Time

func (d releaseDate) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	return []byte(`"` + t.Format("2006-01-02") + `"`), nil
}

func (d *releaseDate) UnmarshalJSON(data []byte) error {
	t, err := time.Parse("2006-01-02", strings.Trim(string(data), "\""))
	if err != nil {
		return err
	}
	*d = releaseDate(t)
	return nil
}

func (d releaseDate) AsDate() time.Time {
	return time.Time(d)
}

type filmDetailResponse struct {
	ResultCode int        `json:"resultCode"`
	Message    string     `json:"message"`
	Data       FilmDetail `json:"data"`
}

func (r *filmDetailResponse) IsSuccess() bool {
	return r.ResultCode == 0
}

type FilmDetail struct {
	RatingData struct {
		Rating     string `json:"rating"`
		RatingIMDb string `json:"ratingIMDb"`
	} `json:"ratingData"`
	HasReleasedFilm int               `json:"hasReleasedFilm"`
	WebURL          string            `json:"webURL"`
	NameRu          string            `json:"nameRu"`
	NameEn          string            `json:"nameEn"`
	PosterURL       string            `json:"posterURL"`
	BigPosterURL    string            `json:"bigPosterUrl"`
	Year            string            `json:"year"`
	FilmLength      string            `json:"filmLength"`
	Country         string            `json:"country"`
	Genre           string            `json:"genre"`
	Description     string            `json:"description"`
	RaitingMPAA     string            `json:"raitingMPAA"`
	RaitingAgeLimit string            `json:"raitingAgeLimit"`
	BudgetData      map[string]string `json:"budgetData"`
	VideoURL        map[string]string `json:"videoURL"`
	Creators        json_creators     `json:"creators"`
}

func (f *FilmDetail) Rating() float64 {
	ratingKP, err := strconv.ParseFloat(f.RatingData.Rating, 32)
	if err != nil {
		return 0
	}
	ratingIMDbFloat, err := strconv.ParseFloat(f.RatingData.RatingIMDb, 32)
	if err != nil {
		return ratingKP
	}
	return (ratingKP + ratingIMDbFloat) / 2
}

type json_creators struct {
	Directors peoples
	Actors    peoples
	Producers peoples
}

func (j *json_creators) UnmarshalJSON(data []byte) error {
	var creators = make([][]json_people, 3, 3)

	if err := json.Unmarshal(data, &creators); err != nil {
		return err
	}
	j.Directors = creators[0]
	j.Actors = creators[1]
	j.Producers = creators[2]
	return nil
}

type peoples []json_people

func (p *peoples) String() string {
	var names = make([]string, 0, len(*p))
	for _, pp := range *p {
		names = append(names, pp.NameRu)
	}

	return strings.Join(names, ", ")
}

type json_people struct {
	NameRu         string    `json:"nameRu"`
	NameEn         string    `json:"nameEn"`
	ProfessionText string    `json:"professionText"`
	ProfessionKey  json_prof `json:"professionKey"`
}

type json_prof string

func (j *json_prof) UnmarshalJSON(data []byte) error {
	for _, prof := range known_profs {
		if bytes.Equal([]byte(prof), bytes.Trim(data, `"`)) {
			*j = prof
			return nil
		}
	}
	return fmt.Errorf("unknown `json_prof`: %s", data)
}

const (
	director json_prof = "director"
	actor    json_prof = "actor"
	producer json_prof = "producer"
)

var known_profs = []json_prof{actor, director, producer}
