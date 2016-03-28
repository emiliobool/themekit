package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyMatcherSetMatchesNothing(t *testing.T) {
	ps := NewMatcherSet()

	assert.False(t, ps.Matches("something"))
}

func TestMatcherSetDoesNotMatchIfExcludesMatch(t *testing.T) {
	ps := NewMatcherSet()

	ps.AddInclude(NewAlwaysMatcher())
	ps.AddExclude(NewAlwaysMatcher())

	assert.False(t, ps.Matches("something"))
}

func TestMatcherSetMatchesWhenIncludeAndNoExcludeMatches(t *testing.T) {
	ps := NewMatcherSet()

	ps.AddInclude(NewAlwaysMatcher())
	ps.AddExclude(NewNeverMatcher())

	assert.True(t, ps.Matches("something"))
}

func TestMatcherSetDoesNotMatchWithoutIncludes(t *testing.T) {
	ps := NewMatcherSet()

	assert.False(t, ps.Matches("something"))
}
