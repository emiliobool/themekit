package filter

import (
	"regexp"

	glob "github.com/ryanuber/go-glob"
)

type Matchers []Matcher

// Matches returns true if any of the contained matchers match the input
func (m Matchers) Matches(input string) bool {
	for _, matcher := range m {
		if matched, _ := matcher.Matches(input); matched {
			return true
		}
	}
	return false
}

type Matcher interface {
	// Matches returns true if this matcher matches the given input
	Matches(string) (bool, error)
}

func NewRegexpMatcher(regex string) (Matcher, error) {
	return regexpMatcher{
		regex: regexp.MustCompile(regex),
	}, nil
}

type regexpMatcher struct {
	regex *regexp.Regexp
}

func (rm regexpMatcher) Matches(input string) (bool, error) {
	return rm.regex.MatchString(input), nil
}

func (rm regexpMatcher) String() string {
	return rm.regex.String()
}

func NewGlobMatcher(pattern string) Matcher {
	return globMatcher{pattern: pattern}
}

type globMatcher struct {
	pattern string
}

func (gm globMatcher) Matches(input string) (bool, error) {
	return glob.Glob(gm.pattern, input), nil
}

func NewAlwaysMatcher() Matcher {
	return alwaysMatcher{}
}

type alwaysMatcher struct{}

func (a alwaysMatcher) Matches(string) (bool, error) {
	return true, nil
}

func NewNeverMatcher() Matcher {
	return neverMatcher{}
}

type neverMatcher struct{}

func (nm neverMatcher) Matches(string) (bool, error) {
	return false, nil
}
