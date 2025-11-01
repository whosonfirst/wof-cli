BOOTSTRAP_VERSION=5.3.8
GEOMAN_VERSION=2.18.2
LEAFLET_VERSION=1.9.4
PROTOMAPS_LEAFLET_VERSION4=4.0.1
PROTOMAPS_LEAFLET_VERSION5=5.1.0
TAGIFY_VERSION=2.31.6

bootstrap:
	curl -s -o css/bootstrap.min-$(BOOTSTRAP_VERSION).css https://cdn.jsdelivr.net/npm/bootstrap@$(BOOTSTRAP_VERSION)/dist/css/bootstrap.min.css
	curl -s -o javascript/bootstrap.bundle.min-$(BOOTSTRAP_VERSION).js https://cdn.jsdelivr.net/npm/bootstrap@$(BOOTSTRAP_VERSION)/dist/js/bootstrap.bundle.min.js

leaflet:
	curl -s -o css/leaflet-$(LEAFLET_VERSION).css https://unpkg.com/leaflet@$(LEAFLET_VERSION)/dist/leaflet.css
	curl -s -o javascript/leaflet-$(LEAFLET_VERSION).js https://unpkg.com/leaflet@$(LEAFLET_VERSION)/dist/leaflet.js

protomaps:
	curl -s -o javascript/protomaps-leaflet-$(PROTOMAPS_LEAFLET_VERSION4).js https://unpkg.com/protomaps-leaflet@$(PROTOMAPS_LEAFLET_VERSION4)/dist/protomaps-leaflet.js
	curl -s -o javascript/protomaps-leaflet-$(PROTOMAPS_LEAFLET_VERSION5).js https://unpkg.com/protomaps-leaflet@$(PROTOMAPS_LEAFLET_VERSION5)/dist/protomaps-leaflet.js

# https://github.com/yairEO/tagify
tagify:
	curl -s -o css/tagify-$(TAGIFY_VERSION).css https://unpkg.com/@yaireo/tagify@$(TAGIFY_VERSION)/dist/tagify.css
	curl -s -o javascript/tagify-$(TAGIFY_VERSION).js https://raw.githubusercontent.com/yairEO/tagify/refs/heads/master/dist/tagify.js

geoman:
	curl -s -o css/leaflet-geoman-$(GEOMAN_VERSION).css https://unpkg.com/@geoman-io/leaflet-geoman-free@$(GEOMAN_VERSION)/dist/leaflet-geoman.css
	curl -s -o javascript/leaflet-geoman-$(GEOMAN_VERSION).js https://unpkg.com/@geoman-io/leaflet-geoman-free@$(GEOMAN_VERSION)/dist/leaflet-geoman.js

wofjs:
	curl -s -o javascript/whosonfirst.bundle.min.js https://raw.githubusercontent.com/whosonfirst/js-whosonfirst/refs/heads/main/dist/whosonfirst.bundle.min.js

wasmjs:
	curl -s -o javascript/sfomuseum.golang.wasm.bundle.js https://raw.githubusercontent.com/sfomuseum/js-sfomuseum-golang-wasm/refs/heads/main/dist/sfomuseum.golang.wasm.bundle.js
