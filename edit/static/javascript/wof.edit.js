var wof = wof || {};

wof.edit = (function () {

    var feature_layer = null;
    var alert_timeout = null;

    const lang_spoken = ["ara", "ara_AE", "ben", "ben_IN", "ben_BD", "dan", "ell", "eng", "eng_GB", "eng_US", "fin", "fra", "ger", "ind", "ita", "jpn", "kan", "kor", "mal", "nld", "nor", "pol", "por", "por_BR", "por_PT", "ron", "rus", "spa", "spa_AR", "spa_MX", "spa_ES", "swe", "tam", "tel", "tha", "tur", "zho", "zho_CN", "zho_TW"];

    const lang_official = ["ara", "ara_AE", "ben", "ben_IN", "ben_BD", "dan", "ell", "eng", "eng_GB", "eng_US", "fin", "fra", "ger", "ind", "ita", "jpn", "kan", "kor", "mal", "nld", "nor", "pol", "por", "por_BR", "por_PT", "ron", "rus", "spa", "spa_AR", "spa_MX", "spa_ES", "swe", "tam", "tel", "tha", "tur", "zho", "zho_CN", "zho_TW"];

    const lang_ok = [{"lang": "afr", "searchBy": ["afr"], "value": "afr"}, {"lang": "ara", "searchBy": ["ara", "Arabic"], "value": "Arabic (ara)"}, {"lang": "arz", "searchBy": ["arz"], "value": "arz"}, {"lang": "ben", "searchBy": ["ben", "Bengali"], "value": "Bengali (ben)"}, {"lang": "bul", "searchBy": ["bul"], "value": "bul"}, {"lang": "cat", "searchBy": ["cat"], "value": "cat"}, {"lang": "ces", "searchBy": ["ces"], "value": "ces"}, {"lang": "cym", "searchBy": ["cym"], "value": "cym"}, {"lang": "dan", "searchBy": ["dan", "Danish"], "value": "Danish (dan)"}, {"lang": "deu", "searchBy": ["deu"], "value": "deu"}, {"lang": "eng", "searchBy": ["eng", "English"], "value": "English (eng)"}, {"lang": "epo", "searchBy": ["epo"], "value": "epo"}, {"lang": "est", "searchBy": ["est"], "value": "est"}, {"lang": "eus", "searchBy": ["eus"], "value": "eus"}, {"lang": "fas", "searchBy": ["fas"], "value": "fas"}, {"lang": "fin", "searchBy": ["fin", "Finnish"], "value": "Finnish (fin)"}, {"lang": "fra", "searchBy": ["fra", "French"], "value": "French (fra)"}, {"lang": "guj", "searchBy": ["guj"], "value": "guj"}, {"lang": "heb", "searchBy": ["heb"], "value": "heb"}, {"lang": "hun", "searchBy": ["hun"], "value": "hun"}, {"lang": "hye", "searchBy": ["hye"], "value": "hye"}, {"lang": "ind", "searchBy": ["ind", "Indonesian"], "value": "Indonesian (ind)"}, {"lang": "ita", "searchBy": ["ita", "Italian"], "value": "Italian (ita)"}, {"lang": "jpn", "searchBy": ["jpn", "Japanese"], "value": "Japanese (jpn)"}, {"lang": "kat", "searchBy": ["kat"], "value": "kat"}, {"lang": "kor", "searchBy": ["kor", "Korean"], "value": "Korean (kor)"}, {"lang": "ltz", "searchBy": ["ltz"], "value": "ltz"}, {"lang": "mar", "searchBy": ["mar"], "value": "mar"}, {"lang": "msa", "searchBy": ["msa"], "value": "msa"}, {"lang": "nav", "searchBy": ["nav"], "value": "nav"}, {"lang": "nld", "searchBy": ["nld", "Dutch"], "value": "Dutch (nld)"}, {"lang": "nno", "searchBy": ["nno"], "value": "nno"}, {"lang": "nob", "searchBy": ["nob"], "value": "nob"}, {"lang": "nor", "searchBy": ["nor", "Norwegian"], "value": "Norwegian (nor)"}, {"lang": "pdc", "searchBy": ["pdc"], "value": "pdc"}, {"lang": "pol", "searchBy": ["pol", "Polish"], "value": "Polish (pol)"}, {"lang": "por", "searchBy": ["por", "Portuguese"], "value": "Portuguese (por)"}, {"lang": "ron", "searchBy": ["ron", "Romanian"], "value": "Romanian (ron)"}, {"lang": "rus", "searchBy": ["rus", "Russian"], "value": "Russian (rus)"}, {"lang": "sco", "searchBy": ["sco"], "value": "sco"}, {"lang": "slk", "searchBy": ["slk"], "value": "slk"}, {"lang": "spa", "searchBy": ["spa", "Spanish"], "value": "Spanish (spa)"}, {"lang": "srp", "searchBy": ["srp"], "value": "srp"}, {"lang": "swe", "searchBy": ["swe", "Swedish"], "value": "Swedish (swe)"}, {"lang": "tam", "searchBy": ["tam", "Tamil"], "value": "Tamil (tam)"}, {"lang": "tgl", "searchBy": ["tgl"], "value": "tgl"}, {"lang": "tha", "searchBy": ["tha", "Thai"], "value": "Thai (tha)"}, {"lang": "tur", "searchBy": ["tur", "Turkish"], "value": "Turkish (tur)"}, {"lang": "ukr", "searchBy": ["ukr"], "value": "ukr"}, {"lang": "urd", "searchBy": ["urd"], "value": "urd"}, {"lang": "vie", "searchBy": ["vie"], "value": "vie"}, {"lang": "yue", "searchBy": ["yue"], "value": "yue"}, {"lang": "zho", "searchBy": ["zho", "Chinese"], "value": "Chinese (zho)"}];
    
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
		
		const form = _self.populate_form(data);

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
			
		    const form = _self.populate_form(data);
			
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

	populate_form: function(data){

	    const form_t = document.querySelector("#edit-form");	    
	    const form = form_t.content.cloneNode(true);

	    self.populate_form_wof_input(form, data);
	    self.populate_form_names(form, data);	    	    
	    self.populate_form_concordances(form, data);
	    
	    return form;
	},

	populate_form_wof_input: function(form, data){

	    const _self = self;
	    
	    const inputs = form.querySelectorAll(".wof-input");
	    const count_inputs = inputs.length;
	    
	    for (var i=0; i < count_inputs; i++){
		
		const input_el = inputs[i];
		const input_id = input_el.getAttribute("id");
		
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
			case "wof:lang_x_spoken":

			    try {
				
				var tag_langs = [];
				
				const str_langs = el.value;
				console.log("VALUE", str_langs);
				
				if (str_langs != ""){
				    tag_langs = JSON.parse(str_langs);
				}
				
			    } catch(err) {
				console.error("Failed to parse tag data", err);
				_self.alert("Failed to parse tag data", err);
				return false;
			    }
			    
			    const count_langs = tag_langs.length;
			    console.log("COUNT", count_langs);
			    
			    if (count_langs == 0){
				
				if ("wof:lang_x_spoken" in data.properties){
				    delete data.properties["wof:lang_x_spoken"];			    
				} else {
				    return;
				}
				
			    } else {
				
				var langs = [];
				
				for (var i=0; i < count_tags; i++){
				    console.log("LANG", i, tag_langs.value);
				    langs.push(tag_langs[i].value);
				}
				
				console.log("SET", langs);
				data.properties["wof:lang_x_spoken"] = langs;
			    }

			    break;
			    
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
			    _self.alert("Failed to format response " + err);
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
		    case "wof:lang_x_spoken":

			console.log("FOO");
			// FORMAT wof:lang_x_spoken for tagify here..

			    try {
				Tagify(input_el, {
				    whiteList: lang_spoken
				});			
			    } catch(err) {
				console.error("SAD", err, input_el);
			    };
			    
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
	    
	},
	
	populate_form_names: function(form, data){

	    const _self = self;

	    const names_group = form.querySelector("#localized-names-group")
	    const names_desc = form.querySelector("#localized-names-description");
	    const names_count = form.querySelector("#localized-names-count");
	    
	    const t = document.querySelector("#localized-names-row");
	    const row = t.content.cloneNode(true);

	    // names_group.appendChild(row);
	    
	    var count_names = 0;
	    
	    for (var k in data.properties){

		if (k.startsWith("name:")){
		    count_names += 1;
		}
	    }

	    switch (count_names){
		case 0:
		    break;
		case 1:
		    names_count.innerText = "1 language";
		    names_desc.style.display = "inline-block";
		    break;
		default:
		    names_count.innerText = count_names + " languages";
		    names_desc.style.display = "inline-block";
		    break;
	    }		    
	},
	
	populate_form_concordances: function(form, data){

	    const _self = self;
	    
	    const concordances_el = form.querySelector("#wof-concordances");
	    const concordances_t = document.querySelector("#wof-concordances-row");	    

	    // Set up wof:concordances button event functions
	    // The order of these functions is relevant, meaning that the "add" function
	    // references both the "update" and "remove" functions so it needs to be defined
	    // last. But wait, there's more! Nestled in between the "update" and "remove" and the
	    // "add" function is a "new_concordances_row" function. Again, this is all necessary
	    // because of scoping wah-wah. Could all of these be moved in to dicrete functions?
	    // Probably, but that is tomorrow's problem right now.
	    
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

	    const new_concordances_row = function(prefix, value){

		const row = concordances_t.content.cloneNode(true);
		
		const wrapper = row.querySelector(".concordance-row");
		wrapper.setAttribute("id", "wof-concordances-" + prefix);
		
		// Note: It is important to update row _before_ appending it to concordances_el
		
		const label_el = row.querySelector(".label");
		label_el.innerText = prefix;	// Lookup name for k here (whosonfirst-sources)
		
		const prefix_el = row.querySelector(".prefix");
		prefix_el.innerText = prefix;
		
		const input_el = row.querySelector(".form-control");
		input_el.setAttribute("name", prefix);
		input_el.setAttribute("value", value);
		
		input_el.onchange = concordances_update_func;

		// START OF this is annoying...
		
		const remove_btn = row.querySelector(".concordance-rm");
		remove_btn.setAttribute("data-prefix", prefix);
		remove_btn.onclick = concordances_remove_func;

		const remove_svg = remove_btn.querySelector("svg");
		remove_svg.setAttribute("data-prefix", prefix);
		remove_svg.onclick = concordances_remove_func;

		const remove_paths = remove_svg.querySelectorAll("path");
		const count_remove = remove_paths.length;
		
		for (var p=0; p < count_remove; p++){
		    const path = remove_paths[p];
		    path.setAttribute("data-prefix", prefix);
		    path.onclick = concordances_remove_func;
		}

		return row;
	    };
	    
	    const concordances_add_func = function(e){
		
		e.stopPropagation()
		
		const dlg = document.createElement("dialog");
		dlg.setAttribute("id", "dialog");
		dlg.setAttribute("class", "dialog");
		
		const exit_func = function(){
		    dlg.close()
		    document.body.removeChild(dlg);		
		};

		const close = document.createElement("div");
		close.setAttribute("class", "dialog-close");
		
		close.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-x-circle" viewBox="0 0 16 16"><path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14m0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16"/><path d="M4.646 4.646a.5.5 0 0 1 .708 0L8 7.293l2.646-2.647a.5.5 0 0 1 .708.708L8.707 8l2.647 2.646a.5.5 0 0 1-.708.708L8 8.707l-2.646 2.647a.5.5 0 0 1-.708-.708L7.293 8 4.646 5.354a.5.5 0 0 1 0-.708"/></svg>';	    
		
		close.onclick = function(e){
		    exit_func();
		    return false;
		};
		
		const body = document.createElement("div");
		body.setAttribute("class", "dialog-body");

		const new_t = document.querySelector("#new-concordances-row");
		const row = new_t.content.cloneNode(true);
		
		body.appendChild(row);

		const buttons = document.createElement("div");
		buttons.setAttribute("id", "new-concordance-buttons");

		const add_btn_func = function(e){

		    _self.start_spinner();
		    
		    const prefix_el = document.querySelector("#new-wof-concordances-name");
		    const value_el = document.querySelector("#new-wof-concordances-value");

		    if (! prefix_el){
			_self.stop_spinner();			
			_self.alert("Failed to derive prefix");
			return false;
		    }

		    if (! value_el){
			_self.stop_spinner();			
			_self.alert("Failed to derive valud");
			return false;
		    }

		    const prefix = prefix_el.value;
		    const value = value_el.value;

		    if (! prefix){
			_self.stop_spinner();			
			_self.alert("Prefix must not be blank");
			return false;
		    }

		    if (! value){
			_self.stop_spinner();			
			_self.alert("Value must not be blank");
			return false;
		    }

		    var data;

		    try {
			const row = document.querySelector("#raw");
			data = JSON.parse(raw.innerText);
		    } catch(err) {
			console.error("Failed to parse raw data", err);
			_self.stop_spinner();
			_self.alert("Failed to parse data, " + err);
			return false;
		    };

		    if (prefix in data.properties["wof:concordances"]){
			_self.stop_spinner();
			_self.alert("There is already a concordance with that prefix");
			return false;
		    }

		    data.properties["wof:concordances"][prefix] = value;
			
		    const str_data = JSON.stringify(data);
		    
		    wof_validate(str_data).then(() => {

			wof_format(str_data).then((fmt_rsp) => {
			    
			    raw.innerText = fmt_rsp;

			    const row = new_concordances_row(prefix, value);
			    concordances_el.prepend(row);		
			    
			    exit_func();
			    _self.stop_spinner();			    
			    _self.feedback("New concordance added", 1500);
			    
			}).catch((err) => {
			    _self.stop_spinner();			    
			    _self.alert("Failed to format raw data, " + err);			    
			    console.error("Failed to format data", err);
			});
			
		    }).catch((err) => {
			_self.stop_spinner();			
			_self.alert("Data validation failed, " + err);
			console.error("Failed to validate data", err);
		    });		    
		    
		    return false;
		};
		
		const add_btn = document.createElement("button");
		add_btn.setAttribute("class", "btn btn-primary");
		add_btn.appendChild(document.createTextNode("Add concordance"));

		add_btn.onclick = function(e){
		    try {
			add_btn_func(e);
		    } catch(err) {
			_self.alert("Failed to add concordance, " + err);
			return false;
		    }
		};
		
		const cancel_btn = document.createElement("button");
		cancel_btn.setAttribute("class", "btn btn");
		cancel_btn.appendChild(document.createTextNode("Cancel"));

		cancel_btn.onclick = function(e){
		    exit_func();
		    return false;
		};
		
		buttons.appendChild(cancel_btn);		
		buttons.appendChild(add_btn);
		
		dlg.appendChild(close);
		dlg.appendChild(body);
		dlg.appendChild(buttons);
		
		document.body.prepend(dlg);	    
		dlg.showModal();
		
		console.log("Add");
		return false;
	    };
	    
	    const add_btn = form.querySelector("#wof-concordances-add");
	    add_btn.onclick = concordances_add_func;
	    
	    const add_svg = add_btn.querySelector("svg");
	    add_svg.onclick = concordances_add_func;
	    
	    const add_paths = add_svg.querySelectorAll("path");
	    const count_add = add_paths.length;
	    
	    for (var p=0; p < count_add; p++){
		const path = add_paths[p];
		path.onclick = concordances_add_func;
	    }

	    // Finally add all the existing concordances
	    
	    for (const k in data.properties["wof:concordances"]){
		const v = data.properties["wof:concordances"][k];
		const row = new_concordances_row(k, v);
		concordances_el.appendChild(row);		
	    }	    

	    
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
	    	    
		delete(data.properties["wof:concordances"][prefix]);
		const str_data = JSON.stringify(data);
		
		wof_validate(str_data).then(() => {
		    wof_format(str_data).then((fmt_rsp) => {

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
		}).catch((err) => {
		    console.error("Data validation failed", err);
		    _self.stop_spinner();
		    reject("Data validation failed, " + err);
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
