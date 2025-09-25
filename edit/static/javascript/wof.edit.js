var wof = wof || {};

wof.edit = (function () {

    var feature_layer = null;
    var alert_timeout = null;
    
    var self = {
	
	init: function() {
	    
	    return new Promise((resolve, reject) => {

		const spinner = document.querySelector("#spinner-svg");
		spinner.style.display = "inline-block";
		
		sfomuseum.golang.wasm.fetch("wasm/wof_edit.wasm").then((rsp) => {
		    spinner.style.display = "none";		    
		    resolve();
		}).catch((err) => {
		    spinner.style.display = "none";		    		    
		    reject("Failed to load wof_placetypes WASM binary " + err);
		});
		
	    });
	},

	feedback: function(msg, ttl){

	    if (! ttl){
		ttl = 5000;
	    }
	    
	    self.alert(msg, ttl);
	},

	alert: function(msg, ttl){

	    const dlg = document.createElement("dialog");
	    dlg.setAttribute("id", "alert");
	    dlg.setAttribute("class", "alert");
	    
	    const close = document.createElement("div");
	    close.setAttribute("class", "alert-close");

		close.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-x-circle" viewBox="0 0 16 16"><path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14m0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16"/><path d="M4.646 4.646a.5.5 0 0 1 .708 0L8 7.293l2.646-2.647a.5.5 0 0 1 .708.708L8.707 8l2.647 2.646a.5.5 0 0 1-.708.708L8 8.707l-2.646 2.647a.5.5 0 0 1-.708-.708L7.293 8 4.646 5.354a.5.5 0 0 1 0-.708"/></svg>';	    

	    close.onclick = function(e){

		if (alert_timeout){
		    clearTimeout(alert_timeout);
		}
		
		dlg.close()
		document.body.removeChild(dlg);		
	    };

	    const body = document.createElement("div");
	    body.setAttribute("class", "alert-body");
	    
	    body.appendChild(document.createTextNode(msg));

	    dlg.appendChild(close);
	    dlg.appendChild(body);

	    document.body.prepend(dlg);	    
	    dlg.showModal();

	    if (ttl){

		if (alert_timeout){
		    clearTimeout(alert_timeout);
		}
		
		alert_timeout = setTimeout(function(){

		    const dlg = document.querySelector("#alert");
		    console.debug("Remove dialog", dlg);
		    
		    if (dlg){
			dlg.close()
			document.body.removeChild(dlg);
		    }
		    
		}, ttl);
	    }
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
		_self.alert("Failed to retrieve list of records to edit, " + err);
		console.error("Failed to retrieve list to edit", err);
	    });
	    
	},

	show: function(uri) {

	    const _self = self;
	    
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
		right.setAttribute("id", "right");
		right.appendChild(raw);

		// Populate form
		
		const form_t = document.querySelector("#edit-form");
		const form = _self.populate_form(form_t, data);

		const form_wrapper = document.createElement("div");
		form_wrapper.setAttribute("id", "form-wrapper");
		form_wrapper.appendChild(form);
		
		right.appendChild(form_wrapper);

		// Set up feature wrapper
		
		var wrapper = document.createElement("div");
		wrapper.setAttribute("id", "feature");

		wrapper.appendChild(left);
		wrapper.appendChild(right);

		// Set up buttons and feedback

		const form_btn = document.createElement("button");
		form_btn.setAttribute("class", "btn btn-primary");
		form_btn.appendChild(document.createTextNode("Form view"));
		form_btn.setAttribute("disabled", "disabled");
		
		const data_btn = document.createElement("button");
		data_btn.setAttribute("class", "btn btn-primary");
		data_btn.appendChild(document.createTextNode("Data view"));

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

		buttons.appendChild(form_btn);
		buttons.appendChild(data_btn);				
		buttons.appendChild(format_btn);		
		buttons.appendChild(validate_btn);
		buttons.appendChild(save_btn);

		// Set up UI

		format_btn.setAttribute("disabled", "disabled");
		validate_btn.setAttribute("disabled", "disabled");
		
		const navbar = document.querySelector("#navbar-content");
		navbar.appendChild(buttons);
		
		const feedback = document.createElement("div");
		feedback.setAttribute("id", "feedback");
		
		const ui = document.createElement("div");
		ui.appendChild(feedback);				
		ui.appendChild(wrapper);
		
		const root = document.getElementById("canvas");
		root.innerHTML = "";
		root.appendChild(ui);

		// Format data and add to pre element
		
		wof_format(str_data).then((fmt_rsp) => {
		    raw.innerText = fmt_rsp;
		}).catch((err) => {
		    _self.alert("Failed to format data, " + err);
		    console.error("Failed to format data", err);
		})

		// Set up map
		
		const bounds = whosonfirst.geojson.deriveBboxAsBounds(data);
		
		const map = L.map(map_id);
		map.fitBounds(bounds);
		
		const osm = L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {});
		osm.addTo(map);

		feature_layer = L.geoJSON(data);
		feature_layer.addTo(map);

		// Set up button interactions

		form_btn.onclick = function(){
		    
		    form_btn.setAttribute("disabled", "disabled");		    
		    data_btn.removeAttribute("disabled");

		    format_btn.setAttribute("disabled", "disabled");
		    validate_btn.setAttribute("disabled", "disabled");
		    
		    // START OF rebuild form

		    var data;
		    
		    try {
			const row = document.querySelector("#raw");
			data = JSON.parse(raw.innerText);
		    } catch(err) {
			_self.alert("Failed to parse raw data, " + err);
			console.error("Failed to parse raw data", err);
			return false;
		    }
			
		    const form_t = document.querySelector("#edit-form");
		    const form = _self.populate_form(form_t, data);
			
		    form_wrapper.innerHTML = "";
		    form_wrapper.appendChild(form);
		    
		    form_wrapper.style.display = "block";

		    // END OF rebuild form
		    
		    raw.style.display = "none";
		    return false;
		};

		data_btn.onclick = function(){
		    data_btn.setAttribute("disabled", "disabled");		    
		    form_btn.removeAttribute("disabled");

		    format_btn.removeAttribute("disabled");
		    validate_btn.removeAttribute("disabled");
		    
		    raw.style.display = "block";
		    form_wrapper.style.display = "none";
		    return false;
		};
		
		format_btn.onclick = function(){

		    var data;

		    try {
			data = JSON.parse(raw.innerText);
		    } catch(err) {
			_self.alert("Failed to parse raw data, " + err);
			console.error("Failed to parse data", err);
			save_btn.setAttribute("disabled", "disabled");
			return;
		    }
		    
		    const str_data = JSON.stringify(data);
		    
		    wof_format(str_data).then((fmt_rsp) => {
			raw.innerText = fmt_rsp;
			save_btn.removeAttribute("disabled");
		    }).catch((err) => {
			_self.alert("Failed to format raw data, " + err);
			console.error("Failed to format data", err)
			save_btn.setAttribute("disabled", "disabled");
		    });
		};

		validate_btn.onclick = function() {

		    var data;

		    try {
			data = JSON.parse(raw.innerText);
		    } catch(err) {
			_self.alert("Failed to parse raw data, " + err);
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
			    _self.alert("Failed to format raw data, " + err);			    
			    console.error("Failed to format data", err);
			    save_btn.setAttribute("disabled", "disabled");						    
			});
			
		    }).catch((err) => {
			_self.alert("Data validation failed, " + err);
			console.error("Failed to validate data", err);
			save_btn.setAttribute("disabled", "disabled");						
		    });		    
		};

		save_btn.onclick = function(){

		    const spinner = document.querySelector("#spinner-svg");
		    spinner.style.display = "inline-block";
		    
		    const raw = document.querySelector("#raw");		    
		    var data;

		    try {
			data = JSON.parse(raw.innerText);
		    } catch(err) {
			spinner.style.display = "none";		    			
			save_btn.setAttribute("disabled", "disabled");
			_self.alert("Failed to parse raw data, " + err);
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

				const bounds = whosonfirst.geojson.deriveBboxAsBounds(data);
				map.fitBounds(bounds);
				map.removeLayer(feature_layer);
				feature_layer = L.geoJSON(data);
				feature_layer.addTo(map);
				
				const str_data = JSON.stringify(data);

				wof_format(str_data).then((fmt_rsp) => {
				    spinner.style.display = "none";		    				    
				    _self.feedback("Data saved", 2000);				    
				    raw.innerText = fmt_rsp;
				}).catch((err) => {
				    spinner.style.display = "none";		    				    
				    _self.alert("Data save but document formatting failed, " + err);
				    console.error("Document saved, failed to format response", err);
				});
				
			    }).catch((err) => {
				spinner.style.display = "none";		    
				_self.alert("Failed to save data, " + err);
				console.error("Failed to save data", err);
			    });
			    
			}).catch((err) => {
			    spinner.style.display = "none";		    			    
			    _self.alert("Data formatting failed, " + err);
			    console.error("Failed to format data", err);
			    save_btn.setAttribute("disabled", "disabled");						    
			});
			
		    }).catch((err) => {
			spinner.style.display = "none";		    			
			_self.alert("Data validation failed, " + err);			
			console.error("Failed to validate data", err);
			save_btn.setAttribute("disabled", "disabled");						
		    });
		};
		
	    }).catch((err) => {
		_self.alert("Failed to retrieve record, " + err);
		console.error("Failed to retrieve record", err)
	    });
	},

	populate_form: function(t, data){

	    const _self = self;
	    
	    const form = t.content.cloneNode(true);
	    
	    const inputs = form.querySelectorAll(".wof-input");
	    const count_inputs = inputs.length;
	    
	    for (var i=0; i < count_inputs; i++){
		
		const input_el = inputs[i];

		input_el.onchange = function(e){
		    
		    const el = e.target;
		    const id = el.getAttribute("id");

		    var data;
		    
		    try {
			const row = document.querySelector("#raw");
			data = JSON.parse(raw.innerText);
		    } catch(err) {
			_self.alert("Failed to parse raw data, " + err);
			console.error("Failed to parse raw data", err);
			return false;
		    }
		    
		    switch (id){
			default:

			    if (el.value == ""){
				delete data.properties[id];
			    } else {
				data.properties[id] = el.value;
			    }
			    
			    break;
		    }

		    const str_data = JSON.stringify(data);

		    wof_validate(str_data).then(() => {
			wof_format(str_data).then((fmt_rsp) => {
			    raw.innerText = fmt_rsp;
			}).catch((err) => {
			    console.error("Failed to format response", err);
			});
		    }).catch((err) => {
			_self.alert("Data validation failed, " + err);			
			console.error("Data validation failed", err);
		    });
		};
		
		const id = input_el.getAttribute("id");
		
		if (! id in data.properties){
		    console.debug("Missing property", id);
		    continue;
		}

		switch (id){
		    case "wof:placetype":
			
			// wof_edit.wasm
			wof_placetypes().then(rsp => {
			    
			    const placetypes = JSON.parse(rsp);
			    const count = placetypes.length;
			    
			    for (var j=0; j < count; j++){
				
				const pt = placetypes[j].name;
				
				const opt = document.createElement("option");
				opt.setAttribute("value", pt);
				
				if (pt == data.properties[id]){
				    opt.setAttribute("selected", "selected");
				}
				
				opt.appendChild(document.createTextNode(pt));
				input_el.appendChild(opt);					
			    }
			    
			}).catch((err) => {
			    _self.alert("There was a problem listing placetypes, " + err);
			    console.error("Failed to derive placetypes", err)
			});
			
			break;

		    case "mz:is_current":
			self.populate_existential_flag(input_el, data.properties[id]);
			break;
		    case "mz:is_funky":
			self.populate_existential_flag(input_el, data.properties[id]);			
			break;
		    default:

			var v = "";

			if (data.properties[id]){
			    v = data.properties[id];
			};
			
			input_el.value = v;
			break;
		}
	    }

	    // wof:concordances
	    
	    const concordances_el = form.querySelector("#wof-concordances");
	    const concordances_t = document.querySelector("#wof-concordances-row");	    

	    const concordances_update_func = function(e){
		const el = e.target;
		
		const prefix = el.getAttribute("name");
		const v = el.value;

		if (v == ""){

		    if (! confirm("Do you want to remove the concordance for " + prefix + "?")){

			// BUT THEN WHAT? replace old value or... ?
			return false;
		    }

		    _self.remove_concordance(prefix).catch((err) => {
			_self.alert("Failed to remove the concordance for " + prefix + ", " + err);
		    });
		    
		    return false;
		}
		
		_self.update_concordance(prefix, v).then(() => {
		    el.value = v;
		}).catch((err) => {
		    _self.alert("Failed to update concordance for " + prefix + ", " + err);
		});
		
		return false;
	    };
	    
	    const concordances_remove_func = function(e){
		
		e.stopPropagation();
		
		const el = e.target;
		const prefix = el.getAttribute("data-prefix");
		
		if (! prefix){
		    _self.alert("Failed to remove concordance, unable to determine prefix");
		    return false;
		}
		
		if (! confirm("Are you sure you want to delete the concordance for " + prefix + "?")){
		    return false;
		}

		_self.remove_concordance(prefix).then(() => {
		    _self.feedback("The concordance for " + prefix + " has been removed", 1500);
		}).catch((err) => {
		    _self.alert("Failed to remove the concordance for " + prefix + ", " + err);
		});
		
		return false;
	    };
	    
	    console.log("CONCORDANCES", concordances_el, concordances_t);

	    // Add concordance button event(s) here...
	    
	    for (const k in data.properties["wof:concordances"]){

		const v = data.properties["wof:concordances"][k];
		const row = concordances_t.content.cloneNode(true);

		const wrapper = row.querySelector(".concordance-row");
		wrapper.setAttribute("id", "wof-concordances-" + k);
		
		// Note: It is important to update row _before_ appending it to concordances_el
		
		const label = row.querySelector(".label");
		label.innerText = k;	// Lookup name for k here (whosonfirst-sources)
		
		const prefix = row.querySelector(".prefix");
		prefix.innerText = k;

		const input = row.querySelector(".form-control");
		input.setAttribute("name", k);
		input.setAttribute("value", v);

		input.onchange = concordances_update_func;

		// START OF this is annoying...
		    
		const remove_btn = row.querySelector(".concordance-rm");
		remove_btn.setAttribute("data-prefix", k);
		remove_btn.onclick = concordances_remove_func;

		const remove_svg = remove_btn.querySelector("svg");
		remove_svg.setAttribute("data-prefix", k);
		remove_svg.onclick = concordances_remove_func;

		const remove_paths = remove_svg.querySelectorAll("path");
		const count_paths = remove_paths.length;
		
		for (var p=0; p < count_paths; p++){
		    const path = remove_paths[p];
		    path.setAttribute("data-prefix", k);
		    path.onclick = concordances_remove_func;
		}
		
		// END OF this is annoying...
		    
		concordances_el.appendChild(row);		
	    }
	    
	    return form;
	},

	populate_existential_flag: function(input_el, flag_v){

	    // TBD: Just do this from a template?
	    
	    const flags = {
		"true": 1,
		"false": 0,
		"unknown": -1,			    
	    };
	    
	    for (const k in flags) {
		const v = flags[k];
		
		const opt = document.createElement("option");
		opt.setAttribute("value", v);
		
		if (v == flag_v){
		    opt.setAttribute("selected", "selected");
		}
		
		opt.appendChild(document.createTextNode(k));
		input_el.appendChild(opt);					
	    }
	    
	},

	update_concordance: function(prefix, v){

	    const _self = self;

	    return new Promise((resolve, reject) => {
		
		_self.start_spinner();
		
		try {
		    const row = document.querySelector("#raw");
		    data = JSON.parse(raw.innerText);
		} catch(err) {
		    console.error("Failed to parse raw data", err);
		    _self.stop_spinner();
		    
		    reject("Failed to parse raw data, " + err);
		    return;
		}
	    	    
		data.properties["wof:concordances"][prefix] = v;
		console.debug("Update concordances", data.properties["wof:concordances"]);
		
		const str_data = JSON.stringify(data);
		
		wof_format(str_data).then((fmt_rsp) => {
		    raw.innerText = fmt_rsp;
		    _self.stop_spinner();
		    
		    resolve();
		    return;
		}).catch((err) => {
		    console.error("Failed to format data", err);
		    _self.stop_spinner();
		    
		    reject("Failed to format raw data, " + err);
		    return;
		});
	    });
	},
	
	remove_concordance: function(prefix){

	    const _self = self;

	    return new Promise((resolve, reject) => {
		
		_self.start_spinner();
		
		try {
		    const row = document.querySelector("#raw");
		    data = JSON.parse(raw.innerText);
		} catch(err) {
		    console.error("Failed to parse raw data", err);
		    _self.stop_spinner();

		    reject("Failed to parse raw data, " + err);		    
		    return;
		}
	    
		try {
		    const row = document.getElementById("wof-concordances-" + prefix);
		    const parent = row.parentNode;
		    parent.removeChild(row);
		} catch(err) {
		    console.error("Failed to remove row", err);
		    _self.stop_spinner();

		    reject("Failed to remove row, " + err);		    
		    return false;
		}
	    
		delete(data.properties["wof:concordances"][prefix]);
		const str_data = JSON.stringify(data);
		
		wof_format(str_data).then((fmt_rsp) => {
		    raw.innerText = fmt_rsp;
		    _self.stop_spinner();

		    resolve();
		    return;
		}).catch((err) => {
		    console.error("Failed to format data", err);
		    _self.stop_spinner();

		    reject("Failed to format raw data, " + err);		    
		    return;
		});

	    });
	},
	
	start_spinner: function(){
		const spinner = document.querySelector("#spinner-svg");
		spinner.style.display = "inline-block";
	},

	stop_spinner: function(){
	    const spinner = document.querySelector("#spinner-svg");
	    spinner.style.display = "none";
	},
	

    };
    
    return self;
    
})();
