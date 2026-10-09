.PHONY: build run test fmt

APP := pingulina
CMD := ./cmd/app
LDFLAGS := -s -w

build:
	CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o $(APP) $(CMD)

run:
	go run $(CMD)

test:
	go test ./...

fmt:
	gofmt -w .
	goimports -w .
