.PHONY: test
test:
	go test -race -v -coverprofile=cover.out ./...
	go tool cover -func=cover.out