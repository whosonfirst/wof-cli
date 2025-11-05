var whosonfirst = whosonfirst || {};
whosonfirst.spelunker = whosonfirst.spelunker || {};

// UPDATE TO USE promises

whosonfirst.spelunker.cache = (function(){
    
    var ttl = 30000;

    var self = {
	
	'get': function(key, on_hit, on_miss){

	    if (typeof(localforage) != 'object'){
		return false;
	    }

	    var fq_key = self.prep_key(key);

	    localforage.getItem(fq_key, function (err, rsp){

		console.log("GET", fq_key, err, rsp);
		
		if ((err) || (! rsp)){
		    // console.log("cache MISS for " + fq_key);
		    on_miss();
		    return false;
		}

		// console.log("cache HIT for " + fq_key);
		// console.log(rsp);

		var data = rsp['data'];

		if (! data){
		    // console.log("cache WTF for " + fq_key);
		    on_miss();
		    return false;
		}

		var dt = new Date();
		var ts = dt.getTime();

		var then = rsp['created'];
		var diff = ts - then;

		if (diff > ttl){
		    // console.log("cache EXPIRED for " + fq_key);
		    self.unset(key);
		    on_miss();
		    return false;
		}

		on_hit(data);
	    });

	    return true;
	},

	'set': function(key, value){

	    if (typeof(localforage) != 'object'){
		return false;
	    }

	    var dt = new Date();
	    var ts = dt.getTime();

	    var wrapper = {
		'data': value,
		'created': ts
	    };

	    key = self.prep_key(key);
	    // console.log("cache SET for " + key);

	    localforage.setItem(key, wrapper);
	    return true;
	},

	'unset': function(key){

	    if (typeof(localforage) != 'object'){
		return false;
	    }

	    key = self.prep_key(key);
	    // console.log("cache UNSET for " + key);

	    localforage.removeItem(key);
	    return true;
	},

	'prep_key': function(key){
	    return key;
	}
    };

    return self;

})();
