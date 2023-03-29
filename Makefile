build_debug:
	cd src/gabo && make build_debug

build: build_debug

build_prod:
	cd src/gabo && make build_prod

# Ref: https://goreleaser.com/install/
build_release:
	goreleaser check && goreleaser build --snapshot --clean

go_lint:
	cd src/gabo && make lint

lint: format go_lint

format:
	cd src/gabo && make format

clean:
	cd src/gabo && make clean

test:
	cd src/gabo && make test
