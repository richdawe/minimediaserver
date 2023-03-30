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
WATCHREGEX = '(\.go|\.tmpl\.html|\.js|\.css|\.mod|\.sum)$$'

watch:
	reflex -r $(WATCHREGEX) -s -- make run

watchtest:
	reflex -r $(WATCHREGEX) -s -- make test

clean:
	rm -fv bin/minimediaserver