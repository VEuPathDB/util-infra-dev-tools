.PHONY: default
default:
	@# Do nothing

.PHONY: docker-build
docker-build:
	@docker build -t veupathdb/tunnel-snake:latest .

.PHONY: go-build
go-build: clean
	@mkdir build
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/tunnel-snake cmd/main.go

.PHONY: clean
clean:
	@rm -rf build