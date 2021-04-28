package cacheablefs

import (
	"io/fs"
	"time"
)

var startAt = time.Now()

// make embedFS is cacheable by Last-Modified header.
type wrappedFS struct{ fs.FS }
type wrappedFile struct{ fs.File }
type wrappedStat struct{ fs.FileInfo }

func (m wrappedStat) ModTime() time.Time {
	return startAt
}

func (w wrappedFile) Stat() (fs.FileInfo, error) {
	d, err := w.File.Stat()
	if err != nil {
		return nil, err
	}
	return wrappedStat{d}, nil
}

func (w wrappedFS) Open(name string) (fs.File, error) {
	f, err := w.FS.Open(name)
	if err != nil {
		return nil, err
	}
	return wrappedFile{f}, nil
}

func Wrap(in fs.FS) fs.FS {
	return wrappedFS{in}
}
