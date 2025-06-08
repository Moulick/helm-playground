.PHONY: browser-test build test serve

test:
	go test go/lib/lib_test.go go/lib/lib.go go/lib/helm_engine.go

build:
	mkdir -p dist
	rm -f dist/*
	cp "$(shell go env GOROOT)/lib/wasm/wasm_exec.js" dist/wasm_exec.js
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o ./dist/lib.wasm go/main/main.go

browser-test: build
	yarn --cwd ./browser-test
	node browser-test/run.js

serve: build
	npx http-server -c-1
