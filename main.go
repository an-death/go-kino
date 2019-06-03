package main

import (
	"net/http"
	"os"
	"strconv"
	"text/template"
)

var PORT int

func init() {
	if port, ok := os.LookupEnv("PORT"); ok {
		PORT, _ = strconv.Atoi(port)
	}
	if PORT == 0 {
		PORT = 8000
	}
}

type Movie struct {
	OriginalName string
	Raiting      int
	NameRu       string
	InfoTable    []Info

	PosterUrl string
	Torrents  []Torrent
}

func (m *Movie) RaitingCollor() string {
	if m.Raiting > 7 {
		return "#3bb33b"
	}
	return "#aaa"
}

func (m *Movie) IsDisplayOrigin() bool {
	return len(m.OriginalName) > 0
}

type Info struct {
	Key string
	Val string
}

type Torrent struct {
	Link string
	Type string
}

var todo = Movie{
	NameRu:       "Семья по-быстрому",
	OriginalName: "Instant Family",
	InfoTable: []Info{
		{"Описание", "DEsc"},
		{"Год", "2018"},
		{"Страна", "США"},
		{"Жанр", "Драма, Комедия"},
		{"Возраст", "16"},
		{"Длительность", "1.58"},
		{"Рейтинг КиноПоиск", "7.3"},
		{"Рейтинг IMDb", "7.4"},
		{"Режисер", "Шон Андерс"},
		{"Актеры", "Марк Оулберг, ,Шон Андерс"},
		{"Дата выхода", "2019-06-01"},
	},
	PosterUrl: "https://st.kp.yandex.net/images/film_iphone/iphone_1108494.jpg?d=20190523114230&width=360",
	Torrents:  []Torrent{{"http://top-tor.org/download/696534", "BDRip 1080p"}},

	Raiting: 10,
}

func handler1(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index.html").ParseFiles("html/index.html", "html/movie.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, []Movie{todo})
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handler1)
	http.ListenAndServe(":8000", nil)
}
