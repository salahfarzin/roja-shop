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