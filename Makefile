.PHONY: build deps test run bundle

build:
	cd daemon; go build -v

deps:
	go get -v -t

test:
	go test -v

run: deps build
	cd daemon; sudo ./daemon

bundle: deps build
	mv daemon/daemon volmex-daemon; tar cvf volmex.tar volmex-daemon volmex.service README.md

