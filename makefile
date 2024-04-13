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

bin/vpdb: clean
	@mkdir -p bin
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -o bin/vpdb \
		-ldflags="-X 'main.Version=$(GIT_TAG)' -X 'main.BuildDate=$(BUILD_DATE)' -X 'main.Commit=$(GIT_COMMIT)'" \
		cmd/main.go
