package kinopoisk

import (
	"strings"
	"time"
)

type responseContainer struct {
	IsSuccess bool         `json:"success"`
	Data      responseData `json:"data"`
}
type responseData struct {
	Items []MovieItem      `json:"items"`
	Stats responseDataStat `json:"stats"`
}
type responseDataStat struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type MovieItem struct {
	Id            int    `json:"id"`
	Slug          string `json:"slug"`
	Title         string `json:"title"`
	OriginalTitle string `json:"originalTitle"`
	Year          int    `json:"year"`
	Poster        struct {
		BaseUrl string `json:"baseUrl"`
		Url     string `json:"url"`
	} `json:"poster"`
	Genres        []named `json:"genres"`
	Countries     []named `json:"countries"`
	Rating        raiting `json:"rating"`
	Expectations  raiting `json:"expectations"`
	CurrentRating string  `json:"currentRating"`
	Serial        bool    `json:"serial"`
	Duration      int     `json:"duration"`
	TrailerId     int     `json:"trailerId"`
	ContextData   struct {
		IsDigital   bool        `json:"isDigital"`
		ReleaseDate releaseDate `json:"releaseDate"`
	} `json:"contextData"`
}

type raiting struct {
	Value float64 `json:"value"`
	Count int     `json:"count"`
	Ready bool    `json:"ready"`
}
type named struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
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
