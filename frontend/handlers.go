package frontend

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/an-death/go-kino/providers/releases"
	"github.com/gin-gonic/gin"
)

var releaseProvider releases.ReleaseProvider

func init() {
	releaseProvider = releases.NewKinopoiskProvider()
}

// func indexHandler(c *gin.Context) {
// 	now := time.Now()
// 	query := url.Values{
// 		"from": []string{now.AddDate(0, -1, 0).Format("2006-01")},
// 		"to":   []string{now.Format("2006-01")},
// 	}
// 	newUrl := url.URL{
// 		Path:     "releases",
// 		RawQuery: query.Encode(),
// 	}
// 	c.Redirect(http.StatusMovedPermanently, newUrl.RequestURI())
// }

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
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

func torrentFileProxy(c *gin.Context) {
	torrentFileUrl := c.Query("url")
	if torrentFileUrl == "" {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	response, err := http.Get(torrentFileUrl)
	if err != nil || response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	reader := response.Body
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	extraHeaders := map[string]string{}
	for headerName, headerValue := range response.Header {
		extraHeaders[headerName] = strings.Join(headerValue, "; ")
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}

func searchGet(c *gin.Context) {
	searchFilmName := c.Query("query")
	if searchFilmName == "" {
		c.HTML(http.StatusOK, "search.html", nil)
		return
	}
	searchQuery := fmt.Sprintf(`http://rutor.info/search/%s/?search_method=0&search_in=0&category=0&s_ad=0`, searchFilmName)
	response, err := http.Get(searchQuery)
	if err != nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	reader := response.Body
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	extraHeaders := map[string]string{}
	for headerName, headerValue := range response.Header {
		extraHeaders[headerName] = strings.Join(headerValue, "; ")
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}

func routes() *gin.Engine {
	router := gin.Default()
	router.Static("/static", "./frontend/static")
	router.LoadHTMLFiles("frontend/html/index.html", "frontend/html/movie.html", "frontend/html/search.html")
	router.GET("/", gin.WrapF(handleHTTP))
	router.GET("/releases", releasesHandler)
	router.GET("/proxy", torrentFileProxy)
	router.GET("/search", searchGet)

	return router
}

func Run(addr string) error {
	router := routes()
	return router.Run(addr)
}
