package spelunker

import (
	"bytes"
	"context"
	"io"

	"github.com/whosonfirst/go-ioutil"
	"github.com/whosonfirst/go-whosonfirst/v4/derivatives"
	"github.com/whosonfirst/go-whosonfirst/v4/uri"
)

// SpelunkerDerivativesProvider implements the `go-whosonfirst/v4/derivatives.Provider` interface
// retrieving records from a `Spelunker` instance.
type SpelunkerDerivativesProvider struct {
	derivatives.Provider
	spelunker Spelunker
}

// NewDerivativesProvider returns an implementation of the `go-whosonfirst/v4/derivatives.Provider` interface
// using 'sp' to retrieve records.
func NewDerivativesProvider(sp Spelunker) derivatives.Provider {

	pr := &SpelunkerDerivativesProvider{
		spelunker: sp,
	}

	return pr
}

// GetFeatures returns an `io.ReadSeekCloser` for the record derived from 'id' and 'uri_args'.
func (pr *SpelunkerDerivativesProvider) GetFeature(ctx context.Context, id int64, uri_args *uri.URIArgs) (io.ReadSeekCloser, error) {

	body, err := pr.spelunker.GetFeatureForId(ctx, id, uri_args)

	if err != nil {
		return nil, err
	}

	br := bytes.NewReader(body)
	return ioutil.NewReadSeekCloser(br)
}
