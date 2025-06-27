package uris

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

var re_wofid = regexp.MustCompile(`^\d+(?:\-alt\-.*)?$`)

// ExpandURICallbackFunc is a callback function to invoked for each path (URI) resolved by ExpandURIsWithCallback
type ExpandURICallbackFunc func(context.Context, string) error

// ExpandURIsWithCallback is a wrapper function for processing one or more URIs and invoking 'cb' for each in the
// order that the URI was processed. Currently supported URI "expansions" are:
//   - If a URI is a bare number (Who's On First ID) it is resolved to its relative path. That path is then prepended
//     with a root "data" folder. Basically it's a short-hand for resolving a specific WOF ID to it's path inside a
//     WOF repo.
//   - If a URI starts with 'repo://' it is assumed to be a whosonfirst/go-whosonfirst-iterate/v2 URI and will be used
//     to create an interator to which every file it encounters will be processed with 'cb'
func ExpandURIsWithCallback(ctx context.Context, cb ExpandURICallbackFunc, uris ...string) error {

	if len(uris) == 0 {
		return nil
	}

	if len(uris) > 1 {

		for _, u := range uris {

			err := ExpandURIsWithCallback(ctx, cb, u)

			if err != nil {
				return err
			}
		}

		return nil
	}

	u := uris[0]

	if strings.HasPrefix(u, "repo://") {

		u = strings.Replace(u, "repo://", "", 1)

		iter, err := iterate.NewIterator(ctx, "repo://")

		if err != nil {
			return err
		}

		for rec, err := range iter.Iterate(ctx, u) {
		    if err != nil {
		       	   return err
	            }

		    err = cb(ctx, rec.Path)

		    if err != nil {
		       return err
		       }
		       }
	}

	if re_wofid.MatchString(u) {

		id, uri_args, err := uri.ParseURI(u)

		if err != nil {
			return fmt.Errorf("Failed to parse URI '%s' in to ID, %w", u, err)
		}

		rel_path, err := uri.Id2RelPath(id, uri_args)

		if err != nil {
			return fmt.Errorf("Failed to derive relative path for ID %d (from URI '%s'), %w", id, u, err)
		}

		u = filepath.Join("data", rel_path)
	}

	abs_u, err := filepath.Abs(u)

	if err != nil {
		return fmt.Errorf("Failed to derive absolute path for '%s', %w", u, err)
	}

	return cb(ctx, abs_u)
}
