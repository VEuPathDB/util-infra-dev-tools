GIT_TAG    := $(shell git describe --tags 2>/dev/null || echo "alpha")
GIT_COMMIT := $(shell git rev-parse HEAD)
BUILD_DATE := $(shell date --rfc-3339=seconds)
CURRENT_OS := $(shell uname | tr '[:upper:]' '[:lower:]')
OS_TARGETS := linux darwin windows

CWD := $(shell pwd)

GO_FILES := $(shell find . -type f -name '*.go')

LD_FLAGS := -X 'main.Version=$(GIT_TAG)' -X 'main.BuildDate=$(BUILD_DATE)' -X 'main.Commit=$(GIT_COMMIT)'

.PHONY: default
default:
	@# do nothing

.PHONY: build
build: bin/vpdb bin/merge-compose

.PHONY: install
install: bin/vpdb
	@cp $< ${HOME}/.local/bin/vpdb
	@mkdir -p "${HOME}/.local/share/vpdb"
	@cp scripts/autocomplete.sh "${HOME}/.local/share/vpdb"

.PHONY: clean
clean:
	@rm -rf bin

.PHONY: release
release: clean $(foreach os,$(OS_TARGETS),bin/vpdb-$(os).zip) #$(foreach os,$(OS_TARGETS),bin/merge-compose-$(os).zip)

.PHONY: publish-local
publish-local: bin/$(CURRENT_OS)/vpdb
	@cp $< $${HOME}/.local/bin/vpdb

## VPDB BUILD

bin/vpdb: bin/$(CURRENT_OS)/vpdb
	@cp $< $@

.SECONDARY:
bin/%/vpdb: $(GO_FILES)
	@mkdir -p bin/$*
	@env CGO_ENABLED=0 GOOS=$* GOARCH=amd64 go build -o $@ -ldflags="$(LD_FLAGS)" cmd/vpdb/main.go

bin/vpdb-%.zip: bin/%/vpdb
	@rm -f $@
	@cd bin/$* \
	  && if [ "$*" = "windows" ]; then \
	    mv vpdb vpdb.exe \
	    && zip -q $(@F) vpdb.exe; \
	  else \
	    cp $(CWD)/scripts/autocomplete.sh .; \
	    zip -q $(@F) vpdb autocomplete.sh; \
	  fi \
	  && mv $(@F) ..
	@rm -rf bin/$*

## MERGE COMPOSE BUILD

bin/merge-compose: bin/$(CURRENT_OS)/merge-compose
	@cp $< $@

.SECONDARY:
bin/%/merge-compose: $(GO_FILES)
	@mkdir -p bin/$*
	@env CGO_ENABLED=0 GOOS=$* GOARCH=amd64 go build -o $@ -ldflags="$(LD_FLAGS)" cmd/merge-compose/main.go

bin/merge-compose-%.zip: bin/%/merge-compose
	@rm -f $@
	@cd bin/$* \
		&& if [ "$*" = "windows" ]; then mv merge-compose merge-compose.exe && zip -q $(@F) merge-compose.exe; else zip -q $(@F) merge-compose; fi \
		&& mv $(@F) ..
	@rm -rf bin/$*
