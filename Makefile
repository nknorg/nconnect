.DEFAULT_GOAL:=local_or_with_proxy

USE_PROXY=GOPROXY=https://goproxy.io
VERSION:=$(shell git describe --abbrev=7 --dirty --always --tags)
BUILD=go build -ldflags "-s -w -X main.Version=$(VERSION) -X main.GOARM=$(GOARM)"
BUILD_DIR=build
BIN_NAME=nConnect
ifdef GOARM
BIN_DIR=$(GOOS)-$(GOARCH)v$(GOARM)
else
BIN_DIR=$(GOOS)-$(GOARCH)
endif

web/dist: $(shell find web/src -type f -not -path "web/src/node_modules/*" -not -path "web/src/build/*")
	@rm -rf web/dist
	-@cd web/src && yarn && yarn build && cp -a ./build ../dist

.PHONY: local
local: web/dist
	$(BUILD) -o $(BIN_NAME)$(EXT) .

.PHONY: local_with_proxy
local_with_proxy: web/dist
	$(USE_PROXY) $(BUILD) -o $(BIN_NAME)$(EXT) .

.PHONY: local_or_with_proxy
local_or_with_proxy:
	${MAKE} local || ${MAKE} local_with_proxy

.PHONY: build
build: web/dist
	rm -rf $(BUILD_DIR)/$(BIN_DIR)
	mkdir -p $(BUILD_DIR)/$(BIN_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) $(BUILD) -o $(BUILD_DIR)/$(BIN_DIR)/$(BIN_NAME)$(EXT) .
	mkdir -p $(BUILD_DIR)/$(BIN_DIR)/web/
	@cp -a web/dist $(BUILD_DIR)/$(BIN_DIR)/web/
	${MAKE} tar

.PHONY: tar
tar:
	cd $(BUILD_DIR) && rm -f $(BIN_DIR).tar.gz && tar --exclude ".DS_Store" --exclude "__MACOSX" -czvf $(BIN_DIR).tar.gz $(BIN_DIR)

.PHONY: zip
zip:
	cd $(BUILD_DIR) && rm -f $(BIN_DIR).zip && zip --exclude "*.DS_Store*" --exclude "*__MACOSX*" -r $(BIN_DIR).zip $(BIN_DIR)

.PHONY: all
all:
	${MAKE} build GOOS=linux GOARCH=amd64
	# ${MAKE} build GOOS=linux GOARCH=arm GOARM=5
	# ${MAKE} build GOOS=linux GOARCH=arm GOARM=6
	# ${MAKE} build GOOS=linux GOARCH=arm GOARM=7
	# ${MAKE} build GOOS=linux GOARCH=arm64
	# ${MAKE} build GOOS=darwin GOARCH=amd64
	# ${MAKE} build GOOS=windows GOARCH=amd64 EXT=.exe
