default:	build

.PHONY: build
build:
	go build -v -o bin/minimediaserver ./cmd

.PHONY: run
run:
	go run ./cmd
	@echo

.PHONY:	test
test:
	go test -v ./...
	@echo

# TODO: code coverage

# Requires 'reflex' from https://github.com/cespare/reflex
# to be in the path.
watch:
	reflex -r '(\.go|\.tmpl\.html)$$' -s -- make run

watchtest:
	reflex -r '(\.go|\.tmpl\.html)$$' -s -- make test

clean:
	rm -fv bin/minimediaserver