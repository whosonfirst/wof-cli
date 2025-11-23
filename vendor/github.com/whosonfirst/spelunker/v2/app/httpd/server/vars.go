package server

import (
	html_template "html/template"
	"sync"

	"github.com/aaronland/go-http/v4/auth"
	"github.com/rs/cors"
	"github.com/whosonfirst/go-whosonfirst-derivatives"
	"github.com/whosonfirst/spelunker/v2"
	wof_http "github.com/whosonfirst/spelunker/v2/http"
)

var run_options *RunOptions

var sp spelunker.Spelunker

var pr derivatives.Provider

var authenticator auth.Authenticator

var uris_table *wof_http.URIs

var html_templates *html_template.Template

var setupCommonOnce sync.Once
var setupCommonError error

var setupWWWOnce sync.Once
var setupWWWError error

var setupAPIOnce sync.Once
var setupAPIError error

var cors_wrapper *cors.Cors
