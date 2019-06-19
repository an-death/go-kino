package releases

import "time"

type Release struct {
	OriginName string
	Raiting    float64
	NameRu     string
	InfoTable  []Info

	PosterUrl string
	Torrents  []Torrent
}

func (m *Release) RaitingCollor() string {
	if m.Raiting > 7 {
		return "#3bb33b"
	}
	return "#aaa"
}

func (m *Release) IsDisplayOrigin() bool {
	return len(m.OriginName) > 0
}

type Info struct {
	Key string
	Val string
}

type Torrent struct {
	Link string
	Type string
}

type ReleaseProvider interface {
	GetReleases(from, to time.Time) []Release
}

func NewReleaseProvider() ReleaseProvider {
	return &mockReleaseProvider{}
}
