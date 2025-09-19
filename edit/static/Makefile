BOOTSTRAP_VERSION=5.3.8
LEAFLET_VERSION=1.9.4

bootstrap:
	curl -s -o css/bootstrap.min.css https://cdn.jsdelivr.net/npm/bootstrap@$(BOOTSTRAP_VERSION)/dist/css/bootstrap.min.css
	curl -s -o javascript/bootstrap.bundle.min.js https://cdn.jsdelivr.net/npm/bootstrap@$(BOOTSTRAP_VERSION)/dist/js/bootstrap.bundle.min.js

leaflet:
	curl -s -o css/leaflet.css https://unpkg.com/leaflet@$(LEAFLET_VERSION)/dist/leaflet.css
	curl -s -o javascript/leaflet.js https://unpkg.com/leaflet@$(LEAFLET_VERSION)/dist/leaflet.js

wofjs:
	curl -s -o javascript/whosonfirst.bundle.min.js https://raw.githubusercontent.com/whosonfirst/js-whosonfirst/refs/heads/main/dist/whosonfirst.bundle.min.js

wasmjs:
	curl -s -o javascript/sfomuseum.golang.wasm.bundle.js https://raw.githubusercontent.com/sfomuseum/js-sfomuseum-golang-wasm/refs/heads/main/dist/sfomuseum.golang.wasm.bundle.js

# curl -s -I https://github.com/whosonfirst/go-whosonfirst-validate/raw/refs/heads/main/www/wasm/wof_validate.wasm | grep location
wasm-validate:
	curl -s -o wasm/wof_validate.wasm https://raw.githubusercontent.com/whosonfirst/go-whosonfirst-validate/refs/heads/main/www/wasm/wof_validate.wasm

# curl -s -I https://github.com/whosonfirst/go-whosonfirst-format/raw/refs/heads/main/www/wasm/wof_format.wasm | grep location
wasm-format:
	curl -s -L -o wasm/wof_format.wasm https://raw.githubusercontent.com/whosonfirst/go-whosonfirst-format/refs/heads/main/www/wasm/wof_format.wasm
