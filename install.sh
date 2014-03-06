#!/usr/bin/env bash
python -c 'import sys; print(sys.version[:3])' > ver

if [ -f "/usr/bin/vhost" ];
then
echo 'Reinstallation.'
else
echo 'Installation.'
fi

echo 'Installing directories'
mkdir -p ~/.vhost/var
mkdir -p ~/.vhost/share

echo 'Installing files'
cp -R --remove-destination etc/config.ini ~/.vhost/
cp -R --remove-destination share/* ~/.vhost/share
cp -R --remove-destination vhost.py /usr/bin/vhost
cp -R --remove-destination include/*.py /usr/lib/python`cat ver`

echo 'Done. Now run vhost -h to get help'
