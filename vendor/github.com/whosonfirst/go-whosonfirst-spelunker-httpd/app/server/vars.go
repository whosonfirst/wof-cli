package server

import (
	html_template "html/template"
	"sync"

	"github.com/rs/cors"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
)

var run_options *RunOptions

var sp spelunker.Spelunker

var authenticator auth.Authenticator

var uris_table *httpd.URIs

var html_templates *html_template.Template

var setupCommonOnce sync.Once
var setupCommonError error

var setupWWWOnce sync.Once
var setupWWWError error

var setupAPIOnce sync.Once
var setupAPIError error

var cors_wrapper *cors.Cors
