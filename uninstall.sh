#!/usr/bin/env sh
echo 'Removing directories'
rm -rf ~/.vhost

echo 'Removing files'
rm /usr/bin/vhost
rm /usr/lib/python2.6/vhost_*.py

echo 'Done'