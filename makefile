GIT_TAG = $(shell git describe --tags --abbrev=0 2>/dev/null || echo "alpha")
GIT_COMMIT = $(shell git rev-parse HEAD)
BUILD_DATE = $(shell date --rfc-3339=seconds)
CURRENT_OS = $(shell uname | tr '[:upper:]' '[:lower:]')
OS_TARGETS = linux darwin windows

.PHONY: default
default:
	@# do nothing

.PHONY: build
build: bin/vpdb

.PHONY: clean
clean:
	@rm -rf bin

.PHONY: release
release: clean $(foreach os,$(OS_TARGETS),bin/vpdb-$(os).zip)

bin/vpdb: bin/$(CURRENT_OS)/vpdb
	@cp $< $@

.SECONDARY:
bin/%/vpdb:
	@mkdir -p bin/$*
	@env CGO_ENABLED=0 GOOS=$* GOARCH=amd64 \
		go build -o $@ \
		-ldflags="-X 'main.Version=$(GIT_TAG)' -X 'main.BuildDate=$(BUILD_DATE)' -X 'main.Commit=$(GIT_COMMIT)'" \
		cmd/main.go

bin/vpdb-%.zip: bin/%/vpdb
	@rm -f $@
	@cd bin/$* \
		&& if [ "$*" = "windows" ]; then mv vpdb vpdb.exe && zip -q $(@F) vpdb.exe; else zip -q $(@F) vpdb; fi \
		&& mv $(@F) ..
	@rm -rf bin/$*
