BINARY_PATH=bin
BINARY_NAME=alien-invasion

build:
	mkdir -p ${BINARY_PATH}
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_PATH}/${BINARY_NAME}-darwin cmd/main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_PATH}/${BINARY_NAME}-linux cmd/main.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_PATH}/${BINARY_NAME}-windows cmd/main.go

deps:
	go get github.com/creamdog/gonfig

test:
	go test -v github.com/pastet89/alien-invasion/city
	go test -v github.com/pastet89/alien-invasion/world
	go test -v github.com/pastet89/alien-invasion/simulator
