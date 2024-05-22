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
