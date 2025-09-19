var wof = wof || {};
wof.edit = wof.edit || {};

wof.edit.api = (function () {

    var self = {

	list: function() {
	    return new Promise((resolve, reject) => {
		
		fetch("/api/list").then(rsp =>
		    rsp.json()
		).then((data) => {
		    resolve(data);
		}).catch((err) => {
		    reject(err);
		});
	    });
	},

	fetch: function(uri){
	    return new Promise((resolve, reject) => {
		
		fetch("/data/" + uri).then(rsp =>
		    rsp.json()
		).then((data) => {
		    resolve(data);
		}).catch((err) => {
		    reject(err);
		});
	    });
	}
	
    };
    
    return self;
})();
