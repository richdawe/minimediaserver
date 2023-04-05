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

.PHONY:	lint
lint:	lint-go lint-js

.PHONY: lint-go
lint-go:
	sudo docker run -t --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v1.52.2 golangci-lint run -v

.PHONY:	lint-js
lint-js:
	npx eslint cmd/static
	@echo

# TODO: Add -race -covermode=atomic later?
# E.g.: when fetching tracks from multiple clients?
.PHONE:	coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Requires 'reflex' from https://github.com/cespare/reflex
# to be in the path.
WATCHREGEX = '(\.go|\.tmpl\.html|\.js|\.css|\.mod|\.sum)$$'

watch:
	reflex -r $(WATCHREGEX) -s -- make run

watchtest:
	reflex -r $(WATCHREGEX) -s -- make test

watchcoverage:
	reflex -r $(WATCHREGEX) -s -- make coverage

clean:
	rm -fv bin/minimediaserver
