VERSION=$(shell git describe --tags)
LDFLAGS=-ldflags "-s -w"

all: linux darwin windows

release: all zip

clean:
	rm -rf bin/* *.zip

upx:
	upx -9 bin/*

linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64 ${LDFLAGS} main.go
	CGO_ENABLED=1 GOOS=linux GOARCH=386 go build -o bin/linux-i386 ${LDFLAGS} main.go

darwin:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o bin/darwin-amd64 ${LDFLAGS} main.go
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o bin/darwin-arm64 ${LDFLAGS} main.go

windows:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o bin/windows-amd64.exe ${LDFLAGS} main.go
	CGO_ENABLED=1 GOOS=windows GOARCH=386 go build -o bin/windows-i386.exe ${LDFLAGS} main.go

build:
	@go build -o tmp/roja-shop

dev:
	@go run main.go
run: build
	@./bin/roja-shop

test:
	@go test -v ./...

migrate:
 	~/go/bin/migrate -database "sqlite://$PWD/storage/database/roja.db" -path "$PWD/migrations" up 
	
watch:
	@~/go/bin/air air.conf
	# @docker run -it --rm \
	# 	-w "/go/src/github.com/cosmtrek/hub" \
	# 	-v .:/go/src/github.com/cosmtrek/hub \
	# 	-p 3000:3000 \
    # 	cosmtrek/air

.PHONY: watch build run test