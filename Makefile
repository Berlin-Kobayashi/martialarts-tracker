.PHONY: test
default: test ;

test:
	@docker run --rm -v ${PWD}/../../../:/go/src/ -w /go/src/github.com/DanShu93/martialarts-tracker golang:1.10 /bin/bash -c "go get ./... && go test ./..."
