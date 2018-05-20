GITHUB_TOKEN=$(shell grep -oE "GITHUB_TOKEN=.+" $(PWD)/env/github_vars.sh | grep -oE "=.+" | grep -Eo "[^=].+")

build:
	go build -v
	go test -v
	go vet

run: build
	@GITHUB_TOKEN=$(GITHUB_TOKEN) go run main.go
