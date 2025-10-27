/**
 * @namespace wof.edit
 * @description Methods for editing Who's On First records
 */

var wof = wof || {};

wof.edit = (function () {

    var map;
    var feature_layer;
    var alert_timeout;

    // Perist the list of currently displayed language-based names
    // when toggling between form and data view. This value is updated
    // in the event handlers for tags_selected below.
    var current_names_langs = [
	"eng",
    ];

    // Perist the list of currently displayed language-based labels
    // when toggling between form and data view. This value is updated
    // in the event handlers for tags_selected below.    
    var current_labels_langs = [
	"eng",
    ];
    
    const lang_spoken = ["ara", "ara_AE", "ben", "ben_IN", "ben_BD", "dan", "ell", "eng", "eng_GB", "eng_US", "fin", "fra", "ger", "ind", "ita", "jpn", "kan", "kor", "mal", "nld", "nor", "pol", "por", "por_BR", "por_PT", "ron", "rus", "spa", "spa_AR", "spa_MX", "spa_ES", "swe", "tam", "tel", "tha", "tur", "zho", "zho_CN", "zho_TW"];

    const lang_official = ["ara", "ara_AE", "ben", "ben_IN", "ben_BD", "dan", "ell", "eng", "eng_GB", "eng_US", "fin", "fra", "ger", "ind", "ita", "jpn", "kan", "kor", "mal", "nld", "nor", "pol", "por", "por_BR", "por_PT", "ron", "rus", "spa", "spa_AR", "spa_MX", "spa_ES", "swe", "tam", "tel", "tha", "tur", "zho", "zho_CN", "zho_TW"];

    const lang_ok = [{"lang": "afr", "searchBy": ["afr"], "value": "afr"}, {"lang": "ara", "searchBy": ["ara", "Arabic"], "value": "Arabic (ara)"}, {"lang": "arz", "searchBy": ["arz"], "value": "arz"}, {"lang": "ben", "searchBy": ["ben", "Bengali"], "value": "Bengali (ben)"}, {"lang": "bul", "searchBy": ["bul"], "value": "bul"}, {"lang": "cat", "searchBy": ["cat"], "value": "cat"}, {"lang": "ces", "searchBy": ["ces"], "value": "ces"}, {"lang": "cym", "searchBy": ["cym"], "value": "cym"}, {"lang": "dan", "searchBy": ["dan", "Danish"], "value": "Danish (dan)"}, {"lang": "deu", "searchBy": ["deu"], "value": "deu"}, {"lang": "eng", "searchBy": ["eng", "English"], "value": "English (eng)"}, {"lang": "epo", "searchBy": ["epo"], "value": "epo"}, {"lang": "est", "searchBy": ["est"], "value": "est"}, {"lang": "eus", "searchBy": ["eus"], "value": "eus"}, {"lang": "fas", "searchBy": ["fas"], "value": "fas"}, {"lang": "fin", "searchBy": ["fin", "Finnish"], "value": "Finnish (fin)"}, {"lang": "fra", "searchBy": ["fra", "French"], "value": "French (fra)"}, {"lang": "guj", "searchBy": ["guj"], "value": "guj"}, {"lang": "heb", "searchBy": ["heb"], "value": "heb"}, {"lang": "hun", "searchBy": ["hun"], "value": "hun"}, {"lang": "hye", "searchBy": ["hye"], "value": "hye"}, {"lang": "ind", "searchBy": ["ind", "Indonesian"], "value": "Indonesian (ind)"}, {"lang": "ita", "searchBy": ["ita", "Italian"], "value": "Italian (ita)"}, {"lang": "jpn", "searchBy": ["jpn", "Japanese"], "value": "Japanese (jpn)"}, {"lang": "kat", "searchBy": ["kat"], "value": "kat"}, {"lang": "kor", "searchBy": ["kor", "Korean"], "value": "Korean (kor)"}, {"lang": "ltz", "searchBy": ["ltz"], "value": "ltz"}, {"lang": "mar", "searchBy": ["mar"], "value": "mar"}, {"lang": "msa", "searchBy": ["msa"], "value": "msa"}, {"lang": "nav", "searchBy": ["nav"], "value": "nav"}, {"lang": "nld", "searchBy": ["nld", "Dutch"], "value": "Dutch (nld)"}, {"lang": "nno", "searchBy": ["nno"], "value": "nno"}, {"lang": "nob", "searchBy": ["nob"], "value": "nob"}, {"lang": "nor", "searchBy": ["nor", "Norwegian"], "value": "Norwegian (nor)"}, {"lang": "pdc", "searchBy": ["pdc"], "value": "pdc"}, {"lang": "pol", "searchBy": ["pol", "Polish"], "value": "Polish (pol)"}, {"lang": "por", "searchBy": ["por", "Portuguese"], "value": "Portuguese (por)"}, {"lang": "ron", "searchBy": ["ron", "Romanian"], "value": "Romanian (ron)"}, {"lang": "rus", "searchBy": ["rus", "Russian"], "value": "Russian (rus)"}, {"lang": "sco", "searchBy": ["sco"], "value": "sco"}, {"lang": "slk", "searchBy": ["slk"], "value": "slk"}, {"lang": "spa", "searchBy": ["spa", "Spanish"], "value": "Spanish (spa)"}, {"lang": "srp", "searchBy": ["srp"], "value": "srp"}, {"lang": "swe", "searchBy": ["swe", "Swedish"], "value": "Swedish (swe)"}, {"lang": "tam", "searchBy": ["tam", "Tamil"], "value": "Tamil (tam)"}, {"lang": "tgl", "searchBy": ["tgl"], "value": "tgl"}, {"lang": "tha", "searchBy": ["tha", "Thai"], "value": "Thai (tha)"}, {"lang": "tur", "searchBy": ["tur", "Turkish"], "value": "Turkish (tur)"}, {"lang": "ukr", "searchBy": ["ukr"], "value": "ukr"}, {"lang": "urd", "searchBy": ["urd"], "value": "urd"}, {"lang": "vie", "searchBy": ["vie"], "value": "vie"}, {"lang": "yue", "searchBy": ["yue"], "value": "yue"}, {"lang": "zho", "searchBy": ["zho", "Chinese"], "value": "Chinese (zho)"}];
    
    const close_svg = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-x-circle" viewBox="0 0 16 16"><path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14m0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16"/><path d="M4.646 4.646a.5.5 0 0 1 .708 0L8 7.293l2.646-2.647a.5.5 0 0 1 .708.708L8.707 8l2.647 2.646a.5.5 0 0 1-.708.708L8 8.707l-2.646 2.647a.5.5 0 0 1-.708-.708L7.293 8 4.646 5.354a.5.5 0 0 1 0-.708"/></svg>';
    
    var self = {

	/**
	* @function init
	* @memberof wof.edit
	* @description Initialize the WOF edit functionality. This fetches, loads and initilized the 'wof_edit' WASM binary.
        * @return {Promise} – 
        */		   	
	init: function() {

	    const _self = self;
	    
	    return new Promise((resolve, reject) => {

		_self.start_spinner();
		
		sfomuseum.golang.wasm.fetch("wasm/wof_edit.wasm").then((rsp) => {

		    _self.stop_spinner();		    

		    const home = document.querySelector("#home-button");
		    
		    home.onclick = function(){
			_self.list(true);
			return false;
		    };
		    
		    resolve();
		}).catch((err) => {
		    _self.stop_spinner();		    		    
		    reject("Failed to load wof_placetypes WASM binary " + err);
		});
		
	    });
	},

	/**
	 * @function list
	 * @memberof wof.edit
	 * @description Fetch the list of WOF records to edit and display them in a list. If the list only contains one record then the 'show' method is immediately dispatched for that record.
         * @return {null}
        */		   			
	list: function(force) {

	    const _self = self;
	    
	    wof.edit.api.list().then((rsp) => {
		
		const count = rsp.length;

		if ((count == 1) && (! force)){
		    _self.show(rsp[0]);
		    return;
		}

		const btns = document.querySelector("#buttons");

		if (btns){
		    btns.parentNode.removeChild(btns);
		}
		
		const items = document.createElement("ul");
		items.setAttribute("id", "list-records");
		
		for (var i=0; i < count; i++) {
		    
		    const uri = rsp[i];
		    
		    const link = document.createElement("a");
		    link.setAttribute("href", uri);
		    link.setAttribute("id", uri);
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

		for (var i=0; i < count; i++) {
		    
		    const uri = rsp[i];

		    wof.edit.api.fetch(uri).then((data) => {
			const props = data.properties;
			const el = document.getElementById(uri);
			el.innerHTML = "";
			el.appendChild(document.createTextNode(props["wof:name"]));
			el.appendChild(document.createTextNode(" ("));
			
			const code = document.createElement("code");
			code.appendChild(document.createTextNode(props["wof:id"]));
			el.appendChild(code);
			el.appendChild(document.createTextNode(")"));
			
		    }).catch((err) => {
			console.error("Failed to update label for record", uri, err);
		    });
		}
		
	    }).catch((err) => {
		_self.alert("Failed to retrieve list of records to edit, " + err);
		console.error("Failed to retrieve list to edit", err);
	    });
	    
	},

	/**
	 * @function show
	 * @memberof wof.edit
	 * @description Display the edit interface for a specific WOF record.
	 * @param {string} uri - The URI of the WOF record to edit.
         * @return {null}
        */		   				
	show: function(uri) {

	    const _self = self;
	    
	    wof.edit.api.fetch(uri).then((data) => {

		// START OF monkey-patching pre-2019 EDTF values
		
		const edtf_props = [
		    "edtf:inception",
		    "edtf:cessation",		    
		];

		const edtf_count = edtf_props.length;

		for (var i=0; i < edtf_count; i++){

		    const prop = edtf_props[i];

		    if (prop in data.properties){

			var v = data.properties[prop];
			
			switch (v){
			    case "open":
				v = "..";
				break;
			    case "uuuu":
				v = "";
				break;
			    default:
				break;
			}

			data.properties[prop] = v;
		    }
		}

		// END OF monkey-patching pre-2019 EDTF values		
		
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
			
		map = L.map(map_id);

		const bounds = whosonfirst.geojson.deriveBboxAsBounds(data);		
		map.fitBounds(bounds);

		// https://geoman.io/docs/leaflet
		
		map.pm.addControls({  
		    position: 'topleft',
		    drawCircle: false,
		    drawMarker: false,		// don't draw default image-based markers		    
		    drawCircleMarker: true,	// draw circle-based markers instead
		    drawPolyline: false,	// there have never been polylines in WOF (or at least I don't think so)
		    drawRectangle: false,	// disabling (in favour of polygons) for the sake of less UI/chrome
		    drawText: false,
		    rotateMode: false,
		    
		});

		map.on("pm:drawend", function(e){
		    _self.update_geometry(map);		   
		});
		
		map.on('pm:remove', function (e) {
		    _self.update_geometry(map);
		});
		
		map.on('pm:globaleditmodetoggled', (e) => {
		    _self.update_geometry(map);		   		    
		});

		//

		const osm = L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {});
		osm.addTo(map);
		
		self.draw_feature_geometry(data);
		
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

		    // Remember: All we're doing here is _formatting_ the data and
		    // not validating it. That might become confusing enough that it
		    // becomes necessary to remove this feature. TBD.
			
		    wof_format(str_data).then((fmt_rsp) => {
			raw.innerText = fmt_rsp;
			const new_data = JSON.parse(fmt_rsp);
			// _self.draw_feature_geometry(new_data);
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

		    _self.export_data(data).then((fmt_rsp) => {
			raw.innerText = fmt_rsp;
			const new_data = JSON.parse(fmt_rsp);
			_self.draw_feature_geometry(new_data);						
		    }).catch((err) => {
			_self.alert("Data validation failed, " + err);
			console.error("Failed to validate (export) data", err);
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
				    _self.draw_feature_geometry(data);				    
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

	/**
	 * @function draw_feature_geometry
	 * @memberof wof.edit
	 * @description Draw (or redraw) the GeoJSON Feature layer on the map
	 * @param {Object} data – The GeoJSON Feature object used to populate the feature layer with
         * @return {null}
        */		   					
	draw_feature_geometry: function(data){

	    if (feature_layer){
		console.debug("Remove feature layer");
		map.removeLayer(feature_layer);
	    }
	    
	    const pt_style = wof.edit.leaflet.style('geom_centroid');
	    
	    const pt_args = {
		style: pt_style,
	    };
	    
	    const pt_handler = wof.edit.leaflet.point(pt_args);
	    const layer_style = wof.edit.leaflet.style('consensus_polygon');
	    
	    const layer_args = {
		style: layer_style,
		pointToLayer: pt_handler,
	    };

	    console.debug("Add feature layer");	    
	    feature_layer = L.geoJSON(data, layer_args);
	    feature_layer.addTo(map);

	    const bounds = whosonfirst.geojson.deriveBboxAsBounds(data);
	    console.debug("Fit map to bounds", bounds, data);
	    map.fitBounds(bounds);	    
	},
	
	/**
	 * @function populate_form
	 * @memberof wof.edit
	 * @description Populate the HTML/DOM edit form elements
	 * @param {Object} data – The GeoJSON Feature object used to populate the form
         * @return {null}
        */		   				
	populate_form: function(data){

	    const form_t = document.querySelector("#edit-form");	    
	    const form = form_t.content.cloneNode(true);

	    self.populate_form_wof_input(form, data);
	    self.populate_form_names(form, data);
	    self.populate_form_labels(form, data);	    	    	    
	    self.populate_form_concordances(form, data);
	    
	    return form;
	},

	/**
	 * @function populate_form_wof_input
	 * @memberof wof.edit
	 * @description Populate the HTML/DOM edit form for input elements with the "wof_input" CSS class
	 * @param {Object} form – The HTML/DOM form element to populate.
	 * @param {Object} data – The GeoJSON Feature object used to populate the form
         * @return {null}
        */		   					
	populate_form_wof_input: function(form, data){

	    const _self = self;
	    
	    const inputs = form.querySelectorAll(".wof-input");
	    const count_inputs = inputs.length;
	    
	    const onchange = function(e){
		    
		const el = e.target;
		const id = el.getAttribute("id");
		const re = el.getAttribute("data-wof-validation-regex");

		if ((re) && (! el.value.match(re))){
		    console.error("Property fails validation (re) test", id);
		    _self.alert(id + " property fails validation test");
		    return false
		}
		
		_self.load_data().then((change_data) => {

		    switch (id){
			case "wof:lang_x_spoken":
			    
			    var tag_langs_spoken = [];
			    
			    try {			   
				
				const str_langs = el.value;
				
				if (str_langs != ""){
				    tag_langs_spoken = JSON.parse(str_langs);
				}
				
			    } catch(err) {
				console.error("Failed to parse tag data", err);
				_self.alert("Failed to parse tag data", err);
				return false;
			    }
			    
			    const count_langs_spoken = tag_langs_spoken.length;
			    
			    if (count_langs_spoken == 0){
				
				if ("wof:lang_x_spoken" in change_data.properties){
				    delete change_data.properties["wof:lang_x_spoken"];			    
				} else {
				    return;
				}
				
			    } else {
				
				var langs = [];
				
				for (var i=0; i < count_langs_spoken; i++){
				    langs.push(tag_langs_spoken[i].value);
				}
				
				change_data.properties["wof:lang_x_spoken"] = langs;
			    }
			    
			    break;
			    
			case "wof:lang_x_official":
			    
			    var tag_langs_official = [];
			    
			    try {
				const str_langs = el.value;
				
				if (str_langs != ""){
				    tag_langs_official = JSON.parse(str_langs);
				}
				
			    } catch(err) {
				console.error("Failed to parse tag data", err);
				_self.alert("Failed to parse tag data", err);
				return false;
			    }
			    
			    const count_langs_official = tag_langs_official.length;
			    
			    if (count_langs_official == 0){
				
				if ("wof:lang_x_official" in change_data.properties){
				    delete change_data.properties["wof:lang_x_official"];			    
				} else {
				    return;
				}
				
			    } else {
				
				var langs = [];
				
				for (var i=0; i < count_langs_official; i++){
				    langs.push(tag_langs_official[i].value);
				}
				
				change_data.properties["wof:lang_x_official"] = langs;
			    }
			    
			    break;			
			default:
			    
			    if (el.value == ""){
				delete change_data.properties[id];
			    } else {
				change_data.properties[id] = el.value;
			    }
			    
			    break;
		    }

		    _self.save_data(change_data).then(() => {
			console.debug("Data saved", id);
		    }).catch((err) => {
			_self.alert("Failed to format response " + err);
			console.error("Failed to format response", err);
		    });
		    
		}).catch((err) => {
		    _self.alert("Failed to load raw data, " + err);
		    console.error("Failed to load raw data", err);
		    return false;
		});
		
	    };
	    
	    for (var i=0; i < count_inputs; i++){
		
		const input_el = inputs[i];
		const input_id = input_el.getAttribute("id");		
		
		if (! input_id in data.properties){
		    console.debug("Missing property", input_id);
		    continue;
		}

		input_el.onchange = onchange;
		
		switch (input_id){

		    case "wof:placetype":
			
			// wof_edit.wasm
			wof_placetypes().then(rsp => {
			    
			    const placetypes = JSON.parse(rsp);
			    const count = placetypes.length;
			    
			    for (var j=0; j < count; j++){
				
				const pt = placetypes[j].name;
				
				const opt = document.createElement("option");
				opt.setAttribute("value", pt);
				
				if (pt == data.properties[input_id]){
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
			self.populate_existential_flag(input_el, data.properties[input_id]);
			break;
		    case "mz:is_funky":
			self.populate_existential_flag(input_el, data.properties[input_id]);			
			break;
		    case "wof:lang_x_spoken":

			if (input_id in data.properties){
			    
			    const count_langs = data.properties[input_id].length;
			    var tag_langs = [];
			    
			    for (var l=0; l < count_langs; l++){
				tag_langs.push({
				    value: data.properties[input_id][l]
				});
			    }
			    
			    const str_langs = JSON.stringify(tag_langs);
			    input_el.value = str_langs;
			}
			
			// Remember: The 'new' is important. That might seem obvious to you
			// but it took me forever to figure that out...
			// https://github.com/yairEO/tagify/blob/master/src/parts/helpers.js#L119
			
			new Tagify(input_el, {
			    whiteList: lang_spoken
			});			
			
			break;

		    case "wof:lang_x_official":

			if (input_id in data.properties){
			    
			    const count_langs = data.properties[input_id].length;
			    var tag_langs = [];
			    
			    for (var l=0; l < count_langs; l++){
				tag_langs.push({
				    value: data.properties[input_id][l]
				});
			    }
			    
			    const str_langs = JSON.stringify(tag_langs);
			    input_el.value = str_langs;
			}
			
			// Remember: The 'new' is important. That might seem obvious to you
			// but it took me forever to figure that out...
			// https://github.com/yairEO/tagify/blob/master/src/parts/helpers.js#L119
			
			new Tagify(input_el, {
			    whiteList: lang_official
			});			
			
			break;			
		    default:

			var v = "";

			if (data.properties[input_id]){
			    v = data.properties[input_id];
			};
			
			input_el.value = v;
			break;
		}
	    }
	    
	},

	/**
	 * @function populate_form_names
	 * @memberof wof.edit
	 * @description Populate the HTML/DOM edit form for localized name input elements.
	 * @param {Object} form – The HTML/DOM form element to populate.
	 * @param {Object} data – The GeoJSON Feature object used to populate the form
         * @return {null}
        */		   						
	populate_form_names: function(form, data){

	    const _self = self;

	    // Count languages and display
	    
	    const names_group = form.querySelector("#localized-names-group")
	    const names_desc = form.querySelector("#localized-names-description");
	    const names_count = form.querySelector("#localized-names-count");

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

	    // Adding and removing languages
	    
	    const t = document.querySelector("#localized-names-row");
	    
	    const names_selected = form.querySelector("#localized-names-selected");

	    const tags_selected = new Tagify(names_selected, {
		whiteList: lang_ok,
	    });

	    // Current name/language tags are added below after we register
	    // the event handlers
	    
	    tags_selected.on("add", function(e){
		
		const lang = e.detail.data.value;
		const id = "localized-names-row-" + lang;
		
		console.debug("Add name(s)", lang, id);

		// Note: We are counting on Tagify.js deduping and not dispatching
		// the same name tag twice to this event.

		    try {
			const node = t.content.cloneNode(true);
			const row = node.querySelector(".localized-names-row");
			
			row.setAttribute("id", id);
			row.setAttribute("data-lang", lang);
			
			const name_el = row.querySelector(".lang-name");
			name_el.innerText = lang;	// LOOKUP NAME HERE
			
			const code_el = row.querySelector(".lang-code");
			code_el.innerText = lang;
			
			const input_els = row.querySelectorAll("input");
			const count_els = input_els.length;

			if (! current_names_langs.includes(lang)){
			    current_names_langs.push(lang);
			}
			
			// START OF if there's a built-in method to do this
			// I can't find it...
			    
			const list_tags = function(tag){

			    const tag_els = tag.getTagElms();
			    const count_tags = tag_els.length;
			    
			    var tags = [];
			    
			    for (var t=0; t < count_tags; t++){
				tags.push(tag_els[t].getAttribute("value"));
			    }

			    return tags;
			};

			// END OF if there's a built-in method to do this
			
			for (var i=0; i < count_els; i++){

			    const el = input_els[i];
			    const spec = el.getAttribute("data-spec");

			    const el_id = "name:" + lang + "_x_" + spec;
			    el.setAttribute("name", el_id);
			    el.setAttribute("id", el_id);

			    const el_tag = new Tagify(el);

			    // Add name tags
			    
			    if (el_id in data.properties){
				console.debug("Add name tags", el_id, data.properties[el_id]);
				el_tag.addTags(data.properties[el_id]);
			    }

			    // When new name tags are added (or updated)

			    const on_edit = function(e){
				
				const tags = list_tags(el_tag);

				if (tags.length == 0){
				    return false;
				}
				
				console.debug("Edit name tag", el_id, tags);

				_self.start_spinner();
				_self.load_data().then((data) => {

				    data.properties[el_id] = tags;
				    
				    _self.save_data(data).then(() => {
					console.debug("Saved name tags on edit", el_id, tags);
					_self.stop_spinner();
				    }).catch((err) => {
					_self.stop_spinner();
					console.error("Failed to save name tags on edit", el_id, err);
					_self.alert("Failed to save data, " + err);
				    });
				    
				}).catch((err) => {
				    _self.stop_spinner();
				    console.error("Failed to load data for name tags on edit", el_id, err);				    
				    _self.alert("Failed to load data, " + err)
				});
			    };
			    
			    el_tag.on("add", on_edit);
			    el_tag.on("change", on_edit);

			    // When existing name tags are removed
			    
			    el_tag.on("remove", function(e){

				if (! e.detail.data){
				    return;
				}

				const tags = list_tags(el_tag);
				console.debug("Remove name tags", el_id, tags);

				_self.start_spinner();
				_self.load_data().then((data) => {

				    if (tags.length == 0){

					if (el_id in data.properties){
					    delete(data.properties[el_id]);
					}
					
				    } else {
					data.properties[el_id] = tags;
				    }
				    
				    _self.save_data(data).then(() => {
					console.debug("Updated name tags on remove", el_id, tags);
					_self.stop_spinner();
				    }).catch((err) => {
					_self.stop_spinner();
					console.error("Failed to update name tags on remove", el_id, err);
					_self.alert("Failed to update name tags on remove, " + err);
				    });
				    
				}).catch((err) => {
				    _self.stop_spinner();
					console.error("Failed to load name tags on remove", el_id, err);				    
				    _self.alert("Failed to load data, " + err)
				});
				
			    });
			    
			}
			
			names_group.prepend(row);
			
		    } catch (err) {
			console.error("Failed to add name(s) row for ", lang, id, err);
		    }
	    });

	    tags_selected.on("remove", function(e){

		if (!e.detail.data){
		    return;
		}

		const lang = e.detail.data.value;
		const id = "localized-names-row-" + lang;
		
		console.debug("Remove name(s)", lang, id);

		const row = document.getElementById(id);

		if (! row){
		    console.warn("Failed to locate name(s) row for language", lang, id);
		    return;
		}

		row.parentNode.removeChild(row);

		if (current_names_langs.includes(lang)){

		    current_names_langs = current_names_langs.filter(function(item) {
			return item !== lang;
		    });
		}
		
	    });

	    console.debug("Add current name/language tags")
	    tags_selected.addTags(current_names_langs);
	},

	/**
	* @function populate_form_labels
	 * @memberof wof.edit
	 * @description Populate the HTML/DOM edit form for localized label input elements.
	 * @param {Object} form – The HTML/DOM form element to populate.
	 * @param {Object} data – The GeoJSON Feature object used to populate the form
        * @return {null}
        */		   
	populate_form_labels: function(form, data){

	    const _self = self;

	    // Count languages and display
	    
	    const labels_group = form.querySelector("#localized-labels-group")
	    const labels_desc = form.querySelector("#localized-labels-description");
	    const labels_count = form.querySelector("#localized-labels-count");

	    var count_labels = 0;
	    
	    for (var k in data.properties){

		if (k.startsWith("name:")){
		    count_labels += 1;
		}
	    }

	    switch (count_labels){
		case 0:
		    break;
		case 1:
		    labels_count.innerText = "1 language";
		    labels_desc.style.display = "inline-block";
		    break;
		default:
		    labels_count.innerText = count_labels + " languages";
		    labels_desc.style.display = "inline-block";
		    break;
	    }		    

	    // Adding and removing languages
	    
	    const t = document.querySelector("#localized-labels-row");
	    
	    const labels_selected = form.querySelector("#localized-labels-selected");

	    const tags_selected = new Tagify(labels_selected, {
		whiteList: lang_ok,
	    });

	    // START OF...

	    tags_selected.on("add", function(e){
		
		const lang = e.detail.data.value;
		const id = "localized-labels-row-" + lang;
		
		console.debug("Add label(s)", lang, id);

		// Note: We are counting on Tagify.js deduping and not dispatching
		// the same language tag twice to this event.
		    
		    try {
			const node = t.content.cloneNode(true);
			const row = node.querySelector(".localized-labels-row");
			
			row.setAttribute("id", id);
			row.setAttribute("data-lang", lang);
		    
			const label_el = row.querySelector(".lang-label");
			label_el.innerText = lang;	// LOOKUP LABEL HERE
			
			const code_el = row.querySelector(".lang-code");
			code_el.innerText = lang;
			
			const input_els = row.querySelectorAll("input");
			const count_els = input_els.length;

			if (! current_labels_langs.includes(lang)){
			    current_labels_langs.push(lang);
			}
			
			// START OF if there's a built-in method to do this
			// I can't find it...
			    
			const list_tags = function(tag){

			    const tag_els = tag.getTagElms();
			    const count_tags = tag_els.length;
			    
			    var tags = [];
			    
			    for (var t=0; t < count_tags; t++){
				tags.push(tag_els[t].getAttribute("value"));
			    }

			    return tags;
			};

			// END OF if there's a built-in method to do this
			
			for (var i=0; i < count_els; i++){

			    const el = input_els[i];
			    const spec = el.getAttribute("data-spec");

			    const el_id = "label:" + lang + "_x_" + spec;
			    el.setAttribute("label", el_id);
			    el.setAttribute("id", el_id);

			    const el_tag = new Tagify(el);

			    // Add label tags
			    
			    if (el_id in data.properties){
				console.debug("Add language tags", el_id, data.properties[el_id]);
				el_tag.addTags(data.properties[el_id]);
			    }

			    // When new language tags are added (or updated)

			    const on_edit = function(e){
				
				const tags = list_tags(el_tag);

				if (tags.length == 0){
				    return false;
				}
				
				console.debug("Edit language tag", el_id, tags);

				_self.start_spinner();
				_self.load_data().then((data) => {

				    data.properties[el_id] = tags;
				    
				    _self.save_data(data).then(() => {
					console.debug("Saved label tags on edit", el_id, tags);
					_self.stop_spinner();
				    }).catch((err) => {
					_self.stop_spinner();
					console.error("Failed to save label tags on edit", el_id, err);
					_self.alert("Failed to save data, " + err);
				    });
				    
				}).catch((err) => {
				    _self.stop_spinner();
				    console.error("Failed to load data for label tags on edit", el_id, err);				    
				    _self.alert("Failed to load data, " + err)
				});
			    };
			    
			    el_tag.on("add", on_edit);
			    el_tag.on("change", on_edit);

			    // When existing language tags are removed
			    
			    el_tag.on("remove", function(e){

				if (! e.detail.data){
				    return;
				}

				const tags = list_tags(el_tag);
				console.debug("Remove label tags", el_id, tags.length, tags);

				_self.start_spinner();
				_self.load_data().then((data) => {

				    if (tags.length == 0){

					if (el_id in data.properties){
					    console.debug("Delete label tag from properties", el_id);
					    delete(data.properties[el_id]);
					}
					
				    } else {
					data.properties[el_id] = tags;
				    }
				    
				    _self.save_data(data).then(() => {
					console.debug("Updated label tags on remove", el_id, tags);
					_self.stop_spinner();
				    }).catch((err) => {
					_self.stop_spinner();
					console.error("Failed to update label tags on remove", el_id, err);
					_self.alert("Failed to update label tags on remove, " + err);
				    });
				    
				}).catch((err) => {
				    _self.stop_spinner();
					console.error("Failed to load label tags on remove", el_id, err);				    
				    _self.alert("Failed to load data, " + err)
				});
				
			    });
			    
			}

			labels_group.prepend(row);
			
		    } catch (err) {
			console.error("Failed to add label(s) row for ", lang, id, err);
		    }
	    });

	    tags_selected.on("remove", function(e){

		if (!e.detail.data){
		    return;
		}

		const lang = e.detail.data.value;
		const id = "localized-labels-row-" + lang;
		
		console.debug("Remove label(s)", lang, id);

		const row = document.getElementById(id);

		if (! row){
		    console.warn("Failed to locate label(s) row for language", lang, id);
		    return;
		}

		row.parentNode.removeChild(row);

		if (current_labels_langs.includes(lang)){

		    current_labels_langs = current_labels_langs.filter(function(item) {
			return item !== lang;
		    });
		}
		
	    });

	    // END OF...
	    
	    console.debug("Add current label tags")
	    tags_selected.addTags(current_labels_langs);	    
	},

	/**
	 * @function populate_form_concordances
	 * @memberof wof.edit
	 * @description Populate the HTML/DOM edit form for wof:concordances input elements.
	 * @param {Object} form – The HTML/DOM form element to populate.
	 * @param {Object} data – The GeoJSON Feature object used to populate the form
         * @return ...
        */		   						
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
		
		close.innerHTML = close_svg;
		
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

		    _self.load_data().then((data) => {

			if (prefix in data.properties["wof:concordances"]){
			    _self.stop_spinner();
			    _self.alert("There is already a concordance with that prefix");
			    return false;
			}
			
			data.properties["wof:concordances"][prefix] = value;
			
			_self.save_data(data).then(() => {
			    
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

	/**
	 * @function populate_existential_flag
	 * @memberof wof.edit
	 * @description Populate an "existential flag" DOM element
	 * @param {Object} input_el – The HTML/DOM input element to populate.
	 * @param {Object} flag_v – The current value of the WOF property for input_el.	   
         * @return {null}
        */		   						
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

	/**
	 * @function update_concordance
	 * @memberof wof.edit
	 * @description Update the value of the (wof:)concordance matching prefix.
	 * @param {string} prefix – The wof:concordances prefix to update.
	 * @param {string} v – The updated value for the wof:concordances prefix.	   
         * @return {Promise}
        */		   						
	update_concordance: function(prefix, v){

	    const _self = self;

	    return new Promise((resolve, reject) => {
		
		_self.start_spinner();

		_self.load_data().then((data) => {

		    data.properties["wof:concordances"][prefix] = v;
		    console.debug("Update concordances", data.properties["wof:concordances"]);

		    _self.save_data(data).then(() => {
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
		    console.error("Failed to parse raw data", err);
		    _self.stop_spinner();
		    
		    reject("Failed to parse raw data, " + err);
		    return;
		});
	    	
	    });
	},
	
	/**
	 * @function remove_concordance
	 * @memberof wof.edit
	 * @description Remove the concordance matching 'prefix' from the underlying #raw HTML element data. If successful the corresponding HTML element for that prefix will be removed from the edit form DOM.
	 * @param {string} prefix – The wof:concordances namespace and prefix to remove.
         * @return {Promise}
        */		   							
	remove_concordance: function(prefix){

	    const _self = self;

	    return new Promise((resolve, reject) => {
		
		_self.start_spinner();

		_self.load_data().then((data) => {

		    delete(data.properties["wof:concordances"][prefix]);

		    _self.save_data(data).then(() => {

			try {
			    const row = document.getElementById("wof-concordances-" + prefix);
			    const parent = row.parentNode;
			    parent.removeChild(row);
			} catch(err) {
			    console.error("Failed to remove wof:concordaces row", prefix, err);
			    _self.stop_spinner();
			    
			    reject("Failed to remove row, " + err);		    
			    return false;
			}
			
			_self.stop_spinner();
			
			resolve();
			return;
		    }).catch((err) => {
			console.error("Failed to save wof:concordances data", prefix, err);
			_self.stop_spinner();
			
			reject("Failed to save concordances data, " + err);		    
			return;
		    });

		}).catch((err) => {
		    console.error("Failed to load data", err);
		    _self.stop_spinner();
		    
		    reject("Failed to load data, " + err);		    
		    return;
		});	    	    
	    });
	},

	/**
	 * @function start_spinner
	 * @memberof wof.edit
	 * @description Display the "spinner" UI element.
         * @return {null}
        */		   								
	start_spinner: function(){
		const spinner = document.querySelector("#spinner-svg");
		spinner.style.display = "inline-block";
	},

	/**
	 * @function stop_spinner
	 * @memberof wof.edit
	 * @description Hide the "spinner" UI element.
         * @return {null}
        */		   									
	stop_spinner: function(){
	    const spinner = document.querySelector("#spinner-svg");
	    spinner.style.display = "none";
	},

	/**
	 * @function save_data
	 * @memberof wof.edit
	 * @description Validate and format data and then store the result in the #raw HTML <pre> element
	 * @param {Object} data - The GeoJSON Feature object data to save.
         * @return {Promise} -
        */		   										
	save_data: function(data) {

	    const _self = self;
	    
	    return new Promise((resolve, reject) => {

		const raw = document.querySelector("#raw");

		if (! raw){
		    reject("Missing #raw data element.");
		}	
		
		self.export_data(data).then((fmt_rsp) => {
		    raw.innerText = fmt_rsp;
		    resolve();
		}).catch((err) => {
		    console.error("Failed to export data", err);
		    reject("Failed to export data, " + err);
		});
	    });
	},

	/**
	 * @function export_data
	 * @memberof wof.edit
	 * @description Validate and format a GeoJSON Feature object.
	 * @param {Object} data - A GeoJSON Feature object.
         * @return {Promise.<string>} - Returns the string representation of 'data' after it has been validated and formatted.
        */		   											
	export_data: function(data){

	    return new Promise((resolve, reject) => {
		
	    	const str_data = JSON.stringify(data);

		// Basic sanity checking
		wof_validate(str_data).then(() => {

		    // Prep various things like derived geometries and relations
		    wof_prepare_feature(str_data).then((str_prepped) => {

			// Make pretty
			wof_format(str_prepped).then((fmt_rsp) => {
			    resolve(fmt_rsp);
			}).catch((err) => {
			    console.error("Data formatting failed", str_data, err);			
			    reject("Data formatting failed, " + err);
			});
			
		    }).catch((err) => {
			console.error("Data prepping failed", data, err);
		    });
		    
		}).catch((err) => {
		    console.error("Data validation failed", data, err);
		    reject("Data validation failed, " + err);
		});
	    });
	},

	/**
	 * @function load_data
	 * @memberof wof.edit
	 * @description Returns the parsed JSON of the GeoJSON Feature data written to the #raw <pre> HTML element.
         * @return {Promise.<Object>} - A GeoJSON Feature object.
        */		   											
	load_data: function(){

	    return new Promise((resolve, reject) => {
		
		const raw = document.querySelector("#raw");

		if (! raw){
		    console.error("Missing #raw data element");
		    reject("Missing #raw data element.");
		}	
		
		try {
		    const data = JSON.parse(raw.innerText);
		    resolve(data);
		} catch(err) {
		    console.error("Failed to parse #raw data", err);
		    reject(err);
		}
	    });
	},

	/**
	 * @function save_geometry
	 * @memberof wof.edit
	 * @description Write updated GeoJSON geometry back to its parent Feature element
	 * @param {Object} geom – A GeoJSON Feature geometry
         * @return {null}
        */		   					
	save_geometry: function(geom){

	    const _self = self;
	    
	    self.load_data().then((data) => {

		data.geometry = geom;
		
		_self.save_data(data).then(() => {
		    console.debug("Geometry updated");
		    return;
		}).catch((err) => {
		    console.error("Failed to save updated geometry", err);
		    _self.alert("Failed to save updated geometry, " + err);
		});
		
	    }).catch((err) => {
		console.error("Failed to load data for updating geometry", err);		
		self.alert("Failed to load data for updating geometry, " + err);
	    });
	    
	},

	/**
	 * @function update_geometry
	 * @memberof wof.edit
	 * @description Derive an update geometry from Leaflet/geoman features and update the parent Feature element with that geometry.
	 * @param {Object} map – A Leaflet L.Map instance
         * @return {null}
        */		   						
	update_geometry: function(map){
	    
	    const feature_group = map.pm.getGeomanLayers(true);
	    const feature_collection = feature_group.toGeoJSON();

	    const features = feature_collection.features;
	    const count = features.length;

	    console.debug("Update geometry, feature count", count);
	    
	    var geom;
	    
	    const _self = self;

	    switch (count) {
		case 0:
		    _self.alert("Geometry can no be empty.");
		    break;
		case 1:

		    geom = features[0].geometry;
		    _self.save_geometry(geom);
		    break;
		    
		default:
		    // merge geometries here...
		    // based on the controls (above) we can have
		    // points
		    // polygons

		    var point_geoms = [];
		    var poly_geoms = [];

		    for (var i=0; i < count; i++){

			const geom = features[i].geometry;

			switch (geom.type){
			    case "Polygon":
				poly_geoms.push(geom.coordinates);
				break;
			    case "MultiPolygon":

				const count_polys = geom.coordinates.length;

				for (var j=0; j < count_polys; j++){
				    poly_geoms.push(geom.coordinates[j]);
				}

				break;
				
			    case "Point":
				point_geoms.push(geom.coordinates);
				break;
			    case "MultiPoint":

				const count_points = geom.coordinates.length;

				for (var j=0; j < count_points; j++){
				    point_geoms.push(geom.coordinates[j]);
				}

				break;
				
			    default:
				_self.alert("Unhandled geometry type, " + geom.type);
				return false;
			}
		    }

		    var geom_points;
		    var geom_polys;
		    
		    const count_points = point_geoms.length;		    
		    const count_polys = poly_geoms.length;

		    console.debug("Merged feature count", "points", count_points, "polygons", count_polys);
		    
		    switch (count_points) {
			case 1:
			    geom_points = {
				type: "Point",
				coordindates: point_geoms[0]
			    };
			    break;
			default:
			    geom_points = {
				type: "MultiPoint",
				coordindates: point_geoms
			    };
		    }

		    switch (count_polys) {
			case 1:
			    geom_polys = {
				type: "Polygon",
				coordindates: point_geoms[0]
			    };
			    break;
			default:
			    geom_polys = {
				type: "MultiPolygon",
				coordindates: point_geoms
			    };
		    }

		    if ((count_points) && (count_polys)){
			geom = {
			    type: "MultiGeometry",
			    "geometries": [
				geom_points,
				geom_polys,
			    ]
			};
		    } else if (count_points){
			geom = geom_points;
		    } else {
			geom = geom_polys;
		    }
				
		    console.debug("Save geometry", geom.type);
		    _self.save_geometry(geom);
		    break;
	    }
	},
	
	/**
	 * @function feedback
	 * @memberof wof.edit
	 * @description Display a feedback message.
	 * @param {string} msg - The feedback message to display.
	 * @param {number} ttl - The number of milliseconds to display the message after which is will be removed. Default is 5000.
         * @return {null}
         */		   		
	feedback: function(msg, ttl){
	    
	    if (! ttl){
		ttl = 5000;
	    }
	    
	    self.alert(msg, ttl);
	},
	
	/**
	 * @function alert
	 * @memberof wof.edit
	 * @description Display an alert message.
	 * @param {string} msg - The message to display.
	 * @param {number} ttl - The number of milliseconds to display the message after which is will be removed.
         * @return {null}
         */		   			
	alert: function(msg, ttl){
	    
	    console.debug(msg);
	    
	    const dlg = document.createElement("dialog");
	    dlg.setAttribute("id", "alert");
	    dlg.setAttribute("class", "alert");
	    
	    const close = document.createElement("div");
	    close.setAttribute("class", "alert-close");
	    
	    close.innerHTML = close_svg;	    
	    
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
	
    };
    
    return self;
    
})();
