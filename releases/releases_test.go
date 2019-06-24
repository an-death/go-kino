package releases

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRelease_RaitingCollor(t *testing.T) {
	tests := []struct {
		name string
		init Release
		want string
	}{
		{name: "green_gt_7", init: Release{Rating: 7.1}, want: "#3bb33b"},
		{name: "gray_lt_7", init: Release{Rating: 5}, want: "#aaa"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.init.RatingColor())
		})
	}
}

func TestReleases_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name   string
		tested Releases
		want   bool
	}{
		{name: "LT", tested: Releases{Release{Date: time.Unix(1, 0)}, Release{Date: time.Unix(2, 0)}}, want: true},
		{name: "GT", tested: Releases{Release{Date: time.Unix(2, 0)}, Release{Date: time.Unix(1, 0)}}, want: false}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.tested.Less(0, 1))
		})
	}
}

func TestReleases_Sorted(t *testing.T) {
	var unsorted = Releases{Release{Date: time.Unix(2, 0)}, Release{Date: time.Unix(1, 0)}}
	var expected = Releases{Release{Date: time.Unix(1, 0)}, Release{Date: time.Unix(2, 0)}}
	sort.Sort(unsorted)
	assert.True(t, sort.IsSorted(unsorted))
	assert.Equal(t, expected, unsorted)
}
