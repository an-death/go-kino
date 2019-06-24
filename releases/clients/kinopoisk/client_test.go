package kinopoisk

import (
	"net/http"
	"testing"
	"time"

	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
)

func mockTime() func() {
	patch := monkey.Patch(time.Now, func() time.Time { return time.Unix(0, 0) })
	return patch.Unpatch
}

func Test_client_createSignature(t *testing.T) {
	var expected = "098f6bcd4621d373cade4e832627b4f6"
	testClient := &client{}
	assert.Equal(t, expected, testClient.createSignature("test"))
}

func Test_client_uriWithUUID(t *testing.T) {
	var testCases = []struct{ name, test, expected string }{
		{"with out any query", "/uri", "/uri?uuid=test_uuid"},
		{"with one query arg", "/uri?first=1", "/uri?first=1&uuid=test_uuid"},
		{"with two query arg", "/uri?first=1&second=2", "/uri?first=1&second=2&uuid=test_uuid"},
		//	{"with boolquery arg", "/uri?first", "/uri?first&uuid=test_uuid"},
	}
	testClient := &client{
		uuid: "test_uuid",
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			want, err := testClient.uriWithUUID(tc.test)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, want)
		})
	}
}

type mockDoer struct {
	savedRequest    *http.Request
	returnsResponse *http.Response
	returnsErr      error
}

func (m *mockDoer) Do(r *http.Request) (*http.Response, error) {
	m.savedRequest = r
	return m.returnsResponse, m.returnsErr
}

func Test_client_Request(t *testing.T) {
	defer mockTime()()
	var mock = &mockDoer{}
	testClient := &client{
		Doer:     mock,
		clientId: "testClient",
		uuid:     "test",
		baseUrl:  "http://some.url",
	}

	var baseUri = "/uri"
	testClient.Request("GET", baseUri)
	wantR := mock.savedRequest

	assert.Equal(t, "0", wantR.Header.Get("X-TIMESTAMP"))
	assert.Equal(t, "some.url", wantR.Host)
	assert.Equal(t, "3a789fb1bbb08405663ebb008d9006f1", wantR.Header.Get("X-SIGNATURE"))
}
