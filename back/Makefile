-include .env

export $(shell sed 's/=.*//' .env)

destination ?= last

define SELECT_APP
	printf "use debug-main ? (y/n): "; \
	read ans; \
	if [ "$$ans" = "y" ] || [ "$$ans" = "Y" ]; then \
		APP="debug"; \
	else \
		APP="server"; \
	fi
endef

.PHONY: run
run:
	echo "generating files..."; \
	go generate -tags=generator ./... && \
	go run ./cmd/server

# run with wgo for live-reload
.PHONY: wgo-run
wgo-run:
	@$(SELECT_APP); \
	echo "generating files..."; \
	go generate -tags=generator ./... && \
	wgo run ./cmd/$$APP

# generate db queries functions
.PHONY: sqlc
sqlc:
	cd ./internal/adapters/postgres && sqlc generate

# run migrations
# add flags="xxx" to command
# ex: make migrate destination="0" => to drop everything
# ex: make migrate => to migrate to latest version
.PHONY: migrate
migrate:
	tern migrate \
	--config internal/adapters/postgres/.migrations/tern.conf \
	--migrations internal/adapters/postgres/.migrations \
	--destination $(destination)

.PHONY: build
build:
	@echo "generating files..."
	go generate -tags=generator ./...
	@echo "compiling..."
	go build -ldflags="-s -w" -o ./bin/app_prod ./cmd/server

.PHONY: start
start:
	./bin/app_prod
