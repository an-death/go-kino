package frontend

import (
	"log"
	"net/http"
	"time"

	"github.com/an-death/go-kino/releases"
	"github.com/gin-gonic/gin"
)

var releaseProvider releases.ReleaseProvider

func init() {
	releaseProvider = releases.NewKinopoiskProvider()
}

func index(c *gin.Context) {
	now := time.Now()

	newReleases := releaseProvider.GetReleases(now.AddDate(0, -1, 0), now)
	c.HTML(http.StatusOK, "index.html", newReleases)
	log.Printf("Request done %s", time.Now().Sub(now))
}

func routes() *gin.Engine {
	router := gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLFiles("html/index.html", "html/movie.html")
	router.GET("/", index)
	return router
}

func Run(addr string) error {
	router := routes()
	return router.Run(addr)
}
