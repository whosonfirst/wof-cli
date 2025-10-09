package edit

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/whosonfirst/go-reader/v2"
)

type editReader struct {
	reader.Reader
	uri_map *sync.Map
	root    *os.Root
}

// Read will open and return an empty `io.ReadSeekCloser` for any value of 'path'.
func (r *editReader) Read(ctx context.Context, path string) (io.ReadSeekCloser, error) {

	fname := filepath.Base(path)

	_, exists := r.uri_map.Load(fname)

	if !exists {
		return nil, fmt.Errorf("Not found")
	}

	return r.root.Open(fname)
}

// Exists returns a boolean value indicating whether 'path' already exists (meaning it will always return false).
func (r *editReader) Exists(ctx context.Context, path string) (bool, error) {

	fname := filepath.Base(path)

	_, exists := r.uri_map.Load(fname)

	if !exists {
		return false, nil
	}

	_, err := r.root.Stat(fname)

	if err != nil {
		return false, err
	}

	return true, nil
}

// ReaderURI returns the value of 'path'.
func (r *editReader) ReaderURI(ctx context.Context, path string) string {
	return path
}
