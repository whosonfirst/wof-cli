package derivatives

import (
	"context"
	"io"

	"github.com/whosonfirst/go-whosonfirst-uri"
)

type NullProvider struct {
	Provider
}

func init() {

	err := RegisterProvider(context.Background(), "null", NewNullProvider)

	if err != nil {
		panic(err)
	}
}

func NewNullProvider(ctx context.Context, uri string) (Provider, error) {
	s := &NullProvider{}
	return s, nil
}

func (s *NullProvider) GetFeature(ctx context.Context, id int64, uri_args *uri.URIArgs) (io.ReadSeekCloser, error) {
	return nil, ErrNotFound
}
