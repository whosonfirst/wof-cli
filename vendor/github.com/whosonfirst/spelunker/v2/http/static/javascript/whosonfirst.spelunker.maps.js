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

	    return  new Promise((resolve, reject) => {
		
		const map_id = map_el.getAttribute("id");
		
		if (_maps[map_id]){
		    resolve(maps[map_id]);
		    return;
		}

		console.debug("Fetch map config");
		
		fetch("/maps.json").then(rsp =>
		    rsp.json()
		).then((cfg) => {

		    console.debug("Retrieved map config", cfg);
		    
		    const map = L.map(map_el);
		    
		    const bbox_pane = map.createPane(bbox_pane_name);
		    bbox_pane.style.zIndex = bbox_pane_zindex;
		    
		    const parent_pane = map.createPane(parent_pane_name);
		    parent_pane.style.zIndex = parent_pane_zindex;
		    
		    const poly_pane = map.createPane(poly_pane_name);
		    poly_pane.style.zIndex = poly_pane_zindex;
		    
		    const centroids_pane = map.createPane(centroids_pane_name);
		    centroids_pane.style.zIndex = centroids_pane_zindex;
		    
		    const tooltips_pane = map.createPane(tooltips_pane_name);
		    tooltips_pane.style.zIndex = tooltips_pane_zindex;

		    console.debug("Configure map provider", cfg.provider);
		    
		    switch (cfg.provider) {
			case "leaflet":
			    
			    var tile_url = cfg.tile_url;
			    
			    var tile_layer = L.tileLayer(tile_url, {
				maxZoom: 19,
			    });
			    
			    tile_layer.addTo(map);
			    break;
			    
			case "protomaps":
			    
			    var tile_url = cfg.tile_url;

			    var pm_args = {
				url: tile_url,
				theme: cfg.protomaps.theme,
			    };

			    if ("max_data_zoom" in cfg){
				pm_args.maxDataZoom = cfg.max_data_zoom;
			    }
			    
			    var tile_layer = protomapsL.leafletLayer(pm_args)
			    tile_layer.addTo(map);
			    break;
			    
			default:
			    reject("Unsupported map tile provider");
			    return;
		    }
		    
		    _maps[map_id] = map;	    
		    resolve(map);
		    
		}).catch((err) => {
		    console.error("Invalid or unsupported map provider", err);
		    reject(err);
		});
	    });
	},
	
	bbox_pane_name: bbox_pane_name,
	parent_pane_name: parent_pane_name,
	poly_pane_name: poly_pane_name,
	centroids_pane_name: centroids_pane_name,
	tooltips_pane_name: tooltips_pane_name,
    };

    return self;
})();

	
