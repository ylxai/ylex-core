# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: ylexd android ios ylexd-cross cvm all test clean
.PHONY: ylexd-linux ylexd-linux-386 ylexd-linux-amd64 ylexd-linux-mips64 ylexd-linux-mips64le
.PHONY: ylexd-linux-arm ylexd-linux-arm-5 ylexd-linux-arm-6 ylexd-linux-arm-7 ylexd-linux-arm64
.PHONY: ylexd-darwin ylexd-darwin-386 ylexd-darwin-amd64
.PHONY: ylexd-windows ylexd-windows-386 ylexd-windows-amd64

GOBIN = ./build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run

ylexd:
	$(GORUN) build/ci.go install ./cmd/ylexd
	@echo "Done building."
	@echo "Run \"$(GOBIN)/ylexd\" to launch ylexd."

all:
	$(GORUN) build/ci.go install

android:
	$(GORUN) build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/ylex.aar\" to use the library."

ios:
	$(GORUN) build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Ylex.framework\" to use the library."

test: all
	$(GORUN) build/ci.go test

lint: ## Run linters.
	$(GORUN) build/ci.go lint

clean:
	env GO111MODULE=on go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "ylem" 2> /dev/null || echo 'Please install ylem'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'
