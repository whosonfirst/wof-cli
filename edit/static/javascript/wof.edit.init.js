window.addEventListener("load", function load(event) {

    wof.edit.init().then((rsp) => {
	console.log("GO");
    }).catch((err) => {
	console.error("Failed to initialize", err);
    });
    
});
