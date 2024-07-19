package reader

import (
	"context"
	"io"
	"os"
)

const STDIN string = "-"

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
