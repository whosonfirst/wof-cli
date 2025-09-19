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

		// Set up buttons and feedback
		
		const validate_btn = document.createElement("button");
		validate_btn.setAttribute("class", "btn btn-primary");
		validate_btn.appendChild(document.createTextNode("Validate"));

		const format_btn = document.createElement("button");
		format_btn.setAttribute("class", "btn btn-primary");
		format_btn.appendChild(document.createTextNode("Format"));

		const save_btn = document.createElement("button");
		save_btn.setAttribute("class", "btn btn-primary");
		save_btn.appendChild(document.createTextNode("Save"));

		const buttons = document.createElement("div");
		buttons.setAttribute("class", "btn-group");
		buttons.setAttribute("id", "buttons");
		
		buttons.appendChild(format_btn);		
		buttons.appendChild(validate_btn);
		buttons.appendChild(save_btn);

		// Set up UI
		
		const feedback = document.createElement("div");
		feedback.setAttribute("id", "feedback");
		
		const ui = document.createElement("div");
		ui.appendChild(feedback);				
		ui.appendChild(buttons);		
		ui.appendChild(wrapper);
		
		
		const root = document.getElementById("canvas");
		root.innerHTML = "";
		root.appendChild(ui);

		wof_format(str_data).then((fmt_rsp) => {
		    raw.innerText = fmt_rsp;
		}).catch((err) => {
		    console.error("Failed to format response", err);
		})

		// Set up map
		
		const bounds = whosonfirst.geojson.deriveBboxAsBounds(data);
		
		const map = L.map(map_id);
		map.fitBounds(bounds);
		
		const osm = L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {});
		osm.addTo(map);

		const feature = L.geoJSON(data);
		feature.addTo(map);

		// Set up button interactions

		format_btn.onclick = function(){

		    var data;

		    try {
			data = JSON.parse(raw.innerText);
		    } catch(err) {
			console.error("Failed to parse data", err);
			save_btn.setAttribute("disabled", "disabled");
			return;
		    }
		    
		    const str_data = JSON.stringify(data);
		    
		    wof_format(str_data).then((fmt_rsp) => {
			raw.innerText = fmt_rsp;
			save_btn.removeAttribute("disabled");
		    }).catch((err) => {
			console.error("Failed to format data", err)
			save_btn.setAttribute("disabled", "disabled");
		    });
		};

		validate_btn.onclick = function() {

		    var data;

		    try {
			data = JSON.parse(raw.innerText);
		    } catch(err) {
			console.error("Failed to parse data", err);
			save_btn.setAttribute("disabled", "disabled");			
			return;
		    }
		    
		    const str_data = JSON.stringify(data);
		    
		    wof_validate(str_data).then(() => {
			console.log("OK")
			save_btn.removeAttribute("disabled");

			wof_format(str_data).then((fmt_rsp) => {
			    raw.innerText = fmt_rsp;			    
			}).catch((err) => {
			    console.error("Failed to format data", err);
			    save_btn.setAttribute("disabled", "disabled");						    
			});
			
		    }).catch((err) => {
			console.error("Failed to validate data", err);
			save_btn.setAttribute("disabled", "disabled");						
		    });		    
		};

		save_btn.onclick = function(){

		    var data;

		    try {
			data = JSON.parse(raw.innerText);
		    } catch(err) {
			save_btn.setAttribute("disabled", "disabled");			
			console.error("Failed to parse data", err);
			return;
		    }
		    
		    const str_data = JSON.stringify(data);
		    
		    wof_validate(str_data).then(() => {

			wof_format(str_data).then((fmt_rsp) => {

			    raw.innerText = fmt_rsp;

			    wof.edit.api.save(uri, fmt_rsp).then(rsp =>
				rsp.json()
			    ).then((data) => {

				console.log("OKAY SAVED");

				// REDRAW MAP HERE...
				    
				const str_data = JSON.stringify(data);

				wof_format(str_data).then((fmt_rsp) => {
				    raw.innerText = fmt_rsp;
				}).catch((err) => {
				    console.error("Document saved, failed to format response", err);
				});
				
			    }).catch((err) => {
				console.error("Failed to save data", err);
			    });
			    
			}).catch((err) => {
			    console.error("Failed to format data", err);
			    save_btn.setAttribute("disabled", "disabled");						    
			});
			
		    }).catch((err) => {
			console.error("Failed to validate data", err);
			save_btn.setAttribute("disabled", "disabled");						
		    });
		};
		
	    }).catch((err) => {
		console.error("SAD", err)
	    });
	},
    };
    
    return self;
})();
