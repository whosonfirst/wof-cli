package opensearch

import (
	"encoding/xml"
)

type OpenSearchImage struct {
	Height int    `xml:"height,attr"`
	Width  int    `xml:"width,attr"`
	URI    string `xml:",chardata"`
}

type OpenSearchURL struct {
	Type       string                    `xml:"type,attr"`
	Method     string                    `xml:"method,attr"`
	Template   string                    `xml:"template,attr"`
	Parameters []*OpenSearchURLParameter `xml:"Param"`
}

type OpenSearchURLParameter struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type OpenSearchDescription struct {
	XMLName       xml.Name         `xml:"OpenSearchDescription"`
	NSMoz         string           `xml:"xmlns:moz,attr"`
	InputEncoding string           `xml:"InputEncoding"`
	NSOpenSearch  string           `xml:"xmlns,attr"`
	ShortName     string           `xml:"ShortName"`
	Description   string           `xml:"Description"`
	Image         *OpenSearchImage `xml:"Image"`
	URL           *OpenSearchURL   `xml:"Url"`
	SearchForm    string           `xml:"moz:searchForm"`
}

func (d *OpenSearchDescription) Marshal() ([]byte, error) {

	enc, err := xml.Marshal(d)

	if err != nil {
		return nil, err
	}

	body := []byte(xml.Header)
	body = append(body, enc...)

	return body, nil
}
