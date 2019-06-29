package main

import (
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/an-death/go-kino/providers/releases"
)

var PORT string
var releaseProvider releases.ReleaseProvider
var infoTpl *template.Template

func init() {
	PORT = os.Getenv("PORT")

	if PORT == "" {
		PORT = "8000"
	}
	releaseProvider = releases.NewKinopoiskProvider()
	infoTpl = template.Must(template.ParseFiles("html/index.html", "html/movie.html"))
}

func handler1(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	newReleases := releaseProvider.GetReleases(now.AddDate(0, -1, 0), now)
	infoTpl.Execute(w, newReleases)
	log.Printf("Request done %s", time.Now().Sub(now))
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handler1)
	http.ListenAndServe(":"+PORT, nil)
}
