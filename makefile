lint:
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify
install: go.sum 
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/sad
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/sacli