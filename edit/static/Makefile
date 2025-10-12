BOOTSTRAP_VERSION=5.3.8
LEAFLET_VERSION=1.9.4
GEOMAN_VERSION=2.18.2
bootstrap:
	curl -s -o css/bootstrap.min.css https://cdn.jsdelivr.net/npm/bootstrap@$(BOOTSTRAP_VERSION)/dist/css/bootstrap.min.css
	curl -s -o javascript/bootstrap.bundle.min.js https://cdn.jsdelivr.net/npm/bootstrap@$(BOOTSTRAP_VERSION)/dist/js/bootstrap.bundle.min.js

leaflet:
	curl -s -o css/leaflet.css https://unpkg.com/leaflet@$(LEAFLET_VERSION)/dist/leaflet.css
	curl -s -o javascript/leaflet.js https://unpkg.com/leaflet@$(LEAFLET_VERSION)/dist/leaflet.js

# https://github.com/yairEO/tagify
tagify:
	curl -s -o css/tagify.css https://unpkg.com/@yaireo/tagify@2.31.6/dist/tagify.css
	curl -s -o javascript/tagify.js https://raw.githubusercontent.com/yairEO/tagify/refs/heads/master/dist/tagify.js
	curl -s -o javascript/tagify.js.map https://raw.githubusercontent.com/yairEO/tagify/refs/heads/master/dist/tagify.js.map

geoman:
	curl -s -o css/leaflet-geoman.css https://unpkg.com/@geoman-io/leaflet-geoman-free@$(GEOMAN_VERSION)/dist/leaflet-geoman.css
	curl -s -o javascript/leaflet-geoman.js https://unpkg.com/@geoman-io/leaflet-geoman-free@$(GEOMAN_VERSION)/dist/leaflet-geoman.js

wofjs:
	curl -s -o javascript/whosonfirst.bundle.min.js https://raw.githubusercontent.com/whosonfirst/js-whosonfirst/refs/heads/main/dist/whosonfirst.bundle.min.js

wasmjs:
	curl -s -o javascript/sfomuseum.golang.wasm.bundle.js https://raw.githubusercontent.com/sfomuseum/js-sfomuseum-golang-wasm/refs/heads/main/dist/sfomuseum.golang.wasm.bundle.js
