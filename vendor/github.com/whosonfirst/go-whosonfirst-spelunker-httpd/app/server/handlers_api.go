package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd/api"
)

func findingAidHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAPIOnce.Do(setupAPI)

	if setupAPIError != nil {
		slog.Error("Failed to set up common configuration", "error", setupAPIError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupAPIError)
	}

	opts := &api.FindingAidHandlerOptions{
		Spelunker: sp,
	}

	h, err := api.FindingAidHandler(opts)

	if err != nil {
		return nil, err
	}

	return cors_wrapper.Handler(h), nil
}

func geoJSONHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAPIOnce.Do(setupAPI)

	if setupAPIError != nil {
		slog.Error("Failed to set up common configuration", "error", setupAPIError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupAPIError)
	}

	opts := &api.GeoJSONHandlerOptions{
		Spelunker: sp,
	}

	h, err := api.GeoJSONHandler(opts)

	if err != nil {
		return nil, err
	}

	return cors_wrapper.Handler(h), nil
}

func geoJSONLDHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAPIOnce.Do(setupAPI)

	if setupAPIError != nil {
		slog.Error("Failed to set up common configuration", "error", setupAPIError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupAPIError)
	}

	opts := &api.GeoJSONLDHandlerOptions{
		Spelunker: sp,
	}

	h, err := api.GeoJSONLDHandler(opts)

	if err != nil {
		return nil, err
	}

	return cors_wrapper.Handler(h), nil
}

func sprHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAPIOnce.Do(setupAPI)

	if setupAPIError != nil {
		slog.Error("Failed to set up common configuration", "error", setupAPIError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupAPIError)
	}

	opts := &api.SPRHandlerOptions{
		Spelunker: sp,
	}

	h, err := api.SPRHandler(opts)

	if err != nil {
		return nil, err
	}

	return cors_wrapper.Handler(h), nil
}

func selectHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAPIOnce.Do(setupAPI)

	if setupAPIError != nil {
		slog.Error("Failed to set up common configuration", "error", setupAPIError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupAPIError)
	}

	// Make this a config/flag
	select_pattern := `properties(?:.[a-zA-Z0-9-_]+){1,}`

	pat, err := regexp.Compile(select_pattern)

	if err != nil {
		slog.Error("Failed to compile select pattern", "pattern", select_pattern, "error", err)
		return nil, fmt.Errorf("Failed to compile select pattern (%s), %w", select_pattern, err)
	}

	opts := &api.SelectHandlerOptions{
		Pattern:   pat,
		Spelunker: sp,
	}

	h, err := api.SelectHandler(opts)

	if err != nil {
		return nil, err
	}

	return cors_wrapper.Handler(h), nil
}

func navPlaceHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAPIOnce.Do(setupAPI)

	if setupAPIError != nil {
		slog.Error("Failed to set up common configuration", "error", setupAPIError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupAPIError)
	}

	opts := &api.NavPlaceHandlerOptions{
		Spelunker:   sp,
		MaxFeatures: 10,
	}

	h, err := api.NavPlaceHandler(opts)

	if err != nil {
		return nil, err
	}

	return cors_wrapper.Handler(h), nil
}

func svgHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAPIOnce.Do(setupAPI)

	if setupAPIError != nil {
		slog.Error("Failed to set up common configuration", "error", setupAPIError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupAPIError)
	}

	sz := api.DefaultSVGSizes()

	opts := &api.SVGHandlerOptions{
		Spelunker: sp,
		Sizes:     sz,
	}

	return api.SVGHandler(opts)
}

func descendantsFacetedHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAPIOnce.Do(setupAPI)

	if setupAPIError != nil {
		slog.Error("Failed to set up common configuration", "error", setupAPIError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupAPIError)
	}

	opts := &api.DescendantsFacetedHandlerOptions{
		Spelunker: sp,
		// Authenticator: authenticator,
	}

	return api.DescendantsFacetedHandler(opts)
}

func placetypeFacetedHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAPIOnce.Do(setupAPI)

	if setupAPIError != nil {
		slog.Error("Failed to set up common configuration", "error", setupAPIError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupAPIError)
	}

	opts := &api.PlacetypeFacetedHandlerOptions{
		Spelunker: sp,
		// Authenticator: authenticator,
	}

	return api.PlacetypeFacetedHandler(opts)
}

func recentFacetedHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAPIOnce.Do(setupAPI)

	if setupAPIError != nil {
		slog.Error("Failed to set up common configuration", "error", setupAPIError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupAPIError)
	}

	opts := &api.RecentFacetedHandlerOptions{
		Spelunker: sp,
		// Authenticator: authenticator,
	}

	return api.RecentFacetedHandler(opts)
}

func hasConcordanceFacetedHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAPIOnce.Do(setupAPI)

	if setupAPIError != nil {
		slog.Error("Failed to set up common configuration", "error", setupAPIError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupAPIError)
	}

	opts := &api.HasConcordanceFacetedHandlerOptions{
		Spelunker: sp,
		// Authenticator: authenticator,
	}

	return api.HasConcordanceFacetedHandler(opts)
}

func searchFacetedHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAPIOnce.Do(setupAPI)

	if setupAPIError != nil {
		slog.Error("Failed to set up common configuration", "error", setupAPIError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupAPIError)
	}

	opts := &api.SearchFacetedHandlerOptions{
		Spelunker: sp,
		// Authenticator: authenticator,
	}

	return api.SearchFacetedHandler(opts)
}

func nullIslandFacetedHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAPIOnce.Do(setupAPI)

	if setupAPIError != nil {
		slog.Error("Failed to set up common configuration", "error", setupAPIError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupAPIError)
	}

	opts := &api.NullIslandFacetedHandlerOptions{
		Spelunker: sp,
		// Authenticator: authenticator,
	}

	return api.NullIslandFacetedHandler(opts)
}
