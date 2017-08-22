.PHONY: build deps test run bundle

build: test
	cd daemon; go build -v

deps:
	dep ensure

test: deps
	go test -v

run: build
	cd daemon; sudo ./daemon

bundle: build
	mv daemon/daemon volmex-daemon; tar cvf volmex.tar volmex-daemon volmex.service README.md

