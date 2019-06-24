package releases

import "time"

type Release struct {
	OriginName string
	NameRu     string
	InfoTable  map[string]string
	Rating     float64
	PosterUrl  string
	WebURL     string
	Date       time.Time
	Torrents   []Torrent
}

func (m *Release) RatingColor() string {
	if m.Rating > float64(7) {
		return "#3bb33b"
	}
	return "#aaa"
}

func (m *Release) IsDisplayOrigin() bool {
	return len(m.OriginName) > 0
}

type Torrent struct {
	Link string
	Type string
}

type ReleaseProvider interface {
	GetReleases(from, to time.Time) []Release
}

type Releases []Release

func (r Releases) Len() int {
	return len(r)
}

func (r Releases) Less(i int, j int) bool {
	return r[i].Date.Before(r[j].Date)
}

func (r Releases) Swap(i int, j int) {
	r[i], r[j] = r[j], r[i]
}
