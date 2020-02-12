default: vet test

vet:
	go vet ./...

test:
	go test ./...

README.md: README.md.tpl $(wildcard *.go)
	becca -package github.com/bsm/conditional
