package filter

import (
	"fmt"
	"os"
	"strings"
)

const (
	DefaultExcludes = "*/.DS_Store,*/.AppleDouble,*.swp,*~"
)

type FileSetBuilder interface {
	InDir(string) FileSetBuilder
	WithInclude(...string) FileSetBuilder
	WithExclude(...string) FileSetBuilder
	WithDefaultExcludes() FileSetBuilder
	withFilesystem(filesystem) FileSetBuilder
	FileSet() (FileSet, error)
}

func NewFileSetBuilder() FileSetBuilder {
	return &filesetBuilder{
		fs:       osFS{},
		includes: make(map[string]struct{}),
		excludes: make(map[string]struct{}),
	}
}

type filesetBuilder struct {
	dir                string
	includes, excludes map[string]struct{}
	fs                 filesystem
}

func (b *filesetBuilder) InDir(path string) FileSetBuilder {
	b.dir = path
	return b
}

func (b *filesetBuilder) WithInclude(includes ...string) FileSetBuilder {
	for _, i := range includes {
		b.includes[i] = struct{}{}
	}
	return b
}

func (b *filesetBuilder) WithExclude(excludes ...string) FileSetBuilder {
	for _, e := range excludes {
		b.excludes[e] = struct{}{}
	}
	return b
}

func (b *filesetBuilder) WithDefaultExcludes() FileSetBuilder {
	b.WithExclude(strings.Split(DefaultExcludes, ",")...)
	return b
}

func (b *filesetBuilder) withFilesystem(fs filesystem) FileSetBuilder {
	b.fs = fs
	return b
}

func (b *filesetBuilder) FileSet() (FileSet, error) {
	dir, err := b.openDir(b.dir)
	if err != nil {
		return FileSet{}, err
	}
	fileSet := &FileSet{
		dir: dir,
	}

	if len(b.includes) == 0 {
		b.includes["*"] = struct{}{}
	}

	for include := range b.includes {
		fileSet.includes = append(fileSet.includes, NewGlobMatcher(include))
	}
	for exclude := range b.excludes {
		fileSet.excludes = append(fileSet.excludes, NewGlobMatcher(exclude))
	}

	return *fileSet, nil
}

func (b *filesetBuilder) openDir(path string) (*os.File, error) {
	if stat, err := b.fs.Stat(b.dir); err != nil {
		return nil, err
	} else if !stat.IsDir() {
		return nil, fmt.Errorf("Provided path %s is not a directory", b.dir)
	}

	return b.fs.Open(b.dir)
}
