_default:
    @just --choose

run:
    go run cmd/lazytsm.go

build:
    go build cmd/lazytsm.go

tidy:
    go mod tidy
