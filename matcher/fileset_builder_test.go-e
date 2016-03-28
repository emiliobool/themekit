package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultExcludes(t *testing.T) {
	fs, _ := NewFileSetBuilder().WithDefaultExcludes().WithInclude("*").FileSet()
	excludedFileNames := []string{"Thumbs.db", ".gitignore", "sub/.gitignore", ".temp~", "foo"}

	for _, name := range excludedFileNames {
		assert.False(t, fs.Matches(name))
	}
}
