package frontend

import (
	"net/http"
	"net/url"
	"time"

	"github.com/an-death/go-kino/providers/releases"
	"github.com/gin-gonic/gin"
)

var releaseProvider releases.ReleaseProvider

func init() {
	releaseProvider = releases.NewKinopoiskProvider()
}

func indexHandler(c *gin.Context) {
	now := time.Now()
	query := url.Values{
		"from": []string{now.AddDate(0, -1, 0).Format("2006-01")},
		"to":   []string{now.Format("2006-01")},
	}
	newUrl := url.URL{
		Path:     "releases",
		RawQuery: query.Encode(),
	}
	c.Redirect(http.StatusMovedPermanently, newUrl.RequestURI())
}

type requestTimeRange struct {
	From time.Time `form:"from" time_format:"2006-01" binding:"required"`
	To   time.Time `form:"to" time_format:"2006-01" binding:"required"`
}

func releasesHandler(c *gin.Context) {
	var rr requestTimeRange
	if err := c.ShouldBindQuery(&rr); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	newReleases := releaseProvider.GetReleases(rr.From, rr.To)
	c.HTML(http.StatusOK, "index.html", newReleases)
}

func routes() *gin.Engine {
	router := gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLFiles("frontend/html/index.html", "frontend/html/movie.html")
	router.GET("/", indexHandler)
	router.GET("/releases", releasesHandler)

	return router
}

func Run(addr string) error {
	router := routes()
	return router.Run(addr)
}
