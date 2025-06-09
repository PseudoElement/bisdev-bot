start-local:
	go run .

start-prod:
	./builds/bot-ubuntu

build-macos:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o ./builds/bot-macos

build-ubuntu:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -o ./builds/bot-ubuntu

build-win64:
	GOOS=windows GOARCH=amd64 go build -o builds/bot-win64.exe