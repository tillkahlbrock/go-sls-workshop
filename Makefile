.PHONY: build clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/create create/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get-url get-url/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/preview preview/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
