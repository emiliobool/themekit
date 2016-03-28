package filter

func NewMatcherSet() *MatcherSet {
	return &MatcherSet{}
}

type MatcherSet struct {
	includes, excludes Matchers
}

// TODO: these don't work because we're on values. ugh.
func (p *MatcherSet) AddInclude(m Matcher) {
	p.includes = append(p.includes, m)
}

func (p *MatcherSet) AddExclude(m Matcher) {
	p.excludes = append(p.excludes, m)
}

// Matches returns true if no exclude and at least one include match the input
func (p *MatcherSet) Matches(input string) bool {
	if p.excludes.Matches(input) {
		return false
	}
	result := p.includes.Matches(input)

	return result
}
