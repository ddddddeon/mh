.PHONY: build

build:
	go build -o bin/mh .

install: 
	cp bin/mh /usr/bin/mh
