#!/usr/bin/env python
import argparse, getpass, sys, logging, configparser, os, subprocess

logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)

configs = (os.getenv('HOME') + '/.vhost/vhost.conf', '/etc/vhost.conf', 'vhost.conf')

config_files = [config for config in configs if os.path.exists(config)]
if len(config_files) == 0:
    logger.error('Vhost is not configured. Configure ~/.vhost.conf first.')
    sys.exit(1)

config = configparser.ConfigParser()
config.read(configs)

if getpass.getuser() != 'root':
    print ('* You must be root to use this program.')
    sys.exit()

def restart_httpd():
    subprocess.call(config.get('commands', 'httpd_reload').split())

def exists(path):
    return os.path.exists(path)

def get_sitename(name):
    return name + config.get('general', 'domain', fallback='.local')

def get_vhost_avail_path(vhost_name):
    ext = config.get('general', 'vhost_file_suffix', fallback='.conf')
    path = config.get('apache', 'dir_hosts_available', fallback='/etc/vhost/sites-available') + '/' + vhost_name + ext
    return path

def get_vhost_enabl_path(vhost_name):
    ext = config.get('general', 'vhost_file_suffix', fallback='.conf')
    path = config.get('apache', 'dir_hosts_enabled', fallback='/etc/vhost/sites-enabled') + '/' + vhost_name + ext
    return path

def find_file(file_set):
    return [file for file in file_set if os.path.exists(file)]

def has_in_hosts(sitename):
    handle = open(config.get('general', 'hosts_file', fallback='/etc/hosts'), 'r')
    contents = handle.read()
    handle.close()
    return sitename in contents

def add_to_hosts(sitename):
    path = config.get('general', 'hosts_file', fallback='/etc/hosts')
    if not has_in_hosts(sitename):
        logger.debug('--> add vhost to %s' % path)
        handle = open(path, 'a')
        handle.write('127.0.0.1         %s\n' % sitename)
        handle.close()
    else:
        logger.debug('--> sitename already in %s' % path)

def remove_from_hosts(sitename):
    path =  config.get('general', 'hosts_file', fallback='/etc/hosts')
    if has_in_hosts(sitename):
        handle = open(path, 'r+') # rw
        contents = handle.readlines()
        new_contents = ''
        for line in contents:
            if sitename not in line:
                new_contents += line

        handle.seek(0)
        handle.write(new_contents)
        handle.truncate()
        handle.close()
        logger.debug('--> removed from %s' % path)
    else:
        logger.debug('--> sitename is not in %s' % path)

def is_enabled(name):
    path = get_vhost_enabl_path(name)
    return exists(path)

