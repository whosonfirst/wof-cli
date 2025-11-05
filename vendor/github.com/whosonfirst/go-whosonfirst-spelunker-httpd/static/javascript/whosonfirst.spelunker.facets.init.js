window.addEventListener("load", function load(event){

    
    var facets_wrapper = document.querySelector("#whosonfirst-facets");

    if (! facets_wrapper){
	console.log("Missing #whosonfirst-facets wrapper");
	return;
    }

    var current_url = facets_wrapper.getAttribute("data-current-url");
    var facets_url = facets_wrapper.getAttribute("data-facets-url");    

    // console.log("FACETS", facets_url);
    
    if ((! current_url) || (! facets_url)){
	return;
    }

    var draw_facets = function(rsp){

	var f = rsp.facet.property;
	var id = "#whosonfirst-facets-" + f;

	var el = document.querySelector(id);

	if (! el){
	    console.log("Unable to find facet wrapper", id);
	    return;
	}

	var f_label = f;

	switch (f) {
	    case "iscurrent":
		f_label = "is current";
		break;
	    case "isdeprecated":
		f_label = "is deprecated";

		var results = rsp.results;
		var count = results.length;

		deprecated = 0;
		not_deprecated = 0;
		
		for (var i=0; i < count; i++){

		    switch (parseInt(results[i].key)){
			case 1:
			    deprecated = results[i].count;
			    break;
			default:
			    not_deprecated += results[i].count;	
			    break;
		    }
		}

		rsp.results = [
		    { key: "0", count: not_deprecated },
		    { key: "1", count: deprecated },		    
		];

		break;
		
	    default:
		break;
	}
	
	var label = document.createElement("h3");
	label.appendChild(document.createTextNode(f_label));
	
	var ul = document.createElement("ul");
	ul.setAttribute("class", "whosonfirst-facets");
	
	var results = rsp.results;
	var count = results.length;

	for (var i=0; i < count; i++){

	    var k = results[i].key;

	    if (k == ""){

		var sp = document.createElement("span");
		sp.setAttribute("class", "hey-look");
		sp.appendChild(document.createTextNode("undefined"));

		var facet_count = results[i].count;
		facet_count = Intl.NumberFormat().format(facet_count);
		
		var sm = document.createElement("small");
		sm.appendChild(document.createTextNode(facet_count));
		
		var item = document.createElement("li");
		item.appendChild(sp);
		item.appendChild(sm);

	    } else {

		var k_label = k;

		switch(f){
		    case "iscurrent":
			
			switch (parseInt(k)){
			    case 0:
				k_label = "not current";
				break;
			    case -1:
				k_label = "unknown";
				break;
			    default:
				k_label = "current";
				break;
			}
			
			break;
			
		    case "isdeprecated":
			
			switch (parseInt(k)){
			    case 0:
				k_label = "valid";
				break;
			    default:
				k_label = "deprecated";
				break;
			}
			
			break;
		default:

			if (f_label == "country"){

			    var country = whosonfirst.spelunker.countries.by_code(k);

			    if (country){
				k_label = country["wof:name"];
			    } else {
				console.log("Unable to determine country name for " + k);
			    }
			}
		}
		
		// Something something something is location.href really safe?
		// https://developer.mozilla.org/en-US/docs/Web/API/URL/URL

		var u = new URL(current_url, location.href);
		u.searchParams.set(f, k)

		var a = document.createElement("a");
		
		a.setAttribute("href", u.toString());
		a.setAttribute("class", "hey-look");

		if (k_label == "deprecated"){
		    a.setAttribute("class", "hey-look deprecated");
		}
		
		a.appendChild(document.createTextNode(k_label));

		var facet_count = results[i].count;
		facet_count = Intl.NumberFormat().format(facet_count);
		
		var sm = document.createElement("small");
		sm.appendChild(document.createTextNode(facet_count));
		
		var item = document.createElement("li");
		item.appendChild(a);
		item.appendChild(sm);
	    }
	    
	    ul.appendChild(item);
	}

	var summary = document.createElement("summary");
	summary.appendChild(document.createTextNode(f_label));
	
	var details = document.createElement("details");
	details.setAttribute("open", "open");
	
	details.appendChild(summary);
	details.appendChild(ul);

	el.appendChild(details);
	
	// el.appendChild(label);
	// el.appendChild(ul);
    };
    
    var fetch_facet = function(f){

	// var url = facets_url + "?&facet=" + f;

	// Something something something is location.href really safe?
	// https://developer.mozilla.org/en-US/docs/Web/API/URL/URL
	
	var u = new URL(facets_url, location.href)
	u.searchParams.set("facet", f);
	var url = u.toString();

	fetch(url)
	    .then((rsp) => rsp.json())
	    .then((data) => {

		var count = data.length;

		for (var i=0; i < count; i++){
		    draw_facets(data[i]);
		}
		
	    }).catch((err) => {
		console.log("SAD", f, err);
	    });
    };
    
    var facets = facets_wrapper.getAttribute("data-facets");
    facets = facets.split(",");

    var count_facets = facets.length;

    for (var i=0; i < count_facets; i++){

	var f = facets[i];
	
	var el = document.createElement("div");
	el.setAttribute("id", "whosonfirst-facets-" + f);
	facets_wrapper.appendChild(el);

	fetch_facet(f);
    }

});
