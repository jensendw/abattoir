dependencies:
	go get github.com/op/go-logging
	go get github.com/antonholmquist/jason
	go get github.com/kelseyhightower/envconfig

build:
	go build

testrun:
	go run config.go connections.go main.go

all: dependencies build
