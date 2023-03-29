build_debug:
	cd src/gabo && make build_debug

build: build_debug

build_prod:
	cd src/gabo && make build_prod

go_lint:
	cd src/gabo && make lint

lint: format go_lint

format:
	cd src/gabo && make format

clean:
	cd src/gabo && make clean

test:
	cd src/gabo && make test
