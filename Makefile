


deps-update:
	go mod tidy

run:
	go run cmd/app/main/main.go -config=config.yaml

test:
	go clean -testcache
	go test ./pkg/...