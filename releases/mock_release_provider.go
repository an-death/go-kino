package releases

import "time"

var todo = Release{
	NameRu:     "Семья по-быстрому",
	OriginName: "Instant Family",
	InfoTable: []Info{
		{"Описание", "DEsc"},
		{"Год", "2018"},
		{"Страна", "США"},
		{"Жанр", "Драма, Комедия"},
		{"Возраст", "16"},
		{"Длительность", "1.58"},
		{"Рейтинг КиноПоиск", "7.3"},
		{"Рейтинг IMDb", "7.4"},
		{"Режисер", "Шон Андерс"},
		{"Актеры", "Марк Оулберг, ,Шон Андерс"},
		{"Дата выхода", "2019-06-01"},
	},
	PosterUrl: "https://st.kp.yandex.net/images/film_iphone/iphone_1108494.jpg?d=20190523114230&width=360",
	Torrents:  []Torrent{{"http://top-tor.org/download/696534", "BDRip 1080p"}},

	Raiting: 10,
}

type mockReleaseProvider struct{}

func (r mockReleaseProvider) GetReleases(_, _ time.Time) []Release {
	return []Release{todo}
}
