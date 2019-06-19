package main

import (
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/an-death/go-kino/releases"
)

var PORT string
var releaseProvider releases.ReleaseProvider

func init() {
	PORT = os.Getenv("PORT")

	if PORT == "" {
		PORT = "8000"
	}
	releaseProvider = releases.NewKinopoiskProvider()
}

func handler1(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index.html").ParseFiles("html/index.html", "html/movie.html")
	if err != nil {
		panic(err)
	}
	now := time.Now()

	newReleases := releaseProvider.GetReleases(now.AddDate(0, -1, 0), now)
	t.Execute(w, newReleases)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handler1)
	http.ListenAndServe(":"+PORT, nil)
}
