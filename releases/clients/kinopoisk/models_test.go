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
			"title": "testTitle",
			"poster": {
				"url": "test_poster_url"
			},
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

func TestJSONParse_json_prof(t *testing.T) {
	var testCases = []struct {
		name     string
		data     []byte
		expected json_prof
	}{
		{name: "director", data: []byte(`"director"`), expected: director},
		{name: "actor", data: []byte(`"actor"`), expected: actor},
		{name: "producer", data: []byte(`"producer"`), expected: producer},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var prof json_prof
			assert.NoError(t, json.Unmarshal(tc.data, &prof))
			assert.Equal(t, tc.expected, prof)
		})
	}
}

func TestJSONParse_json_people(t *testing.T) {
	var expected = json_people{
		NameRu:         "Имя",
		NameEn:         "NameEn",
		ProfessionText: "Режисёр",
		ProfessionKey:  director,
	}
	var testData = []byte(`{
	"nameRu": "Имя",
	"nameEn": "NameEn",
	"professionText": "Режисёр",
	"professionKey": "director"
	}`)
	var testPeople json_people

	assert.NoError(t, json.Unmarshal(testData, &testPeople))
	assert.Equal(t, expected, testPeople)
}

func TestJSONParse_json_creators(t *testing.T) {
	var expected = json_creators{
		Directors: []json_people{
			json_people{ProfessionKey: director},
		},
		Actors: []json_people{
			json_people{ProfessionKey: actor},
		},
		Producers: []json_people{
			json_people{ProfessionKey: producer},
			json_people{ProfessionKey: producer},
		},
	}
	var testData = []byte(`
		[
			[{"professionKey": "director"}],
			[{"professionKey":"actor"}],
			[{"id":1, "professionKey":"producer"}, {"id":2, "professionKey":"producer"}]
		]
	`)
	var testCreators json_creators

	assert.NoError(t, json.Unmarshal(testData, &testCreators))
	assert.Equal(t, expected, testCreators)
}

func TestStringer_peoples(t *testing.T) {
	var testCases = []struct {
		name     string
		expected string
		actual   peoples
	}{
		{name: "one", expected: "People1", actual: peoples{json_people{NameRu: "People1"}}},
		{name: "three", expected: "People1, People2, People3",
			actual: peoples{
				json_people{NameRu: "People1"},
				json_people{NameRu: "People2"},
				json_people{NameRu: "People3"},
			}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.actual.String())
		})
	}
}

func TestFilmDetail_Rating(t *testing.T) {
	var expected = 10.0
	var testData = []byte(`{
"ratingData": {
	"rating": "10",
	"ratingIMDb": "10.0"
}}`)
	var actual FilmDetail

	assert.NoError(t, json.Unmarshal(testData, &actual))
	assert.Equal(t, expected, actual.Rating())
}
