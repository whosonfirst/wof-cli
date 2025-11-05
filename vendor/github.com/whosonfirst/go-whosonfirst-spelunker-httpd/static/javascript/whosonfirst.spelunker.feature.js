var whosonfirst = whosonfirst || {};
whosonfirst.spelunker = whosonfirst.spelunker || {};

whosonfirst.spelunker.feature = (function(){
    
    var cache_ttl = 30000;
    
    var self = {
	
	fetch: function(wofid, uri_args){

	    if (wofid < 0){

		return new Promise((resolve, reject) => {
		    reject("Not a valid WOF ID to fetch");
		});
	    }
	    
	    var _self = self;
	    var _url = whosonfirst.spelunker.uri.id2abspath(wofid, uri_args);

	    return new Promise((resolve, reject) => {	

		var on_hit = function(f){
		    console.log("FETCH HIT", _url, f);
		    resolve(f);
		};
		
		var on_miss = function(){
		    console.log("FETCH MISS", _url);
		    
		    _self._refresh(_url).then((rsp) => {
			resolve(rsp);
		    }).catch((err) => {
			reject(err);
		    });
		};

		console.log("FETCH CACHE", _url);
		whosonfirst.spelunker.cache.get(_url, on_hit, on_miss);
	    });
	    
	},

	'_refresh': function(url){

	    var _self = self;
	    var _url = url;

	    console.log("_REFRESH", url);
	    
	    return new Promise((resolve, reject) => {
		
		fetch(_url).then((rsp) => rsp.json())
			   .then((feature) => {
			       console.log("_REFRESH OK", _url);
			      whosonfirst.spelunker.cache.set(_url, feature);
			      resolve(feature);
			   }).catch((err) => {
			       console.log("_REFRESH ERR", _url, err);			       
			      console.log("Failed to fetch source", _url, err);
			      reject(err);
			  });
	    });

	}
    };

    return self;

})();
