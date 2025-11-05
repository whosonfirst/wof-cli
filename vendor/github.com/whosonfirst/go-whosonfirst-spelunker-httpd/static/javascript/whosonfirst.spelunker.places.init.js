window.addEventListener("load", function load(event){

    var places = document.querySelectorAll(".whosonfirst-places-list li");

    if (! places){
	console.log("No places");
	return;
    }
    
    var count_places = places.length;

    // console.log("Count", count_places);
    
    if (count_places == 0){
	return;
    }
    
    var coords = [];
    var names = [];
    var links = [];
    
    for (var i=0; i < count_places; i++) {

	var el = places[i];
	// console.log(el);
	
	var lat = parseFloat(el.getAttribute("data-latitude"));
	var lon = parseFloat(el.getAttribute("data-longitude"));
	var id = parseInt(el.getAttribute("data-id"));		

	if ((! lat) || (!lon)){
	    console.log("Invalid coordinates", i, lat, lon);
	    continue;
	}

	if (! id){
	    console.log("Invalid ID", i, id);
	}
	
	var n = el.querySelector(".wof-place-name");

	if ((! n) || (n.innerText == "")){
	    console.log("Invalid name", i);
	    continue;
	}

	coords.push([ lon, lat ]);
	names[ JSON.stringify(coords[i]) ] = n.innerText + " (" + id + ")";
	links[ JSON.stringify(coords[i]) ] = n.getAttribute("href");	
    }

    // console.log("Coords", coords);
    // console.log("Names", names);
    
    var f = {
	"type": "Feature",
	// These get handled in whosonfirst.spelunker.leaflet.handlers.js (point)
	"properties": {
	    "lflt:label_names": names,
	    "lflt:label_links": links,
	},
	"geometry": {
	    "type": "MultiPoint",
	    "coordinates": coords,
	},
    };

    // console.log("Feature", f);
    
    var map_el = document.querySelector("#map");
    map_el.style.display = "block";
    
    const map = whosonfirst.spelunker.maps.map(map_el);

    switch (coords.length){
	case 0:
	// Null Island
	    coords.push([ 0.0, 0.0 ]);
	    f.geometry.coords = coords;
	    map.setView([coords[0][1], coords[0][0]], 3);
	    break;
	case 1:
	    // TO DO: set zoom based on placetype or mz:min/max_zoom (requires fetching the record...)
	    map.setView([coords[0][1], coords[0][0]], 12);
	    break;
	default:

	    // START OF wrap me in a common function
	    
	    var bounds = whosonfirst.spelunker.geojson.derive_bounds(f);
	    var sw = bounds[0];
	    var ne = bounds[1];

	    // TO DO: set zoom based on placetype or mz:min/max_zoom (requires fetching all the records so... maybe not?)	    
	    if ((sw[0] == ne[0]) && (sw[1] == ne[1])){
		map.setView(sw, 12);
	    } else {
		map.fitBounds(bounds);
	    }

	    // END OF wrap me in a common function
	    
	    break;
    }
    
    var pt_handler_layer_args = {
	pane: whosonfirst.spelunker.maps.centroids_pane_name,
	tooltips_pane: whosonfirst.spelunker.maps.tooltips_pane_name,	
    };
		
    var pt_handler = whosonfirst.spelunker.leaflet.handlers.point(pt_handler_layer_args);
    var lbl_style = whosonfirst.spelunker.leaflet.styles.search_centroid();

    var points_layer_args = {
	style: lbl_style,
	pointToLayer: pt_handler,
	pane: whosonfirst.spelunker.maps.centroids_pane_name,
    }
    
    var points_layer = L.geoJSON(f, points_layer_args);
    points_layer.addTo(map);

});
