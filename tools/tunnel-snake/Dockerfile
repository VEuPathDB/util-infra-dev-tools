FROM golang:1.22.2-alpine3.19 AS build

WORKDIR /tunnel
COPY . .
RUN mkdir build && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/tunnel-snake cmd/main.go

FROM alpine:3.19.1

COPY --from=build /tunnel/build/tunnel-snake /usr/bin/tunnel-snake
CMD tunnel-snake
