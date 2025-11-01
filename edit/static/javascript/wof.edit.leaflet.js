/**
 * @namespace wof.edit.styles
 * @description Leaflet style and handler definitions for GeoJSON Features. (At some point this should just be merged back in to whosonfirst/js-whosonfirst.
 */

var wof = wof || {};
wof.edit = wof.edit || {};

wof.edit.leaflet = (function () {

    const styles = {
	bbox:  {
	    "color": "#000000",
	    "weight": .5,
	    "opacity": 1,
	    "fillColor": "#000000",
	    "fillOpacity": .4,
	},
	
	label: {
	    "color": "#fff",
	    "weight": 3,
	    "opacity": 1,
	    "radius": 10,
	    "fillColor": "#ff0099",
	    "fillOpacity": 0.8
	},
	
	'math_centroid':  {
	    "color": "#fff",
	    "weight": 2,
	    "opacity": 1,
	    "radius": 6,
	    "fillColor": "#ff7800",
	    "fillOpacity": 0.8
	},

	'geom_centroid': {
	    "color": "#fff",
	    "weight": 3,
	    "opacity": 1,
	    "radius": 10,
	    //"fillColor": "#32cd32",
	    "fillColor": "#1c5894",	    	    
	    "fillOpacity": 0.8
	},

	'search_centroid': {
	    "color": "#000",
	    "weight": 2,
	    "opacity": 1,
	    "radius": 6,
	    "fillColor": "#fe1e9f",
	    // "fillColor": "#0BBDFF",
	    "fillOpacity": 1
	},

	'breach_polygon': {
	    "color": "#ffff00",
	    //"color": "#002EA7",
	    "weight": 1.5,
	    "dashArray": "5, 5",
	    "opacity": 1,
	    "fillColor": "#002EA7",
	    "fillOpacity": 0.1
	},
		
	'consensus_polygon': {
	    "color": "#ff0066",
	    // "color": "#1c5894",	    	    
	    "weight": 2,
	    "opacity": 1,
	    "fillColor": "#ff69b4",
	    "fillOpacity": 0.4
	},

	'parent_polygon': {
	    "color": "#000",
	    "weight": 1,
	    "opacity": 1,
	    "fillColor": "#00308F",
	    "fillOpacity": 0.5
	},
    };
    
    var self = {
	
	point: function(layer_args){
	    
	    return function(feature, latlon){
		const m = L.circleMarker(latlon, layer_args);
		return m;
	    };
	},

	style: function(label){
	    return styles[label];
	},

    };
    
    return self;
    
})();
