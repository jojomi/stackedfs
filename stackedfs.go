package stackedfs

import (
	"fmt"
	"io/fs"
)

// implements fs.FS
type StackedFS struct {
	filesystems []fs.FS
}

// NewStackedFS returns a stacked filesystem. The order is important, a file is returned from the
// *first* filesystem that does not return an error when queried for the filename.
func NewStackedFS(filesystems ...fs.FS) *StackedFS {
	sfs := StackedFS{
		filesystems: filesystems,
	}
	return &sfs
}

func (s *StackedFS) AddFS(filesystem fs.FS) {
	s.filesystems = append(s.filesystems, filesystem)
}

func (s *StackedFS) Open(name string) (fs.File, error) {
	var (
		file fs.File
		err  error
	)
	for _, f := range s.filesystems {
		file, err = f.Open(name)
		if err == nil {
			return file, err
		}
	}
	return nil, fmt.Errorf("file not found in any of the underlaying %d filesystems of this StackedFS: %s", len(s.filesystems), name)
}
