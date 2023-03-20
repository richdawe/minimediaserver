default:	build

build:
	go build -v -o bin/minimediaserver ./cmd

run:
	go run ./cmd

# Requires 'reflex' from https://github.com/cespare/reflex
# to be in the path.
watch:
	reflex -r '\.go$$' -s -- go run ./cmd

clean:
	rm -fv bin/minimediaserver