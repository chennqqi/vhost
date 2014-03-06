#!/usr/bin/env python
# Copyright (C) 2011 Alex Oleshkevich <alex.oleshkevich@gmail.com>
#
# Authors:
#  Alex Oleshkevich
#
# This program is free software; you can redistribute it and/or modify it under
# the terms of the GNU General Public License as published by the Free Software
# Foundation; version 3.
#
# This program is distributed in the hope that it will be useful, but WITHOUT
# ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
# FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more
# details.
#
# You should have received a copy of the GNU General Public License along with
# this program; if not, write to the Free Software Foundation, Inc.,
# 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA

import os

class vhost_class():
    """
        Constructor
    """
    def __init__(self, name, config):
        self.ini = config
        self.hosts_dir = self.ini.get_option('apache', 'hosts_dir')
        self.enabled_dir = self.ini.get_option('apache', 'enabled_hosts')
        self.hosts_file = self.ini.get_option('apache', 'hosts_file')
        self.doc_roots = self.ini.get_option('apache', 'doc_roots')
        self.postfix = self.ini.get_option('general', 'postfix')
        self.name = name
        self.use_name = self.name + self.postfix
        self.destination = self.hosts_dir + '/' + self.use_name
        
        """ status flags """
        self.db_created = False
        self.db_user_created = False
        self.svn_created = False
        self.hg_create = False
        
    """
        Is vhost already exists
    """
    def is_exists(self):
        if (os.path.exists(self.destination)):
            return True
        else:
            return False
            
    """
        Create a directory structure
    """
    def create_directories(self):
        dirs = ['www', 'tmp', 'log']
        for i in dirs:
            dir = '%s/%s/%s' % (self.doc_roots, self.use_name, i)
            try:
                os.makedirs(dir)
            except OSError as (errno, errstr):
                print '%s - already exists' % dir
        
    """
        Save vhost template
    """
    def install_vhost_template(self, vhost_template):
        template = open(vhost_template).read().replace('$name', self.use_name).replace('$projects', self.doc_roots)
        output = open(self.destination, 'w')
        output.write(template)
        output.close()
        
    """
        Save default index.htm as test index file
    """
    def install_index_html(self, html_template):
        template = open(html_template).read().replace('$title', self.use_name)
        dest = '%s/%s/www/index.htm' % (self.doc_roots, self.use_name)
        output = open(dest, 'w')
        output.write(template)
        output.close()
        
    """
        Write lines into hosts file
    """
    def install_hosts(self):
        import vhost_hosts as hosts
        hosts = hosts.hosts(self.hosts_file)
        hosts.add_line('%s        %s' % (self.ini.get_option('general', 'ip'), self.use_name))
        
    """
        Generate password
    """
    def make_password(self):
        import hashlib, time
        return hashlib.md5(str(time.time())).hexdigest()[0:8]
        
    """
        Create MySQL database
    """
    def create_database(self, name, create_user = False):
        binary = self.ini.get_option('mysql', 'executable')
        user = self.ini.get_option('mysql', 'username')
        password = self.ini.get_option('mysql', 'password')
        
        os.system('%s -u%s -p%s --execute="CREATE DATABASE IF NOT EXISTS %s"' % (binary, user, password, name))
        
        """ create user if requested """
        if (create_user):
            user_pass = self.make_password();
            self.mysql_user_password = user_pass
            
            print 'Creating mysql user'
            os.system('%s -u%s -p%s --execute="CREATE USER %s@localhost IDENTIFIED BY \'%s\'"' % (binary, user, password, name, user_pass))
            
            print 'Granting all privilegies'
            os.system('%s -u%s -p%s --execute="GRANT ALL PRIVILEGES on %s.* to %s@localhost IDENTIFIED BY \'%s\'"' % (binary, user, password, name, name, user_pass))
            self.db_user_created = True
            
        self.db_created = True
        
    """
        Collect and save installation data into file and display as result
    """
    def save_data(self, filename):
        from ConfigParser import ConfigParser
        config = ConfigParser()
        config.add_section('install_data')
        config.set('install_data', 'install_dir', self.destination)
        config.set('install_data', 'name', self.name)
        config.set('install_data', 'postfix', self.postfix)
        
        """ save data if db was created """
        if (self.db_created == True):
            config.set('install_data', 'mysql_created', '1')
            config.set('install_data', 'mysql_database', self.name)
            config.set('install_data', 'mysql_username', self.name)
        else:
            config.set('install_data', 'mysql_created', '0')
            
        f = open(filename, 'w')
        config.write(f)
        
        print '\n -----------------------------------------'
        print 'Virtual host %s was successfully created' % self.use_name
        print 'domain: %s' % self.use_name
        
        if (self.db_created):
            print 'mysql database: %s' % self.name
            if (self.db_user_created):
                print 'mysql username: %s' % self.name
                print 'mysql password: %s' % self.mysql_user_password
    
    """
        Backup database
    """
    def backup_database(self, folder):
        user = self.ini.get_option('mysql', 'username')
        password = self.ini.get_option('mysql', 'password')
        
        filename = folder + '/db.sql' 
        os.system('mysqldump -u%s -p%s %s > %s' % (user, password, self.name, filename))

    def backup_host(self, folder):
        import shutil
        _from = self.hosts_dir + '/' + self.use_name
        target = folder + '/' + self.use_name
        shutil.copyfile(_from, target)
        
    """
        Backup www data
    """
    def backup_data(self, folder):
        archiver = 'cd %s && tar -cjf %s ./'
        target = folder + '/www.tar.bz2'
        source = self.doc_roots + '/' + self.use_name
        
        os.system(archiver % (source, target));
        
    """
        Backup meta data
    """
    def backup_meta(self, folder):
        import shutil
        source = self.ini.get_option('general', 'dir_dumps') + self.use_name + '.ini'
        target = folder + '/metadata.ini'
        
        shutil.copyfile(source, target)

    """
        remove vhost
    """
    def remove(self, name):
        host_file = self.hosts_dir + '/' + self.use_name
        info_file = self.ini.get_option('general', 'dir_dumps') + self.use_name + '.ini'
        
        if (os.path.exists(host_file) == True):
            os.remove(host_file)
        else:
            print 'No such vhost'
            
        if (os.path.exists(info_file) == True):
            os.remove(info_file)
        
    """
        remove vhost line from hosts
    """
    def uninstall_hosts(self, line):
        import vhost_hosts as hosts
        hosts = hosts.hosts(self.hosts_file)
        hosts.remove_line(line)

    def apache_reload(self):
        os.system(self.ini.get_option('general', 'reload_command'))

    def enable(self, name):
        target = self.enabled_dir + '/' + self.use_name + '.conf'
        os.symlink(self.destination, target)
        self.install_hosts()

    def disable(self, name):
        target = self.enabled_dir + '/' + self.use_name
        os.remove(target)
        self.uninstall_hosts('%s        %s' % (self.ini.get_option('general', 'ip'), self.use_name))

    def is_enabled(self, name):
        return os.path.exists(self.enabled_dir + '/' + self.use_name)

    def drop_database(self):
        binary = self.ini.get_option('mysql', 'executable')
        user = self.ini.get_option('mysql', 'username')
        password = self.ini.get_option('mysql', 'password')

        os.system('%s -u%s -p%s --execute="DROP DATABASE IF EXISTS %s"' % (binary, user, password, self.name))

    def restore_db(self, db_file):
        binary = self.ini.get_option('mysql', 'executable')
        user = self.ini.get_option('mysql', 'username')
        password = self.ini.get_option('mysql', 'password')

        os.system('%s -u%s -p%s --execute="CREATE DATABASE IF NOT EXISTS %s"' % (binary, user, password, self.name))
        os.system('%s -u%s -p%s -D%s < %s' % (binary, user, password, self.name, db_file))

    def restore_data(self, arc):
        www_dir = self.doc_roots + '/' + self.use_name
        tarc = www_dir + '/www.tar.bz2'
        command = 'cp %s %s && cd %s && tar xjf %s && rm %s' % (arc, www_dir, www_dir, tarc, tarc)
        os.system(command)
        self.chown()

    def restore_hosts(self):
        source = self.ini.get_option('backup', 'backup_dir') + '/' + self.use_name + '/' + self.use_name
        command = 'cp %s %s ' % (source, self.hosts_dir)
        os.system(command)

    def restore_metadata(self, metadata_file):
        target = self.ini.get_option('general', 'dir_dumps') + '/' + self.use_name + '.ini'
        os.system('cp %s %s' % (metadata_file, target))
        self.chown(target)

    def remove_user_data(self):
        target = self.doc_roots + '/' + self.use_name
        os.system('rm -rf %s' % target)

    def chown(self, target = None):
        if (target == None):
            target = self.doc_roots + '/' + self.use_name
            
        user = self.ini.get_option('general', 'user')
        group = self.ini.get_option('general', 'group')
            
        os.system('chown -R %s:%s %s' % (user, group, target))

    """
        Magic getter
    """
    def __getitem__(self, key):
        return self.__dict__[key]
