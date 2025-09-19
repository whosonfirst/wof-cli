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
		    console.log("Jump to show here...")
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

	    wof.edit.api.fetch(uri).then((rsp) => {
		console.log("OK", rsp);
	    }).catch((err) => {
		console.error("SAD", err)
	    });
	},
    };
    
    return self;
})();
