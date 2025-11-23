package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-http-maps/v2"
	opensearch_http "github.com/aaronland/go-http/v4/opensearch" // as in the browser search widget not the document store
	"github.com/whosonfirst/spelunker/v2/http/templates/javascript"
	"github.com/whosonfirst/spelunker/v2/http/templates/text"
	"github.com/whosonfirst/spelunker/v2/http/www"
)

func robotsTxtHandlerFunc(ctx context.Context) (http.Handler, error) {

	t, err := text.LoadTemplates(ctx)

	if err != nil {
		slog.Error("Failed to load text templates", "error", err)
		return nil, fmt.Errorf("Failed to load templates, %w", err)
	}

	opts := &www.RobotsTxtHandlerOptions{
		Templates: t,
	}

	return www.RobotsTxtHandler(opts)
}

func staticHandlerFunc(ctx context.Context) (http.Handler, error) {

	http_fs := http.FS(run_options.StaticAssets)
	fs_handler := http.FileServer(http_fs)

	return http.StripPrefix(run_options.URIs.Static, fs_handler), nil
}

func urisJSHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	js_templates, err := javascript.LoadTemplates(ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to load JavaScript templates, %w", err)
	}

	opts := &www.URIsJSHandlerOptions{
		Templates: js_templates,
		URIs:      uris_table,
	}

	return www.URIsJSHandler(opts)
}

func openSearchHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	search_url, err := uris_table.Abs(uris_table.Search)

	if err != nil {
		slog.Error("Failed to derive search URL", "error", err)
		return nil, fmt.Errorf("Failed to derive search URL")
	}

	// Make me a config? I mean, really, the entire OS description
	// blob should be made configurable so punting for now...

	image_url := `data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAABGdBTUEAALGPC/xhBQAAA5VJREFUWAntVl1IVFEQnlWz8id01cyCCKyIIoiIiLB86OdFi3oIiUCE6MeyoN6MfqCHHsyEqIhASQolgyIhCiGMrAiCCC0wHyKwolJL1NTM1L7Pvcc75+667oPgyw7s7jl35sx8M/Oduesbk/YxmUGJmcHY46GjAKIViLwCubNFWjJFuhaK3PSLZER+NBzRI/NSlCDyJENkdbxIWqxIUaLIfny8ciIpYHcUv6FkCc5m46NkagB7EbwKGcf51DEsX/219zvniFSkimzB75l5to67EoD6mCVyN93ShQeQg4yrETzGE/z1kMgzfIwsRlbVaWYnUtnvrrk6ieBXAI5+lsVZuskBLILT+0Ab7wk+jMF5uNt1MgvLOgRPdVwNjIpc7nP1+ajIxRR3328P3tAACPIOnGbY/Rr3crpH5O2w67AczjeAoEZuIPtOgKCkwT0JqyvY/i+gc75DAziLHuYop+bILTgvU9kVzBU5nmy0IqxOhdITXLoniYY/rj1WwQDWoqalIUjUMChy4Jd7eBXKVInstNQMiHwZCTwh2wtBYC2jAFgLGyU2ALY7FOMbgXr3TxFDfD+O1eNaJqnjY3Be1uu6LkFldOmpYfAP4VpQAMRrwHwtzHxHl8igQx7yg6TLttks9bBrVc43e1rYA16cAn88olKAJg+M1VIHxAw+4ASn7hqu01aPHZ+Xq95zz1ukpRg357PTHnXcBrCSd8qRbiDeh7IrwkspynoQd9orbTB6afoDJedCgrq+bM8ChOIj6p7On/BgA1ihysq+O4DHrY8g8AV1nx+h5EY0sVj5e5gfycq1D5E5Jfke+YRpuNyNo6xwUA+JVpX6MQS/qoKfRy+fq0n4xsme3m6DH+scHrGKWvzIvh1Z5XZOPHWh8FEbSGSGD0mWAuTMuliVvRaz4BzYfh0ZGeEUTIQtb9Ae5+p9ha+NHSKsKl9eWQjeBNCXwJVel1M+608p32YslRH2juUz0oS2bAN6JvwQZc7DIKKMwI5ENWUn4zch+DtVxYBl0Lfdgipk16Ear4O/R9RdXe4syERGRmIBUgfPB8gIgvO4DYClOYTrwomlpRnBt8Npt3re6+kv7b+zv8j8hcMJ7WOStQ2ARg/A7kKM3N8IMISANahKDpx+8wSknZbH2K//IdI8ddn1MZsDWsN7zCr3qay1nmv+U1qK2cEr26huhdcuzH5yAGEOTacquAXT6T0CX1EA0Qr8B2vx8tDeL73aAAAAAElFTkSuQmCC`

	desc_opts := &opensearch_http.BasicOpenSearchDescriptionOptions{
		Name:           "Who's On First Spelunker Search",
		Description:    "Search for places in the Who's On First Spelunker",
		QueryParameter: "q",
		ImageURI:       image_url,
		SearchTemplate: search_url,
		SearchForm:     search_url,
	}

	desc, err := opensearch_http.BasicOpenSearchDescription(desc_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create basic OpenSearch description, %w", err)
	}

	return opensearch_http.OpenSearchHandler(desc)
}

func indexHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.TemplateHandlerOptions{
		Authenticator: authenticator,
		Templates:     html_templates,
		TemplateName:  "index",
		PageTitle:     "",
		URIs:          uris_table,
	}

	return www.TemplateHandler(opts)
}

func aboutHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.TemplateHandlerOptions{
		Authenticator: authenticator,
		Templates:     html_templates,
		TemplateName:  "about",
		PageTitle:     "About",
		URIs:          uris_table,
	}

	return www.TemplateHandler(opts)
}

func descendantsHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.DescendantsHandlerOptions{
		Spelunker:     sp,
		Authenticator: authenticator,
		Templates:     html_templates,
		URIs:          uris_table,
	}

	return www.DescendantsHandler(opts)
}

func recentHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.RecentHandlerOptions{
		Spelunker:     sp,
		Authenticator: authenticator,
		Templates:     html_templates,
		URIs:          uris_table,
	}

	return www.RecentHandler(opts)
}

func idHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.IdHandlerOptions{
		Spelunker:     sp,
		Authenticator: authenticator,
		Templates:     html_templates,
		URIs:          uris_table,
	}

	return www.IdHandler(opts)
}

func placetypesHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.PlacetypesHandlerOptions{
		Spelunker:     sp,
		Authenticator: authenticator,
		Templates:     html_templates,
		URIs:          uris_table,
	}

	return www.PlacetypesHandler(opts)
}

func hasPlacetypeHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.HasPlacetypeHandlerOptions{
		Spelunker:     sp,
		Authenticator: authenticator,
		Templates:     html_templates,
		URIs:          uris_table,
	}

	return www.HasPlacetypeHandler(opts)
}

func hasConcordanceHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.HasConcordanceHandlerOptions{
		Spelunker:     sp,
		Authenticator: authenticator,
		Templates:     html_templates,
		URIs:          uris_table,
	}

	return www.HasConcordanceHandler(opts)
}

func concordancesHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.ConcordancesHandlerOptions{
		Spelunker:     sp,
		Authenticator: authenticator,
		Templates:     html_templates,
		URIs:          uris_table,
	}

	return www.ConcordancesHandler(opts)
}

func searchHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.SearchHandlerOptions{
		Spelunker:     sp,
		Authenticator: authenticator,
		Templates:     html_templates,
		URIs:          uris_table,
	}

	return www.SearchHandler(opts)
}

func nullIslandHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.NullIslandHandlerOptions{
		Spelunker:     sp,
		Authenticator: authenticator,
		Templates:     html_templates,
		URIs:          uris_table,
	}

	return www.NullIslandHandler(opts)
}

func mapConfigHandlers(ctx context.Context) (http.Handler, http.Handler, string, error) {

	opts := &maps.AssignMapConfigHandlerOptions{
		MapProvider:          map_provider,
		MapTileURI:           map_tile_uri,
		ProtomapsTheme:       protomaps_theme,
		ProtomapsMaxDataZoom: protomaps_max_data_zoom,
	}

	map_cfg, err := maps.MapConfigFromOptions(opts)

	if err != nil {
		return nil, nil, "", err
	}

	map_cfg_handler := maps.MapConfigHandler(map_cfg)

	var tile_url_handler http.Handler
	var tile_url string

	if map_cfg.TileURLHandler != nil {
		tile_url_handler = map_cfg.TileURLHandler
		tile_url = map_cfg.TileURL
	}

	return map_cfg_handler, tile_url_handler, tile_url, nil
}
