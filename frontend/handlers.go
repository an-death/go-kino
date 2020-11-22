package frontend

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

func proxy(c *gin.Context) {
	target := c.Query("url")
	if target == "" {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	response, err := http.Get(target)
	if err != nil || response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	proxyResponse(response, c)
}

func proxyResponse(response *http.Response, c *gin.Context) {
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
		fmt.Println(err)
		c.Status(http.StatusServiceUnavailable)
		return
	}
	result, err := replaceUrls(response.Body)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusServiceUnavailable)
		return
	}

	oldBody, newBody := response.Body, result
	defer oldBody.Close()
	response.Body = newBody
	proxyResponse(response, c)
}

func replaceUrls(source io.Reader) (io.ReadCloser, error) {
	doc, err := goquery.NewDocumentFromReader(source)
	if err != nil {
		return nil, err
	}

	base, _ := url.Parse("http://rutor.info/")
	doc.Find("#index").Find("a").Each(func(_ int, link *goquery.Selection) {
		href, ok := link.Attr("href")
		if ok {
			u, err := url.Parse(href)
			if err != nil {
				panic(err)
			}
			if u.IsAbs() {
				link.ReplaceWithHtml("/proxy?url=" + href)
			} else {
				local, _ := url.Parse("/proxy")
				q := local.Query()
				q.Set("url", base.ResolveReference(u).String())
				local.RawQuery = q.Encode()

				link.ReplaceWithHtml(local.String())
			}
		}
	})
	ctn, err := doc.Html()
	if err != nil {
		return nil, err
	}

	return ioutil.NopCloser(strings.NewReader(ctn)), nil
}

func routes() *gin.Engine {
	router := gin.Default()
	router.Static("/static", "./frontend/static")
	router.LoadHTMLFiles("frontend/html/index.html", "frontend/html/movie.html", "frontend/html/search.html")
	router.GET("/", searchGet)
	router.GET("/proxy", proxy)
	router.GET("/search", searchGet)

	return router
}

func Run(addr string) error {
	router := routes()
	return router.Run(addr)
}
