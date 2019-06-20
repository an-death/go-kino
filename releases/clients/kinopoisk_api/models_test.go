package kinopoisk

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const apiResponse string = `{
   "success":true,
   "data":{
      "items":[
         {
            "id":1,
            "contextData":{
               "releaseDate":"2019-06-01"
            }
         }
      ],
      "stats":{
         "total":19,
         "limit":1000,
         "offset":0
      }
   }
}
`
const emptyApiResponse string = `
{
   "success":true,
   "data":{
      "items":[],
      "stats":{
         "total":0,
         "limit":1000,
         "offset":0
      }
   }
}
`

func Test_releaseDateMarshalJSON(t *testing.T) {
	const expected = "2019-12-23"
	date, err := time.Parse("2006-01-02", expected)
	assert.NoError(t, err, "test not initialize: wrong expected date")
	actual, err := releaseDate(date).MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `"`+expected+`"`, string(actual))
}

func Test_releaseDateUnmarshalJSON(t *testing.T) {
	const expected = "2019-06-01"
	var d releaseDate
	err := d.UnmarshalJSON([]byte(expected))
	assert.NoError(t, err)
	assert.Equal(t, expected, time.Time(d).Format("2006-01-02"))
}

func Test_releaseDateAsDate(t *testing.T) {
	expected, err := time.Parse("2006-01-02", "2019-12-13")
	assert.NoError(t, err, "test not initialize: wrong expected date")
	assert.Equal(t, expected, releaseDate(expected).AsDate())
}

func TestJSONParseModel(t *testing.T) {
	var testCases = []struct {
		name        string
		apiResponse string
	}{
		{name: "empty", apiResponse: emptyApiResponse},
		{name: "filled", apiResponse: apiResponse},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			var response responseContainer
			expected := []byte(tc.apiResponse)
			assert.NoError(t, json.Unmarshal(expected, &response))
			t.Log(response)
			data, err := json.Marshal(response)
			t.Log(string(data))
			assert.NoError(t, err)
			assert.JSONEq(t, tc.apiResponse, string(data))
		})
	}
}
