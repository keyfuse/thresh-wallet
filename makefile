export GOPATH := $(shell pwd)

build:
	@echo "--> Building..."
	@mkdir -p bin/
	go build -v -o bin/threshwallet-server src/cmd/server.go
	go build -v -o bin/threshwallet-client src/cmd/client.go
	@chmod 755 bin/*

buildosx:
	@echo "--> Building osx library..."
	go get -v golang.org/x/mobile/cmd/gobind
	go get -v golang.org/x/mobile/cmd/gomobile
	./bin/gomobile bind -target=ios library

buildandroid:
	@echo "--> Building android library..."
	go get -v golang.org/x/mobile/cmd/gobind
	go get -v golang.org/x/mobile/cmd/gomobile
	./bin/gomobile bind -target=android library

clean:
	@echo "--> Cleaning..."
	@go clean
	@rm -f bin/*

test:
	@echo "--> Testing..."
	@$(MAKE) testxlog
	@$(MAKE) testproto
	@$(MAKE) testserver
	@$(MAKE) testlibrary

testxlog:
	go test -v xlog
testproto:
	go test -v proto
testlibrary:
	go test -v library
testserver:
	go test -v server


# code coverage
allpkgs =	xlog proto library server
coverage:
	go build -v -o bin/gotestcover \
	src/vendor/github.com/pierrre/gotestcover/*.go;
	bin/gotestcover -coverprofile=coverage.out -v $(allpkgs)
	go tool cover -html=coverage.out

check:
	go get -v github.com/golangci/golangci-lint/cmd/golangci-lint
	bin/golangci-lint run -D errcheck src/proto/... src/library/... src/server/... ../src/client/...

.PHONY: build clean install fmt test coverage
