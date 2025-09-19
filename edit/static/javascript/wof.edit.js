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
		    })
		    
		}).catch((err) => {
		    reject("Failed to load update wof_format WASM binary " + err);
		});

		});
	},
    };
    
    return self;
})();
