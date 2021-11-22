.PHONY: test,build,install

INSTALLPATH=/usr/local/bin/gonanny

test:
	go test ./...

build:
	go build -o gonanny main.go

install: build
	sudo install gonanny ${INSTALLPATH}

uninstall:
	sudo rm ${INSTALLPATH}
