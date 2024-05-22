package show

import (
	"net/http"

	"github.com/paulmach/orb/geojson"
)

func dataHandler(fc *geojson.FeatureCollection) http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		enc_json, err := fc.MarshalJSON()

		if err != nil {
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-type", "application/json")
		rsp.Write(enc_json)
		return
	}

	return http.HandlerFunc(fn)
}
