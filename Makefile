## hocon Makefile:
#:

SHELL          = /bin/sh
GO            ?= go
CFG           ?= .env
PRG           ?= $(shell basename $$PWD)

SOURCES        = $(shell find . -maxdepth 3 -mindepth 1 -path ./var -prune -o -name '*.go')

# do not include docker-compose.yml from dcape
# docker-compose can't build docker image if it included
DCAPE_DC_USED  = false

VERSION       ?= $(shell git describe --tags --always)
# Last project tag
RELEASE       ?= $(shell git describe --tags --abbrev=0 --always)

APP_ROOT      ?= .
APP_SITE      ?= $(PRG).dev.lan
APP_PROTO     ?= http

APP_TOKEN     ?= $(shell < /dev/urandom tr -dc A-Za-z0-9 | head -c16; echo)

IMAGE         ?= $(PRG)
IMAGE_VER     ?= latest

# Hardcoded in docker-compose.yml service name
DC_SERVICE    ?= app

# docker app for change inside containers
DOCKER_BIN    ?= docker

# Docker-compose project name (container name prefix)
PROJECT_NAME  ?= $(PRG)

# dcape network connect to, must be set in .env
DCAPE_NET     ?= dcape_default

GOGENS_IMG    ?= ghcr.io/apisite/gogens:latest

# ------------------------------------------------------------------------------
# .env template (custom part)
# inserted in .env.sample via 'make config'
define CONFIG_CUSTOM
# ------------------------------------------------------------------------------
# app custom config, generated by make config

APP_PROTO=$(APP_PROTO)

APP_TOKEN=$(APP_TOKEN)

# dcape network connect to, must be set in .env
DCAPE_NET=$(DCAPE_NET)

endef

.PHONY: buf.lock dep build run lint test clean



# ------------------------------------------------------------------------------
## Code generate operations
#:

buf.lock:
	@id=$$(docker create $(GOGENS_IMG)) ; \
	docker cp $$id:/app/$@ $@ ; \
	docker rm -v $$id

gen-go:
	docker run --rm  -v `pwd`:/mnt/pwd -w /mnt/pwd $(GOGENS_IMG) --debug generate --template buf.gen.yaml --path proto

buf:
	docker run --rm -it  -v `pwd`:/mnt/pwd -w /mnt/pwd $(GOGENS_IMG) $(CMD)

gen-js: static/js/api.js

static/js/api.js: zgen/ts/proto/service.pb.ts
	docker run --rm  -v `pwd`:/mnt/pwd -w /mnt/pwd --entrypoint /go/bin/esbuild $(GOGENS_IMG)  \
	zgen/ts/proto/service.pb.ts --bundle --outfile=/mnt/pwd/static/js/api.js --global-name=AppAPI

#	--sourcemap --target=chrome58 \
#	--minify --sourcemap --target=chrome58,firefox57,safari11,edge16 \

gen: gen-go gen-js

# ------------------------------------------------------------------------------
## Compile operations
#:

$(PRG): $(SOURCES)
	$(GO) build -ldflags "-X main.version=$(VERSION)" ./cmd/$(PRG)

run: $(PRG)
	./$(PRG) --token $(APP_TOKEN)

## Format go sources
fmt:
	$(GO) fmt ./...

## Run vet
vet:
	$(GO) vet ./...

## Run linter
lint:
	golint ./...

## Run more linters
lint-more:
	golangci-lint run ./...

## Run tests
test: coverage.out

## Run tests and fill coverage.out
cov: coverage.out

# internal target
coverage.out: $(SOURCES)
	GIN_MODE=release $(GO) test -test.v -test.race -coverprofile=$@ -covermode=atomic ./...

## Open coverage report in browser
cov-html: cov
	$(GO) tool cover -html=coverage.out

## Clean coverage report
cov-clean:
	rm -f coverage.*

## Changes from last tag
changelog:
	@echo Changes since $(RELEASE)
	@echo
	@git log $(RELEASE)..@ --pretty=format:"* %s"

# ------------------------------------------------------------------------------
## Docker operations
#:

docker: $(PRG)
	docker build -t $(PRG) .

# ------------------------------------------------------------------------------


## Build docker image
docker-build: CMD="build --no-cache --force-rm $(DC_SERVICE)"
docker-build: dc

## Remove docker image & temp files
docker-clean:
	[ "$$($(DOCKER_BIN) images -q $(DC_IMAGE) 2> /dev/null)" = "" ] || $(DOCKER_BIN) rmi $(DC_IMAGE)

clean: ## Remove previous builds
	@rm -f $(PRG)

# ------------------------------------------------------------------------------
# Find and include DCAPE/apps/drone/dcape-app/Makefile
DCAPE_COMPOSE   ?= dcape-compose
DCAPE_MAKEFILE  ?= $(shell docker inspect -f "{{.Config.Labels.dcape_app_makefile}}" $(DCAPE_COMPOSE))
ifeq ($(shell test -e $(DCAPE_MAKEFILE) && echo -n yes),yes)
  include $(DCAPE_MAKEFILE)
else
  include /opt/dcape-app/Makefile
endif
