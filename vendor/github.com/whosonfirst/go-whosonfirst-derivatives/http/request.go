package http

import (
	"context"
	"fmt"
	"io"
	
	"github.com/whosonfirst/go-whosonfirst-derivatives"	
	wof_http "github.com/whosonfirst/go-whosonfirst/http"	
)

func FeatureFromRequestURI(ctx context.Context, prv derivatives.Provider, req_uri *wof_http.URI) ([]byte, error) {

	wof_id := req_uri.Id

	r, err := prv.GetFeature(ctx, wof_id, req_uri.URIArgs)	

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve feature for %d, %w", wof_id, err)
	}

	defer r.Close()

	return io.ReadAll(r)
}
