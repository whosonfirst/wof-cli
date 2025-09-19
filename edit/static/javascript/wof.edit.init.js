window.addEventListener("load", function load(event) {

    wof.edit.init().then((rsp) => {

	wof.edit.list();
	
    }).catch((err) => {
	console.error("Failed to initialize", err);
    });
    
});
