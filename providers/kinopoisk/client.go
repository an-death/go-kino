package kinopoisk

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/an-death/go-kino/providers"
)

const (
	KINOPOISK_API_SALT = "IDATevHDS7"
)

type client struct {
	providers.Doer
	baseUrl        string
	clientId, uuid string
}

func NewAPIClient(baseUrl string, do providers.Doer) providers.APIClient {
	var clientId = make([]byte, 16, 16)
	var uuid = make([]byte, 12, 12)
	rand.Read(clientId)
	rand.Read(uuid)
	return &client{
		Doer:     do,
		clientId: fmt.Sprintf("%x", clientId),
		uuid:     fmt.Sprintf("%x", uuid),
		baseUrl:  baseUrl,
	}
}

func (c *client) Request(method, uri string) (*http.Response, error) {
	var err error
	uri, err = c.uriWithUUID(uri)
	if err != nil {
		return nil, err
	}

	timeNow := time.Now()
	timestamp := strconv.FormatInt(timeNow.UnixNano(), 10)
	signature := c.createSignature(uri + timestamp + KINOPOISK_API_SALT)
	request, err := http.NewRequest(method, c.baseUrl+uri, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "Android client (6.0.1 / api23), ru.kinopoisk/4.6.5 (86)")
	request.Header.Add("Image-Scale", "3")
	request.Header.Add("device", "android")
	request.Header.Add("ClientId", c.clientId)
	request.Header.Add("countryID", "2")
	request.Header.Add("cityID", "1")
	request.Header.Add("Android-Api-Version", "23")
	request.Header.Add("clientDate", timeNow.Format("03:04 02.01.2006"))
	request.Header.Add("X-TIMESTAMP", timestamp)
	request.Header.Add("X-SIGNATURE", signature)

	return c.Doer.Do(request)
}

func (c *client) uriWithUUID(baseURI string) (string, error) {
	uri, err := url.Parse(baseURI)
	if err != nil {
		return "", nil
	}
	q := uri.Query()
	q.Add("uuid", c.uuid)
	uri.RawQuery = q.Encode()
	return uri.String(), nil
}

func (c *client) createSignature(from string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(from)))
}
