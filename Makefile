build:
	@go build -trimpath -o webserver github.com/pieterclaerhout/go-webserver/cmd/webserver

run: build
	@DEBUG=1 PORT=8080 ./webserver
