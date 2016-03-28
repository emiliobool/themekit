package filter

func NewPatternSet() *PatternSet {
	return &PatternSet{}
}

type PatternSet struct {
	includes, excludes Matchers
}

// TODO: these don't work because we're on values. ugh.
func (p *PatternSet) AddInclude(m Matcher) {
	p.includes = append(p.includes, m)
}

func (p *PatternSet) AddExclude(m Matcher) {
	p.excludes = append(p.excludes, m)
}

// Matches returns true if no exclude and at least one include match the input
func (p *PatternSet) Matches(input string) bool {
	if p.excludes.Matches(input) {
		return false
	}
	result := p.includes.Matches(input)

	return result
}
