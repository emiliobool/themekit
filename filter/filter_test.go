package filter

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	m := Must(NewRegexpMatcher(""))
	assert.NotNil(t, m)
}

func TestGlobMatcher(t *testing.T) {

	type td struct {
		input    string
		expected bool
	}
	data := []struct {
		pattern string
		tests   []td
	}{
		{
			"", // Empty glob doesn't really do much
			[]td{
				{"", true},
				{"tmp", false},
			},
		},
		{
			"tmp",
			[]td{
				{"tmp", true},
				{"tmp/foo.txt", false},
				{"tmp/foo.tmp", false},
				{"bar/tmp.txt", false},
			},
		},
		{
			".DS_Store",
			[]td{
				{".DS_Store", true},
			},
		},
		{
			"tmp/*",
			[]td{
				{"tmp", false},
				{"tmp/", true},
				{"tmp/foo.txt", true},
				{"bar/tmp/bar.txt", false},
			},
		},
		{
			"foo/**/bar",
			[]td{
				{"foo/bar", false},
				{"foo/blargh/bar", true},
				{"foo/some/other/stuff/bar", true},
			},
		},
		{
			"**/*.png",
			[]td{
				{"foo.png", false},
				{"bar/foo.png", true},
				{"foo/png/bar.txt", false},
			},
		},
	}

	for _, d := range data {
		matcher := NewGlobMatcher(d.pattern)

		for _, test := range d.tests {
			input := test.input
			expected := test.expected

			matched, _ := matcher.Matches(input)
			assert.Equal(t, expected, matched, fmt.Sprintf("Pattern \"%s\" matched \"%s\": %t but was supposed to be %t", d.pattern, input, matched, expected))
		}
	}
}
