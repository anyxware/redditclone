.SILENT:
.PHONY:

build:
	go build -o .bin/main cmd/redditclone/main.go

run: build
	.bin/main

clear:
	rm -rf .bin

test_handler:
	go test -v ./internal/handler -coverprofile=internal/handler/cover.out \
		&& go tool cover -html=internal/handler/cover.out -o internal/handler/coverage.html \
		&& rm -f internal/handler/cover.out

test_mysqlrepo:
	go test -v ./internal/repository/mysqlrepo -coverprofile=internal/repository/mysqlrepo/cover.out \
		&& go tool cover -html=internal/repository/mysqlrepo/cover.out -o internal/repository/mysqlrepo/coverage.html \
		&& rm -f internal/repository/mysqlrepo/cover.out

test_mongorepo:
	go test -v ./internal/repository/mongorepo -coverprofile=internal/repository/mongorepo/cover.out \
		&& go tool cover -html=internal/repository/mongorepo/cover.out -o internal/repository/mongorepo/coverage.html \
		&& rm -f internal/repository/mongorepo/cover.out

migration:
	go run cmd/migration/main.go