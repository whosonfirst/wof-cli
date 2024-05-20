package reader

import (
	"context"
	"io"
	"os"
)

func ReadCloserFromURI(ctx context.Context, uri string) (io.ReadCloser, bool, error) {

	r, err := os.Open(uri)

	if err != nil {
		return nil, false, err
	}

	return r, false, nil
}
