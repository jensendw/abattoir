dependencies:
	go get github.com/op/go-logging
	go get github.com/antonholmquist/jason
	go get github.com/kelseyhightower/envconfig

build:
	go build

all: dependencies build
