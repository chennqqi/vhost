.ONESHELL:

build:
	go build

install: build
install:
	sudo cp vhost /usr/local/bin/vhost
	mkdir -p ~/.vhost
	cp -r shared/* ~/.vhost
	cp config.yaml ~/.vhost/config.yaml

arch-pkg:
	makepkg .


.PHONY: build install arch-pkg