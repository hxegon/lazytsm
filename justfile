_default:
    @just --choose

run: build
    ./bin/lazytsm

build:
    go build -o bin/lazytsm

tidy:
    go mod tidy
