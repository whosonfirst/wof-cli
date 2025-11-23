package opensearch

type BasicOpenSearchDescriptionOptions struct {
	Name           string
	Description    string
	QueryParameter string
	SearchTemplate string
	SearchForm     string
	ImageURI       string
}

func BasicOpenSearchDescription(opts *BasicOpenSearchDescriptionOptions) (*OpenSearchDescription, error) {

	im := &OpenSearchImage{
		Height: DEFAULT_IMAGE_HEIGHT,
		Width:  DEFAULT_IMAGE_WIDTH,
		URI:    opts.ImageURI,
	}

	params := []*OpenSearchURLParameter{
		&OpenSearchURLParameter{
			Name:  opts.QueryParameter,
			Value: DEFAULT_SEARCHTERMS,
		},
	}

	u := &OpenSearchURL{
		Type:       DEFAULT_URL_TYPE,
		Method:     DEFAULT_URL_METHOD,
		Template:   opts.SearchTemplate,
		Parameters: params,
	}

	desc := &OpenSearchDescription{
		NSMoz:         NS_MOZILLA_SEARCH,
		NSOpenSearch:  NS_OPENSEARCH,
		InputEncoding: "UTF-8",
		ShortName:     opts.Name,
		Description:   opts.Description,
		Image:         im,
		URL:           u,
		SearchForm:    opts.SearchForm,
	}

	return desc, nil
}
