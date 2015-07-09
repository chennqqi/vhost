#!/usr/bin/env bash

if [[ ! -d ~/.vhost ]] ; then
    mkdir -p ~/.vhost
fi

cp -vR share ~/.vhost/share
cp -v vhost.conf ~/.vhost/vhost.conf
sudo cp -v vhost.py /usr/bin/vhost

xdg-open ~/.vhost/vhost.conf