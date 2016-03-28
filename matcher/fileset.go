package matcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	RelativeToDir = "relative"
	Unchanged     = "unchanged"
	Absolute      = "absolute"
)

type FileSetSettings struct {
	PathsRelativeTo string
}

type FileSet struct {
	MatcherSet
	FileSetSettings
	dir *os.File
}

// Takes a path and makes it relative to the basepath of this fileset
func (fs FileSet) relativeToDir(path string) (string, error) {
	base := fs.dir.Name()
	result, err := filepath.Rel(base, path)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (fs FileSet) AllFiles() ([]string, error) {
	var result []string
	files, _ := fs.Files()
	for f := range files {
		result = append(result, f)
	}

	return result, nil
}

func (fs FileSet) relativizer() func(string) (string, error) {
	switch fs.PathsRelativeTo {
	case Absolute:
		return func(s string) (string, error) {
			return filepath.Abs(s)
		}
	case Unchanged:
		return func(s string) (string, error) { return s, nil }
	default:
		return fs.relativeToDir
	}
}

func (fs FileSet) Files() (chan string, error) {
	results := make(chan string)

	go fs.walkAllMatches(results, fs.relativizer())

	return results, nil
}

func (fs FileSet) walkAllMatches(results chan string, relative func(string) (string, error)) {
	defer close(results)
	filepath.Walk(fs.dir.Name(), func(path string, fi os.FileInfo, err error) error {
		matched := fs.Matches(path)
		if !matched || fi.IsDir() {
			return nil
		}

		relativePath, _ := relative(path)
		results <- relativePath

		return nil
	})
}

func (fs FileSet) String() string {
	incs := []string{}
	for _, m := range fs.includes {
		incs = append(incs, fmt.Sprintf("%s", m))
	}
	excl := []string{}
	for _, m := range fs.excludes {
		excl = append(excl, fmt.Sprintf("%s", m))
	}
	return fmt.Sprintf("FileSet{dir=%s; includes=%s; excludes=%s}", fs.dir.Name(), strings.Join(incs, ","), strings.Join(excl, ","))
}

type filesystem interface {
	Open(name string) (*os.File, error)
	Stat(name string) (os.FileInfo, error)
}

type osFS struct{}

func (osFS) Open(name string) (*os.File, error)    { return os.Open(name) }
func (osFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }
