package http

import (
	_ "fmt"
	"github.com/aaronland/go-http-rewrite"
	"github.com/sfomuseum/go-http-opensearch"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	_ "log"
	"net/http"
)

type AppendPluginsOptions struct {
	Plugins map[string]*opensearch.OpenSearchDescription
}

func AppendPluginsHandler(next http.Handler, opts *AppendPluginsOptions) http.Handler {

	var cb rewrite.RewriteHTMLFunc

	cb = func(n *html.Node, w io.Writer) {

		if n.Type == html.ElementNode && n.Data == "head" {

			for uri, d := range opts.Plugins {

				link_rel := html.Attribute{"", "rel", "search"}
				link_type := html.Attribute{"", "type", "application/opensearchdescription+xml"}
				link_href := html.Attribute{"", "href", uri}
				link_title := html.Attribute{"", "title", d.ShortName}

				link := html.Node{
					Type:      html.ElementNode,
					DataAtom:  atom.Link,
					Data:      "link",
					Namespace: "",
					Attr: []html.Attribute{
						link_rel,
						link_type,
						link_href,
						link_title,
					},
				}

				n.AppendChild(&link)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			cb(c, w)
		}
	}

	return rewrite.RewriteHTMLHandler(next, cb)
}
