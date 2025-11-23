BOOTSTRAP_VERSION=5.3.8
LEAFLET_VERSION=1.9.4
PROTOMAPS_LEAFLET_VERSION4=4.0.1
PROTOMAPS_LEAFLET_VERSION5=5.1.0

bootstrap:
	curl -s -o css/bootstrap.min-$(BOOTSTRAP_VERSION).css https://cdn.jsdelivr.net/npm/bootstrap@$(BOOTSTRAP_VERSION)/dist/css/bootstrap.min.css
	curl -s -o javascript/bootstrap.bundle.min-$(BOOTSTRAP_VERSION).js https://cdn.jsdelivr.net/npm/bootstrap@$(BOOTSTRAP_VERSION)/dist/js/bootstrap.bundle.min.js

leaflet:
	curl -s -o css/leaflet-$(LEAFLET_VERSION).css https://unpkg.com/leaflet@$(LEAFLET_VERSION)/dist/leaflet.css
	curl -s -o javascript/leaflet-$(LEAFLET_VERSION).js https://unpkg.com/leaflet@$(LEAFLET_VERSION)/dist/leaflet.js

protomaps:
	curl -s -o javascript/protomaps-leaflet-$(PROTOMAPS_LEAFLET_VERSION4).js https://unpkg.com/protomaps-leaflet@$(PROTOMAPS_LEAFLET_VERSION4)/dist/protomaps-leaflet.js
	curl -s -o javascript/protomaps-leaflet-$(PROTOMAPS_LEAFLET_VERSION5).js https://unpkg.com/protomaps-leaflet@$(PROTOMAPS_LEAFLET_VERSION5)/dist/protomaps-leaflet.js
