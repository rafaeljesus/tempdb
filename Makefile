.PHONY: all deps test

all: deps test

deps:
	@go get -u github.com/golang/dep/cmd/dep
	@dep ensure

test:
	@go vet && go test -v -race -cover
