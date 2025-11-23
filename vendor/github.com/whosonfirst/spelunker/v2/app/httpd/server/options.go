package server

import (
	"context"
	"flag"
	"fmt"
	html_template "html/template"
	io_fs "io/fs"
	"net/url"

	"github.com/aaronland/go-http/v4/route"
	"github.com/mitchellh/copystructure"
	"github.com/sfomuseum/go-flags/flagset"
	sfom_funcs "github.com/sfomuseum/go-template/funcs"
	wof_http "github.com/whosonfirst/spelunker/v2/http"
	"github.com/whosonfirst/spelunker/v2/http/static"
	wof_funcs "github.com/whosonfirst/spelunker/v2/http/templates/funcs"
	"github.com/whosonfirst/spelunker/v2/http/templates/html"
)

type RunOptions struct {
	ServerURI         string                            `json:"server_uri"`
	SpelunkerURI      string                            `json:"spelunker_uri"`
	AuthenticatorURI  string                            `json:"authenticator_uri"`
	URIs              *wof_http.URIs                    `json:"uris"`
	HTMLTemplates     []io_fs.FS                        `json:"templates,omitemtpy"`
	HTMLTemplateFuncs html_template.FuncMap             `json:"template_funcs,omitempty"`
	StaticAssets      io_fs.FS                          `json:"static_assets,omitempty"`
	CustomHandlers    map[string]route.RouteHandlerFunc `json:"custom_handlers,omitempty"`
	Verbose           bool                              `json:"verbose"`
}

func (o *RunOptions) Clone() (*RunOptions, error) {

	v, err := copystructure.Copy(o)

	if err != nil {
		return nil, fmt.Errorf("Failed to create local run options, %w", err)
	}

	new_opts := v.(*RunOptions)

	// Things that aren't handled by copystructure
	// TBD...

	new_opts.HTMLTemplates = o.HTMLTemplates
	new_opts.HTMLTemplateFuncs = o.HTMLTemplateFuncs
	new_opts.StaticAssets = o.StaticAssets

	return new_opts, nil
}

func RunOptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "SPELUNKER")

	if err != nil {
		return nil, fmt.Errorf("Failed to assign flags from environment variables, %w", err)
	}

	return RunOptionsFromParsedFlags(ctx)
}

func RunOptionsFromParsedFlags(ctx context.Context, args ...string) (*RunOptions, error) {

	if root_url == "" {
		root_url = server_uri
	}

	root_u, err := url.Parse(root_url)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse root_url '%s', %w", root_url, err)
	}

	uris_table = wof_http.DefaultURIs()
	uris_table.RootURL = root_u.String()

	t_funcs := html_template.FuncMap{
		"IsAvailable": sfom_funcs.IsAvailable,
		// "Add":              sfom_funcs.Add,
		"JoinPath": sfom_funcs.JoinPath,
		// "QRCodeB64":        sfom_funcs.QRCodeB64,
		// "QRCodeDataURI":    sfom_funcs.QRCodeDataURI,
		// "IsEven":           sfom_funcs.IsEven,
		// "IsOdd":            sfom_funcs.IsOdd,
		"FormatStringTime": sfom_funcs.FormatStringTime,
		"FormatUnixTime":   sfom_funcs.FormatUnixTime,
		"GjsonGet":         sfom_funcs.GjsonGet,
		// https://github.com/golang/go/issues/57773
		"URIForId":         wof_http.URIForIdSimple,
		"URIForRecent":     wof_http.URIForRecentSimple,
		"NameForSource":    wof_funcs.NameForSource,
		"FormatNumber":     wof_funcs.FormatNumber,
		"AppendPagination": wof_funcs.AppendPagination,
		"IsAPlacetype":     wof_funcs.IsAPlacetype,
	}

	opts := &RunOptions{
		ServerURI:         server_uri,
		AuthenticatorURI:  authenticator_uri,
		SpelunkerURI:      spelunker_uri,
		URIs:              uris_table,
		HTMLTemplates:     []io_fs.FS{html.FS},
		HTMLTemplateFuncs: t_funcs,
		StaticAssets:      static.FS,
		Verbose:           verbose,
	}

	return opts, nil
}
