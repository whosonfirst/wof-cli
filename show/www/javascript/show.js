window.addEventListener("load", function load(event){

    // Null Island
    var map = L.map('map').setView([0.0, 0.0], 12);
    
    var tile_layer = L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
	maxZoom: 19,
    });
    
    tile_layer.addTo(map);

    fetch("/features.geojson")
	.then((rsp) => rsp.json())
	.then((f) => {

	    var raw_el = document.querySelector("#raw");

	    if (raw_el){

		// Something something something WASM
		// https://github.com/whosonfirst/go-whosonfirst-format-wasm
		// Which expects Feature elements rather than FeatureCollections
		// Something something something WASM
		
		var str_f = JSON.stringify(f, "", " ");
		
		var pre = document.createElement("pre");
		pre.appendChild(document.createTextNode(str_f));
		raw_el.appendChild(pre);
	    }

	    var pt_handler = whosonfirst.spelunker.leaflet.handlers.point({});

	    var poly_style = whosonfirst.spelunker.leaflet.styles.consensus_polygon();	    
	    // var lbl_style = whosonfirst.spelunker.leaflet.styles.label_centroid();
	    
	    var geojson_args = {
		style: poly_style,
		pointToLayer: pt_handler,		
	    };
	    
	    var geojson_layer = L.geoJSON(f);
	    geojson_layer.addTo(map);

	    var bounds = whosonfirst.spelunker.geojson.derive_bounds(f);

	    var sw = bounds[0];
	    var ne = bounds[1];

	    if ((sw[0] == ne[0]) && (sw[1] == ne[1])){
		map.setView(sw, 12);
	    } else {
		map.fitBounds(bounds);
	    }
	    
	}).catch((err) => {
	    console.log("SAD", err);
	});
    
});
