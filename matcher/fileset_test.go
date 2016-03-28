package matcher

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFilesetReturnsErrorIfStatFails(t *testing.T) {
	_, err := NewFileSetBuilder().InDir("fixtures/does not exist").WithDefaultExcludes().WithExclude("*.txt", "*.min.js").FileSet()
	assert.NotNil(t, err)
}

func TestFilesetDoesNotMatchWithoutAnyIncludes(t *testing.T) {
}

func TestFilesetWithFixtureDir(t *testing.T) {
	fs, err := NewFileSetBuilder().InDir("fixtures/simple").WithDefaultExcludes().WithExclude("*.txt", "*.min.js", "*.tmp").FileSet()
	assert.Nil(t, err)

	files, err := fs.AllFiles()

	assert.Equal(t, []string{"src/images/foo.png", "src/js/foo.js", "src/js/tmp.js"}, files)
}

type mockedFS struct {
	mock.Mock
}

func (m mockedFS) Open(name string) (*os.File, error) {
	args := m.Called(name)

	file, _ := args.Get(0).(*os.File)
	return file, args.Error(1)
}

func (m mockedFS) Stat(name string) (os.FileInfo, error) {
	args := m.Called(name)

	file, _ := args.Get(0).(os.FileInfo)
	return file, args.Error(1)
}
