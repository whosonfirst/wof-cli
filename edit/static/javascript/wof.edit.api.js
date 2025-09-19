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
	},

	save: function(uri, body) {

	    return new Promise((resolve, reject) => {

		const onload = function (rsp) {
		    const target = rsp.target;
		    
		    if (target.readyState != 4) {
			return;
		    }
		    
		    const status_code = target["status"];
		    const status_text = target["statusText"];
		    
		    if (status_code != 200) {
			reject(status_text);
			return;
		    }
		    
		    const raw = target["responseText"];
		    const do_rsp = new Response(raw);
		    
		    resolve(do_rsp);
		    return true;
		};
		
		const onerror = function (err) {
		    reject(err);
		};
		
		const onabort = function (err) {
		    reject(err);
		};
		
		const req = new XMLHttpRequest();
		req.addEventListener("load", onload);
		req.addEventListener("error", onerror);
		req.addEventListener("abort", onabort);
		
		req.open("POST", "/api/save/" + uri, true);
		req.send(body);
	    });	    
	},
	
    };
    
    return self;
})();
