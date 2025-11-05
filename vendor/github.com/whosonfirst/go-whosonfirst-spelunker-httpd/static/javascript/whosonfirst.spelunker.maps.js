var whosonfirst = whosonfirst || {};
whosonfirst.spelunker = whosonfirst.spelunker || {};

whosonfirst.spelunker.maps = (function(){

    const bbox_pane_name = "bbox";
    const bbox_pane_zindex = 1000;
    
    const parent_pane_name = "parent";
    const parent_pane_zindex = 2000;
    
    const poly_pane_name = "polygon";
    const poly_pane_zindex = 3000;
    
    const centroids_pane_name = "centroids";
    const centroids_pane_zindex = 4000;

    const tooltips_pane_name = "tooltips";
    const tooltips_pane_zindex = 4001;
    
    var _maps = {};
    
    var self = {

	map: function(map_el){

	    const map_id = map_el.getAttribute("id");

	    if (_maps[map_id]){
		return maps[map_id];
	    }
	    
	    var tiles_url = map_el.getAttribute("data-tiles-url");	    
	    tiles_url = decodeURIComponent(tiles_url);

	    const map = L.map(map_el);

	    var bbox_pane = map.createPane(bbox_pane_name);
	    bbox_pane.style.zIndex = bbox_pane_zindex;
	    
	    var parent_pane = map.createPane(parent_pane_name);
	    parent_pane.style.zIndex = parent_pane_zindex;
	    
	    var poly_pane = map.createPane(poly_pane_name);
	    poly_pane.style.zIndex = poly_pane_zindex;
	    
	    var centroids_pane = map.createPane(centroids_pane_name);
	    centroids_pane.style.zIndex = centroids_pane_zindex;

	    var tooltips_pane = map.createPane(tooltips_pane_name);
	    tooltips_pane.style.zIndex = tooltips_pane_zindex;
	    
	    var layer = protomapsL.leafletLayer({url: tiles_url, theme: 'white'});
	    layer.addTo(map);

	    _maps[map_id] = map;	    
	    return map;
	},

	bbox_pane_name: bbox_pane_name,
	parent_pane_name: parent_pane_name,
	poly_pane_name: poly_pane_name,
	centroids_pane_name: centroids_pane_name,
	tooltips_pane_name: tooltips_pane_name,
    };

    return self;
})();

	
