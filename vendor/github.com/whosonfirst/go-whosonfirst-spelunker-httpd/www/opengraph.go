package www

// https://ogp.me/

// OpenGraph describes a struct containing OpenGraph metadata to pass down to
// templates and include as HTML <meta> tags
type OpenGraph struct {
	Type        string
	SiteName    string
	Title       string
	Description string
	Image       string
}
