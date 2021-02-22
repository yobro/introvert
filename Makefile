GO ?= go
GOFMT ?= $(GO)fmt
GOOPTS ?= 
GO111MODULE := on 

WEB_APP_PATH = web/app
WEB_APP_SOURCE_FILES = $(shell find $(WEB_APP_PATH)/src/ $(WEB_APP_PATH)/tsconfig.json)
WEB_APP_OUTPUT_DIR = web/static/app/dist
WEB_APP_NODE_MODULES_PATH = $(WEB_APP_PATH)/node_modules
WEB_APP_BUILD_SCRIPT = ./scripts/build_web_app.sh

$(WEB_APP_OUTPUT_DIR): $(WEB_APP_NODE_MODULES_PATH) $(WEB_APP_SOURCE_FILES) $(WEB_APP_BUILD_SCRIPT)
	@echo ">> building React app"
	@$(WEB_APP_BUILD_SCRIPT)


assets: $(WEB_APP_OUTPUT_DIR)
	@echo ">> writing assets"
	cd web && GO111MODULE=$(GO111MODULE) GOOS= GOARCH= $(GO) generate -x -v && cd -
	@$(GOFMT) -w ./web

build: assets
	@echo ">> building binary"
	# $(GO) build -tags builtinassets ./cmd/introvert/main.go
	$(GO) build ./cmd/introvert/main.go

clean:
	go clean
	rm -f web/assets_vfsdata.go main