{{ define "whosonfirst_spelunker_uris" -}}
var whosonfirst = whosonfirst || {};
whosonfirst.spelunker = whosonfirst.spelunker || {};

whosonfirst.spelunker.uris = (function(){

    var _table = {{ .Table }};

    var self = {
	
	abs_root_url: function(){

	    var root = _table.root_url;

	    if (! root.endsWith("/")){
		root += "/";
	    }
	    
	    return root;
	},
	
	table: function(){
	    return _table;
	},
    };

    return self;
})();
{{ end -}}
