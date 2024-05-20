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
