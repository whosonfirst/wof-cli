package derivatives

import (
	"context"
	"io"
	"net/url"

	_ "github.com/whosonfirst/go-reader-findingaid/v2"
	_ "github.com/whosonfirst/go-reader-github/v2"

	"github.com/whosonfirst/go-reader/v2"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

type ReaderProvider struct {
	Provider
	reader reader.Reader
}

func init() {

	err := RegisterProvider(context.Background(), "reader", NewReaderProvider)

	if err != nil {
		panic(err)
	}
}

func NewReaderProvider(ctx context.Context, uri string) (Provider, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	reader_uri := q.Get("reader-uri")

	r, err := reader.NewReader(ctx, reader_uri)

	if err != nil {
		return nil, err
	}

	s := &ReaderProvider{
		reader: r,
	}
	return s, nil
}

func (s *ReaderProvider) GetFeature(ctx context.Context, id int64, uri_args *uri.URIArgs) (io.ReadSeekCloser, error) {

	rel_path, err := uri.Id2RelPath(id, uri_args)

	if err != nil {
		return nil, err
	}

	return s.reader.Read(ctx, rel_path)
}
