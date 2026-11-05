include ../../.env

email ?= didi@test.com

.PHONY: city-create
city-create:
	@curlie -v POST ":$(APP_PORT)/city" email=$(email) password=sisilafamille name=new-city mayor_name=dani

.PHONY: city-delete
city-delete:
	@curlie -v DELETE ":$(APP_PORT)/city/$(id)"

.PHONY: update-zoning-industrial
update-zoning-industrial:
	@curlie -v PATCH ":$(APP_PORT)/city/$(id)/zoning/update" super_kind_name=Industrial top_left:='{"q": 0, "r": 0}' bottom_right:='{"q": 31, "r": 31}'

.PHONY: update-zoning-residential
update-zoning-residential:
	@curlie -v PATCH ":$(APP_PORT)/city/$(id)/zoning/update" super_kind_name=Residential top_left:='{"q": 32, "r": 0}' bottom_right:='{"q": 63, "r": 31}'

.PHONY: revert-zoning
revert-zoning:
	@curlie -v PATCH ":$(APP_PORT)/city/$(id)/zoning/revert" top_left:='{"q": 0, "r": 0}' bottom_right:='{"q": 127, "r": 31}'

.PHONY: get-tiles
get-tiles:
	@curlie -v GET ":$(APP_PORT)/debug/city/$(id)/tiles?q=60&r=0"

.PHONY: get-tiles-by-anchor
get-tiles-by-anchor:
	@curlie -v GET ":$(APP_PORT)/debug/city/$(id)/tiles-by-anchor?anchor_q=0&anchor_r=0"

.PHONY: get-city-tiles
get-city-tiles:
	@curlie -v GET ":$(APP_PORT)/city/$(id)/all-tiles"

.PHONY: destroy-artifact
destroy-artifact:
	@curlie -v PATCH ":$(APP_PORT)/city/$(id)/construction/destroy" position:='{"q": 0, "r": 0}'

.PHONY: build-artifact
build-artifact:
	@curlie -v PATCH ":$(APP_PORT)/city/$(id)/construction/build" anchor:='{"q": 0, "r": 0}' artifact_name=MeteoCenter height:=5

.PHONY: build-pathway
build-pathway:
	@curlie -v PATCH ":$(APP_PORT)/city/$(id)/construction/build-pathway" coordinates:='[{"q": 0, "r": 0}]' artifact_name=Road

# todo
.PHONY: get-tiles-spatial-values
.PHONY: get-city-stats
.PHONY: get-flagship
.PHONY: get-kind-budget
