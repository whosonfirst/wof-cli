window.addEventListener("load", function load(event){
    
    var map = L.map('map').setView([51.505, -0.09], 13);
    
    var tile_layer = L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
	maxZoom: 19,
	attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
    });
    
    tile_layer.addTo(map);

    fetch("/features.geojson")
	.then((rsp) => rsp.json())
	.then((f) => {
	    console.log(f)

	    var geojson_layer = L.geoJSON(f);
	    geojson_layer.addTo(map);

	    var bounds = whosonfirst.spelunker.geojson.derive_bounds(f);
	    map.fitBounds(bounds);
	    
	}).catch((err) => {
	    console.log("SAD", err);
	});
    
});