def _create(args):
    logger.info('Create: %s', args.name)

    if os.path.exists(get_vhost_avail_path(args.name)):
        logger.warning('--> vhost exists. you may want to enable this?')
        sys.exit()

    sitesroot = config.get('general', 'sites_dir', fallback='/var/www')
    template = find_file((os.getenv('HOME') + '/.vhost/share/vhost.conf', '/etc/vhost/share/vhost.conf', 'share/vhost.conf'))[0]
    contents = open(template).read()
    contents = contents\
        .replace('%name%', get_sitename(args.name))\
        .replace('%sitesdir%', sitesroot)\
        .replace('%ip%', args.ip)

    # handle --subdir switch
    subdir = args.subdir or ''
    if len(subdir) > 0 and subdir[0] != '/':
        subdir = '/' + subdir
    contents = contents.replace('%subdir%', subdir)

    if args.ssl:
        ssl_part = find_file((os.getenv('HOME') + '/.vhost/share/ssl.conf', '/etc/vhost/ssl.conf', 'share/ssl.conf'))[0]
        if not ssl_part:
            raise Exception('Could not find SSL part.')

        cert_file = config.get('ssl', 'cert_file', fallback=None)
        key_file = config.get('ssl', 'key_file', fallback=None)

        if cert_file is None or key_file is None:
            raise Exception('Either ssl.cert_file or ssl.key_file is not found in config.')

        if not os.path.exists(cert_file):
            raise Exception('File not found: %s' % cert_file)

        if not os.path.exists(key_file):
            raise Exception('File not found: %s' % key_file)

        ssl_content = open(ssl_part).read()
        ssl_content = ssl_content.replace('%cert_file%', cert_file).replace('%key_file%', key_file)

        contents = contents.replace('%ssl%', ssl_content)
    else:
        contents = contents.replace('%ssl%', '')

    # open file and write contents
    avail_vhost_path = get_vhost_avail_path(args.name)

    output = open(avail_vhost_path, 'w')
    output.write(contents)
    output.close()

    dirs = ('log', 'www' + subdir, 'tmp')
    for dir in dirs:
        new_dir = '%s/%s/%s' % (sitesroot, get_sitename(args.name), dir)
        logger.info('Create directory: %s' % new_dir)
        if os.path.exists(new_dir):
            logger.warning('--> already exists')
        else:
            os.makedirs(new_dir)

    if args.sample:
        sample_index = '%s/%s/www/index.html' % (sitesroot, get_sitename(args.name))
        index_template = find_file((os.getenv('HOME') + '/.vhost/share/index.html', '/etc/vhost/share/index.html', 'share/index.html'))[0]
        index_contents = open(index_template).read()
        index_contents = index_contents.replace('%site%', get_sitename(args.name))
        out = open(sample_index, 'w')
        out.write(index_contents)
        out.close()
        logger.info('--> add index.html: %s' % sample_index)

def _enable(args):
    logger.info('Enable: %s' % args.name)
    if is_enabled(args.name):
        logger.warning('--> already enabled')
    else:
        if exists(get_vhost_avail_path(args.name)):
            os.symlink(get_vhost_avail_path(args.name), get_vhost_enabl_path(args.name))
            add_to_hosts(get_sitename(args.name))
            restart_httpd()
        else:
            logger.error('--> vhost does not exists')

def _disable(args):
    logger.info('Disable: %s' % args.name)
    if exists(get_vhost_avail_path(args.name)):
        if not is_enabled(args.name):
            logger.warning('--> not enabled')
        else:
            os.remove(get_vhost_enabl_path(args.name))
            remove_from_hosts(get_sitename(args.name))
            restart_httpd()
    else:
        logger.error('--> vhost does not exists')

def main():
    parser = argparse.ArgumentParser(
        prog='vhost',
        description='Handy helper for easy PHP development.',
        epilog = 'Bug reports send to alex.oleshkevich@gmail.com'
    )

    group = parser.add_mutually_exclusive_group()

    # create
    group.add_argument('-c', '--create', help='create a new vhost', action='store_true', dest='create', default=False)
    parser.add_argument('--subdir', help='point document root this subdirectory', action='store', dest='subdir', default=None)
    parser.add_argument('--sample', help='add sample index.html file', action='store_true', dest='sample', default=None)
    parser.add_argument('--ip', help='bind to that IP address', action='store', dest='ip', default='*')
    parser.add_argument('--ssl', help='use SSL for that vhost', action='store_true', dest='ssl', default=False)

    # enable
    group.add_argument('-e', '--enable', help='enable existing vhost', action='store_true', dest='enable', default=False)

    # disable
    group.add_argument('-d', '--disable', help='disable vhost', action='store_true', dest='disable', default=False)

    # vhost name
    parser.add_argument('name', action='store', default=False, help='vhost name')

    args = parser.parse_args()

    try:
        if args.create:
            _create(args)
            _enable(args)
        elif args.enable:
            _enable(args)
        elif args.disable:
            _disable(args)
    except Exception as e:
        print(e)

if __name__ == '__main__':
    try:
        main()
    except OSError as e:
        logger.critical(e.strerror)
        sys.exit(e.errno)
