GO_FILES = $(shell find . -type f -name '*.go')

.PHONY: default
default:
	@# do nothing

bin/ssh-compose: $(GO_FILES)
	@mkdir -p bin
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/ssh-compose cmd/main.go
