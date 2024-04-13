GIT_TAG = $(shell git describe --tags --abbrev=0 2>/dev/null || echo "alpha")
GIT_COMMIT = $(shell git rev-parse HEAD)
BUILD_DATE = $(shell date --rfc-3339=seconds)


.PHONY: default
default:
	@# do nothing

.PHONY: build
build: bin/vpdb

.PHONY: clean
clean:
	@rm -rf bin

bin/linux/vpdb: clean
	@mkdir -p bin/linux
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -o $@ \
		-ldflags="-X 'main.Version=$(GIT_TAG)' -X 'main.BuildDate=$(BUILD_DATE)' -X 'main.Commit=$(GIT_COMMIT)'" \
		cmd/main.go

bin/darwin/vpdb: clean
	@mkdir -p bin/darwin
	@env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 \
		go build -o $@ \
		-ldflags="-X 'main.Version=$(GIT_TAG)' -X 'main.BuildDate=$(BUILD_DATE)' -X 'main.Commit=$(GIT_COMMIT)'" \
		cmd/main.go

bin/vpdb-linux.tar.gz: bin/linux/vpdb
	@rm bin/$@
	@cd bin/linux && tar -czf $(@F) vpdb && cp $(@F) ..

bin/vpdb-darwin.tar.gz: bin/darwin/vpdb