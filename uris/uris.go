package uris

import (
	"context"
	"fmt"
	"io"
	_ "os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

// To do : alt files...
var re_wofid = regexp.MustCompile(`^\d+$`)

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

		iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {
			return cb(ctx, path)
		}

		iter, err := iterator.NewIterator(ctx, "repo://", iter_cb)

		if err != nil {
			return err
		}

		return iter.IterateURIs(ctx, u)
	}

	if re_wofid.MatchString(u) {

		id, err := strconv.ParseInt(u, 10, 64)

		if err != nil {
			return fmt.Errorf("Failed to parse URI '%s' in to ID, %w", u, err)
		}

		rel_path, err := uri.Id2RelPath(id)

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
