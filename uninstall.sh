#!/usr/bin/env sh
python -c 'import sys; print(sys.version[:3])' > ver
echo 'Removing directories'
rm -rf ~/.vhost

echo 'Removing files'
rm /usr/bin/vhost
rm /usr/lib/python`cat ver`/vhost_*.py

rm ver
echo 'Done'
