.PHONY: build run test fmt

APP := pingulina
CMD := .
LDFLAGS := -s -w
ARGS ?=

build:
	CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o $(APP) $(CMD)

run:
	go run $(CMD) $(ARGS)

test:
	go test ./...

fmt:
	gofmt -w .
	goimports -w .
