_default:
    @just --choose

run:
    go run .

build:
    go build .

tidy:
    go mod tidy
