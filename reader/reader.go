package reader

import (
	"context"
	"io"
	"os"

	"github.com/whosonfirst/go-ioutil"
	wof_reader "github.com/whosonfirst/go-reader/v2"
)

const READER_SCHEME string = "wof-cli"
const STDIN string = "-"

func init() {
	ctx := context.Background()
	err := wof_reader.RegisterReader(ctx, READER_SCHEME, NewReader)

	if err != nil {
		panic(err)
	}
}

func ReadCloserFromURI(ctx context.Context, uri string) (io.ReadCloser, bool, error) {

	if uri == STDIN {
		return os.Stdin, true, nil
	}

	r, err := os.Open(uri)

	if err != nil {
		return nil, false, err
	}

	return r, false, nil
}

func BytesFromURI(ctx context.Context, uri string) ([]byte, error) {

	r, is_stdin, err := ReadCloserFromURI(ctx, uri)

	if err != nil {
		return nil, err
	}

	if !is_stdin {
		defer r.Close()
	}

	return io.ReadAll(r)
}

// Reader implements the `whosonfirst/go-reader/v2.Reader` interface.
type Reader struct {
	wof_reader.Reader
}

func NewReader(ctx context.Context, uri string) (wof_reader.Reader, error) {

	r := &Reader{}
	return r, nil
}

// Read will open and return an empty `io.ReadSeekCloser` for any value of 'path'.
func (r *Reader) Read(ctx context.Context, path string) (io.ReadSeekCloser, error) {

	rdr, _, err := ReadCloserFromURI(ctx, path)

	if err != nil {
		return nil, err
	}

	return ioutil.NewReadSeekCloser(rdr)
}

// ReaderURI returns the value of 'path'.
func (r *Reader) ReaderURI(ctx context.Context, path string) string {
	return path
}
