package update

// Move this in to go-whosonfirst-export... TBD

import (
	"context"

	"github.com/paulmach/orb/geojson"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type UpdateFeatureOptions struct {
	StringProperties  multi.KeyValueString
	Int64Properties   multi.KeyValueInt64
	Float64Properties multi.KeyValueFloat64
	BooleanProperties multi.KeyValueBool
	Geometry          *geojson.Geometry
	IfMissing         bool
}

func UpdateFeature(ctx context.Context, body []byte, opts *UpdateFeatureOptions) (bool, []byte, error) {

	changed := false

	for _, p := range opts.StringProperties {

		path := p.Key()
		new_value := p.Value()

		update := true

		old_rsp := gjson.GetBytes(body, path)

		if old_rsp.Exists() {

			old_value := old_rsp.String()

			if old_value == new_value {
				update = false
			}

			if update && opts.IfMissing {
				update = false
			}
		}

		if update {

			new_body, err := sjson.SetBytes(body, path, new_value)

			if err != nil {
				return false, nil, err
			}

			body = new_body
			changed = true
		}
	}

	for _, p := range opts.Int64Properties {

		path := p.Key()
		new_value := p.Value()

		update := true

		old_rsp := gjson.GetBytes(body, path)

		if old_rsp.Exists() {

			old_value := old_rsp.Int()

			if old_value == new_value {
				update = false
			}

			if update && opts.IfMissing {
				update = false
			}
		}

		if update {

			new_body, err := sjson.SetBytes(body, path, new_value)

			if err != nil {
				return false, nil, err
			}

			body = new_body
			changed = true
		}
	}

	for _, p := range opts.Float64Properties {

		path := p.Key()
		new_value := p.Value()

		update := true

		old_rsp := gjson.GetBytes(body, path)

		if old_rsp.Exists() {

			old_value := old_rsp.Float()

			if old_value == new_value {
				update = false
			}

			if update && opts.IfMissing {
				update = false
			}
		}

		if update {

			new_body, err := sjson.SetBytes(body, path, new_value)

			if err != nil {
				return false, nil, err
			}

			body = new_body
			changed = true
		}
	}

	for _, p := range opts.BooleanProperties {

		path := p.Key()
		new_value := p.Value()

		update := true

		old_rsp := gjson.GetBytes(body, path)

		if old_rsp.Exists() {

			old_value := old_rsp.Bool()

			if old_value == new_value {
				update = false
			}

			if update && opts.IfMissing {
				update = false
			}
		}

		if update {

			new_body, err := sjson.SetBytes(body, path, new_value)

			if err != nil {
				return false, nil, err
			}

			body = new_body
			changed = true
		}
	}

	if opts.Geometry != nil {

		new_body, err := sjson.SetBytes(body, "geometry", opts.Geometry)

		if err != nil {
			return false, nil, err
		}

		body = new_body
		changed = true
	}

	return changed, body, nil
}
