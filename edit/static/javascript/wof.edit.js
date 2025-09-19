var wof = wof || {};

wof.edit = (function () {

    var self = {
	
	init: function() {
	    
	    return new Promise((resolve, reject) => {

		sfomuseum.golang.wasm.fetch("wasm/wof_format.wasm").then((rsp) => {
		    
		    sfomuseum.golang.wasm.fetch("wasm/wof_validate.wasm").then((rsp) => {
			
			resolve();
			
		    }).catch((err) => {
			reject("Failed to wof_validate WASM binary " + err);
		    });
		    
		}).catch((err) => {
		    reject("Failed to load update wof_format WASM binary " + err);
		});

		});
	},

	list: function() {

	    const _self = self;
	    
	    wof.edit.api.list().then((rsp) => {
		
		const count = rsp.length;

		if (count == 1){
		    _self.show(rsp[0]);
		    return;
		}
		
		const items = document.createElement("ul");
		
		for (var i=0; i < count; i++) {
		    
		    const uri = rsp[i];
		    
		    const link = document.createElement("a");
		    link.setAttribute("href", uri);
		    link.appendChild(document.createTextNode(uri));
		    
		    link.onclick = function(e){
			const el = e.target;
			const uri = el.getAttribute("href");
			_self.show(uri);
			return false;
		    }
		    
		    const item = document.createElement("li");
		    item.appendChild(link);

		    items.appendChild(item);
		}
		
		const root = document.getElementById("canvas");
		root.innerHTML = "";
		root.appendChild(items);
		
	    }).catch((err) => {
		console.error("Failed to retrieve list to edit", err);
	    });
	    
	},

	show: function(uri) {

	    console.log("SHOW", uri);

	    wof.edit.api.fetch(uri).then((data) => {

		console.log("OK", data);

		const str_data = JSON.stringify(data);

		const map_id = "map";
		
		const map_el = document.createElement("div");
		map_el.setAttribute("id", map_id);

		var raw = document.createElement("pre");
		raw.setAttribute("id", "raw");
		raw.setAttribute("contentEditable", "plaintext-only");
		raw.innerText = str_data;
		
		const left = document.createElement("div");
		left.appendChild(map_el);

		const right = document.createElement("div");
		right.appendChild(raw);
		
		var wrapper = document.createElement("div");
		wrapper.setAttribute("id", "feature");

		wrapper.appendChild(left);
		wrapper.appendChild(right);

		const root = document.getElementById("canvas");
		root.innerHTML = "";
		root.appendChild(wrapper);

		wof_format(str_data).then((fmt_rsp) => {
		    raw.innerText = fmt_rsp;
		}).catch((err) => {
		    console.error("Failed to format response", err);
		})

		const bounds = whosonfirst.geojson.deriveBboxAsBounds(data);
		console.log("BOUNDS", bounds);
		
		const map = L.map(map_id);
		map.fitBounds(bounds);
		
		const osm = L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {});
		osm.addTo(map);

		const feature = L.geoJSON(data);
		feature.addTo(map);

		console.log("WOO");
		
	    }).catch((err) => {
		console.error("SAD", err)
	    });
	},
    };
    
    return self;
})();
