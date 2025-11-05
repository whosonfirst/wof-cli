// failed to bind label because TypeError: m.bindLabel is not a function

window.addEventListener("load", function load(event){

    // Because we still don't have a place to send/store these data
    whosonfirst.spelunker.yesnofix.enabled(false);
    
    // START OF wrap me in a webcomponent

    const bbox_pane_name = "bbox";
    const bbox_pane_zindex = 1000;
    
    const parent_pane_name = "parent";
    const parent_pane_zindex = 2000;
    
    const poly_pane_name = "polygon";
    const poly_pane_zindex = 3000;
    
    const centroids_pane_name = "centroids";
    const centroids_pane_zindex = 4000;
    
    try {

	var map_el = document.querySelector("#map");
	var svg_el = document.querySelector("#map-svg");	
	var wof_id = map_el.getAttribute("data-wof-id");

	var alt_source = map_el.getAttribute("data-wof-alt-source");
	var alt_function = map_el.getAttribute("data-wof-alt-function");
	var alt_extra = map_el.getAttribute("data-wof-alt-extra");		

	var uri_args = {
	    alt: true,
	    source: alt_source,
	    function: alt_function,
	    extra: alt_extra,
	};

	whosonfirst.spelunker.feature.fetch(wof_id, uri_args).then((f) => {

	    var props;
	    var pretty;
	    
	    try {

		props = f.properties;
		
		var props_el = document.querySelector("#whosonfirst-properties");
		props_el.innerText = JSON.stringify(props, "", " ");
	
		pretty = whosonfirst.spelunker.properties.render(props);	
	    } catch(err) {
		console.log("Failed to render properties", err);
	    }
	    
	    try {
		var wrapper = document.querySelector("#props-wrapper");
		wrapper.appendChild(pretty);
		
		var raw = wrapper.children[0];
		raw.style.display = "none";
		
		//wrapper.replaceChild(pretty, raw);
		
		var toggle = document.querySelector("#props-toggle");
		toggle.style.display = "inline-block";
		
		var toggle_raw = document.querySelector("#props-toggle-raw");
		toggle_raw.style.display = "block";
		
		toggle_raw.onclick = function(){	    
		    raw.style.display = "block";
		    pretty.style.display = "none";	    
		    toggle_raw.style.display = "none";
		    toggle_pretty.style.display = "block";
		};
		
		var toggle_pretty = document.querySelector("#props-toggle-pretty");
		
		toggle_pretty.onclick = function(){	    
		    raw.style.display = "none";
		    pretty.style.display = "block";	    
		    toggle_raw.style.display = "block";
		    toggle_pretty.style.display = "none";
		};
		
	    } catch(err){
		console.log("Failed to install pretty properties", err);
	    }
	    
	    map_el.style.display = "block";

	    const map = whosonfirst.spelunker.maps.map(map_el);
	    
	    if (f.geometry.type == "Point"){

		var coords = f.geometry.coordinates;
		
		var pt = [ coords[1], coords[0] ];
		var zm = Math.max(12, f.properties["mz:min_zoom"]);
		map.setView(pt, zm);
		
	    } else {
		var bounds = whosonfirst.spelunker.geojson.derive_bounds(f);
		map.fitBounds(bounds);
	    }
	    
	    // http://localhost:8080/id/1259472055
	    if (f.geometry.type == "Point"){

		var pt_handler_layer_args = {
		    pane: whosonfirst.spelunker.maps.centroids_pane_name,
		    tooltips_pane: whosonfirst.spelunker.maps.tooltips_pane_name,
		};
		
		var pt_handler = whosonfirst.spelunker.leaflet.handlers.point(pt_handler_layer_args);
		var lbl_style = whosonfirst.spelunker.leaflet.styles.label_centroid();

		var layer_args = {
		    style: lbl_style,
		    pointToLayer: pt_handler,
		    pane: centroids_pane_name,
		}
		
		whosonfirst.spelunker.leaflet.draw_point(map, f, layer_args);
		return;
	    }

	    var bbox_style = whosonfirst.spelunker.leaflet.styles.bbox();

	    var bbox_layer_args = {
		style: bbox_style,
		pane: bbox_pane_name,
	    }
	    
	    whosonfirst.spelunker.leaflet.draw_bbox(map, f, bbox_layer_args);
	    
	    var poly_style = whosonfirst.spelunker.leaflet.styles.consensus_polygon();
	    
	    var poly_layer_args = {
		style: poly_style,
		pane: poly_pane_name,
	    };
	    
	    whosonfirst.spelunker.leaflet.draw_poly(map, f, poly_layer_args);

	    var props = f.properties;

	    var pt_handler_layer_args = {
		pane: whosonfirst.spelunker.maps.centroids_pane_name,
		tooltips_pane: whosonfirst.spelunker.maps.tooltips_pane_name,
	    };
	    
	    var pt_handler = whosonfirst.spelunker.leaflet.handlers.point(pt_handler_layer_args);
	    
	    if ((props["lbl:longitude"]) && (props["lbl:latitude"])){
		
		var lbl_centroid = [ props["lbl:longitude"], props["lbl:latitude" ] ];
		
		var lbl_f = {
		    "type": "Feature",
		    "properties": { "lflt:label_text": "label centroid" },
		    "geometry": { "type": "Point", "coordinates": lbl_centroid }
		};
		
		var lbl_style = whosonfirst.spelunker.leaflet.styles.label_centroid();
		
		var lbl_layer_args = {
		    style: lbl_style,
		    pointToLayer: pt_handler,
		    pane: centroids_pane_name,
		};
		
		whosonfirst.spelunker.leaflet.draw_point(map, lbl_f, lbl_layer_args);		
	    }

	    if ((props["geom:longitude"]) && (props["geom:latitude"])){
		
		var math_centroid = [ props["geom:longitude"], props["geom:latitude" ] ];
		
		var math_f = {
		    "type": "Feature",
		    "properties": { "lflt:label_text": "math centroid" },
		    "geometry": { "type": "Point", "coordinates": math_centroid }
		};	    
		
		var math_style = whosonfirst.spelunker.leaflet.styles.math_centroid();
		
		var math_layer_args = {
		    style: math_style,
		    pointToLayer: pt_handler,
		    pane: centroids_pane_name,
		};
		
		whosonfirst.spelunker.leaflet.draw_point(map, math_f, math_layer_args);
	    }
	    
	    // Draw parent here...

	    var parent_id = f.properties["wof:parent_id"];

	    if ((parent_id) && (parent_id > 0)){
		
		// console.log("Fetch parent", parent_id);
		
		whosonfirst.spelunker.feature.fetch(parent_id).then((parent_f) => {
		    
		    if (! parent_f.geometry.type.endsWith("Polygon")){
			return;
		    }
		    
		    var parent_style = whosonfirst.spelunker.leaflet.styles.parent_polygon();
		    
		    var parent_layer_args = {
			style: parent_style,
			pane: parent_pane_name,
		    };
		    
		    whosonfirst.spelunker.leaflet.draw_poly(map, parent_f, parent_layer_args);
		    
		}).catch((err) => {
		    console.log("Failed to fetch parent record", parent_id, err);
		})
	    }

	    document.querySelector("#whosonfirst-alt-properties").style.display = "grid";
	    
	}).catch((err) => {
	    console.log("Failed to initialize map", err);
	    throw(err);
	});
	
    } catch (err) {
	console.log("Failed to initialize map", err);
	svg_el.style.display = "block";	    	
    }
    
    // END OF wrap me in a webcomponent    
    
    whosonfirst.spelunker.namify.namify_selector(".props-uoc");
    whosonfirst.spelunker.namify.namify_selector(".wof-namify");
    whosonfirst.spelunker.namify.namify_selector(".yesnofix-uoc");
    
});
