package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyPatternSetMatchesNothing(t *testing.T) {
	ps := NewPatternSet()

	assert.False(t, ps.Matches("something"))
}

func TestPatternSetDoesNotMatchIfExcludesMatch(t *testing.T) {
	ps := NewPatternSet()

	ps.AddInclude(NewAlwaysMatcher())
	ps.AddExclude(NewAlwaysMatcher())

	assert.False(t, ps.Matches("something"))
}

func TestPatternSetMatchesWhenIncludeAndNoExcludeMatches(t *testing.T) {
	ps := NewPatternSet()

	ps.AddInclude(NewAlwaysMatcher())
	ps.AddExclude(NewNeverMatcher())

	assert.True(t, ps.Matches("something"))
}

func TestPatternSetDoesNotMatchWithoutIncludes(t *testing.T) {
	ps := NewPatternSet()

	assert.False(t, ps.Matches("something"))
}
