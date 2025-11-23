var whosonfirst = whosonfirst || {};
whosonfirst.spelunker = whosonfirst.spelunker || {};

// UPDATE TO USE whosonfirst.spelunker.cache

whosonfirst.spelunker.namify = (function(){
    
    var cache_ttl = 30000;

    var self = {
	
	'init': function(){

	},
	
	namify_selector: function(q){
	    
	    var resolver = whosonfirst.spelunker.uri.id2abspath;

	    var els = document.querySelectorAll(q);
	    var count = els.length;

	    for (var i=0; i < count; i++){
		self.namify_el(els[i], resolver);
	    }
	},

	namify_el: function(el, resolver){

	    var wofid = el.getAttribute("data-wof-id");

	    if (! wofid){	
		console.log("node is missing data-wof-id attribute");
		return;
	    }

	    if (el.textContent != wofid){
		console.log("node has not-a-wof-id body");
		return;
	    }

	    var url = resolver(wofid);

	    var on_hit = function(feature){
		// console.log("NAMIFY FROM CACHE");
		self.apply_namification(el, feature);
	    };
	    
	    var on_miss = function(){
		// console.log("INVOKING ON MISS FOR " + url);
		self.namify_el_from_source(url, el);
	    };

	    if (! self.cache_get(url, on_hit, on_miss)){
		self.namify_el_from_source(url, el);
	    }

	},

	'namify_el_from_source': function(url, el){

	    var _self = self;
	    var _el = el;
	    var _url = url;
	    
	    fetch(url).then((rsp) => rsp.json())
		      .then((feature) => {
			  _self.apply_namification(_el, feature);
			  _self.cache_set(_url, feature);		
		      }).catch((err) => {
			  console.log("Failed to fetch source", _url, err);
		      });

	},

	'apply_namification': function(el, feature){

		var props = feature['properties'];

		// to account for whosonfirst-brands which needs to be updated
		// to grow a 'properties' hash... (20160319/thisisaaronland)

		if (! props){
		    props = feature;
		}

	        // console.log(props);
	    
		var label = props['wof:label'];

		if ((! label) || (label == '')){

		    var possible = [
			'wof:name',
			'wof:brand_name'
		    ];

		    var count = possible.length;

		    for (var i = 0; i < count; i++) {

			var k = possible[i];

			if (props[k]){
			    label = props[k];
			    break;
			}
		    }
		}

		var enc_label = whosonfirst.spelunker.php.htmlspecialchars(label);
		el.innerHTML = enc_label;
	},

	'cache_get': function(key, on_hit, on_miss){

	    if (typeof(localforage) != 'object'){
		return false;
	    }

	    var fq_key = self.cache_prep_key(key);

	    localforage.getItem(fq_key, function (err, rsp){

		// console.log("GET", fq_key, err, rsp);
		
		if ((err) || (! rsp)){
		    // console.log("cache MISS for " + fq_key);
		    on_miss();
		    return;
		}

		// console.log("cache HIT for " + fq_key);
		// console.log(rsp);

		var data = rsp['data'];

		if (! data){
		    // console.log("cache WTF for " + fq_key);
		    on_miss();
		    return;
		}

		var dt = new Date();
		var ts = dt.getTime();

		var then = rsp['created'];
		var diff = ts - then;

		if (diff > cache_ttl){
		    // console.log("cache EXPIRED for " + fq_key);
		    self.cache_unset(key);
		    on_miss();
		    return;
		}

		on_hit(data);
	    });

	    return true;
	},

	'cache_set': function(key, value){

	    if (typeof(localforage) != 'object'){
		return false;
	    }

	    var dt = new Date();
	    var ts = dt.getTime();

	    var wrapper = {
		'data': value,
		'created': ts
	    };

	    key = self.cache_prep_key(key);
	    // console.log("cache SET for " + key);

	    localforage.setItem(key, wrapper);
	    return true;
	},

	'cache_unset': function(key){

	    if (typeof(localforage) != 'object'){
		return false;
	    }

	    key = self.cache_prep_key(key);
	    // console.log("cache UNSET for " + key);

	    localforage.removeItem(key);
	    return true;
	},

	'cache_prep_key': function(key){
	    return key + '#namify';
	}
    };

    return self;

})();
