build_debug:
	GO111MODULE=on go build -v -o bin/${BINARY_NAME} ./cmd/gabo

build_prod:
	# Shrink binary by removing symbol and DWARF table
	# Ref: https://lukeeckley.com/post/useful-go-build-flags/
	GO111MODULE=on go build -v -ldflags="-s -w" -o bin/${BINARY_NAME} ./cmd/gabo

go_lint:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go vet ./...
	# Deprecated
	# golint -set_exit_status ./...
	go tool fix ./..
	golangci-lint run

go_format:
	go fmt ./...

format: go_format

lint: go_lint

clean:
	GO111MODULE=on go clean --modcache
	rm -rf bin/*

test:
	GO111MODULE=on go test ./... -v

update_go_deps:
	GO111MODULE=on go get -t -u ./...

