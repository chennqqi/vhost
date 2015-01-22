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

import sys, os, getpass
import vhost_args as args
import vhost_config as config
import vhost_class

global dir_dumps
global vhost
global ini

dir_install = os.getenv("HOME") + '/.vhost/'
file_config = dir_install + '/config.ini'
file_vhost_template = dir_install + '/share/vhost.conf'
file_html_template = dir_install + '/share/index.htm'
dir_dumps = dir_install + '/var/'

"""
Entry point
"""
def main():
    if (getpass.getuser() != 'root'):
        print '* You must be root to use this program.'
        sys.exit()
    
    parser = args.get_parser();
    options = parser.parse_args()
    name = options.name
    ini = config.config(file_config)
    ini.set_option('general', 'dir_install', dir_install)
    ini.set_option('general', 'dir_dumps', dir_dumps)
    
    vhost = vhost_class.vhost_class(name, ini)
    
    """ create a new vhost """
    if (options.create == True):
        """ Check for existing vhost """
        if (vhost.is_exists()):
            print '* Virtual host %s already exists.' % vhost.use_name
            sys.exit()
        
        print '* Ok. Fast forward ---->'
        print '* Creation of virtual host: %s' % vhost.use_name
        
        """ creating a directory structure """
        print '* Making directory structure'
        vhost.create_directories()
        
        """ parse template and replace vars """
        print '* Installing vhost template'
        vhost.install_vhost_template(file_vhost_template)
        
        """ Which way to use? bind9 or hosts file? """
        if (ini.get_option('apache', 'dns_mode') == 'hosts'):
            print '* Name-based vhosts enabled. Using %s' % vhost.hosts_file
            vhost.install_hosts()
       
        """ install index.htm """
        install_index = ini.get_option('general', 'install_index_html') == '1'
        if (options.noindex != None):
            install_index = options.noindex
            
        if (install_index == True):
            print '* Installing index.htm file'
            vhost.install_index_html(file_html_template)
            
        """ create database + user with all privilegies on it """
        if (options.database ==  True):
            mysql_binary = ini.get_option('mysql', 'executable')
            if (os.path.exists == False):
                print '* No mysql found on this machine'
                sys.exit()
                
            """ detect if mysql user should also created """
            create_user = ini.get_option('mysql', 'also_create_user') == '1'
            if (options.db_create_user != None):
                create_user = options.db_create_user
                
            """ create database """
            vhost.create_database(name, create_user)

        vhost.chown()
        vhost.enable(name)
        vhost.apache_reload()
        vhost.save_data(dir_dumps + vhost.use_name + '.ini')
            
        """ perform deletion """
    elif (options.purge == True):
        if (vhost.is_exists() == False):
            print '* Vhost does not exists'
            sys.exit()
        
        print '* Virtual host %s will be deleted' % vhost.use_name
        print '* Loading installation data'

        try:
            install = config.config(dir_dumps + vhost.use_name + '.ini')
        except:
            pass
        
        if (options.backup == True):
            print '* Backup before purging will be done'
            backup(vhost, ini)

        if (options.database == True):
            choice = raw_input('** Database is to be deleted. Continue? (y/n)')
            if (choice == 'y'):
                vhost.drop_database()
            else:
                print '* Skipped'

        if (options.remove_data == True):
            choice = raw_input('** WWW data is to be deleted. Continue? (y/n)')
            if (choice == 'y'):
                vhost.remove_user_data()
            else:
                print '* Skipped'
                
        if (vhost.is_enabled(name)):
            vhost.disable(name)

        print '* Removing vhost file'
        vhost.remove(name)
        vhost.uninstall_hosts(name)
        vhost.apache_reload()
        
        """ perform backup """
    elif (options.backup == True):
        backup(vhost, ini)
        
        """ perform restore from backup """
    elif (options.restore == True):
        print '* Restoring of %s' % name
        path = ini.get_option('backup', 'backup_dir') + '/' + vhost.use_name
        if (os.path.exists(path) == False):
            print '** No backup found'
            sys.exit()

        if (os.path.exists(path + '/db.sql') == False):
            print '** No database backup found. Skipped.'
        else:
            vhost.restore_db(path + '/db.sql')
            print '* Database restored'

        if (os.path.exists(path + '/www.tar.bz2') == False):
            print '** No user-data backup archive found'
            sys.exit()
        else:
            if (os.path.exists(vhost.doc_roots + '/' + vhost.use_name) == False):
                print '* Restoring directory structure'
                os.makedirs(vhost.doc_roots + '/' + vhost.use_name)
                
            vhost.restore_data(path + '/www.tar.bz2')
            print '* Data restored'

        if (os.path.exists(path + '/metadata.ini') == False):
            print '** No metadata backup found. Skipped'
        else:
            vhost.restore_metadata(path + '/metadata.ini')
            print '* Metadata restored'

        vhost.restore_hosts()
        vhost.install_hosts()
        vhost.enable(name)
        vhost.apache_reload()

    elif (options.enable == True):
        if not vhost.is_exists():
            print '* Vhost %s does not exists.' % name
        elif vhost.is_enabled(name):
            print '* Vhost %s is already enabled' % name
        else:
            vhost.enable(name)
            vhost.apache_reload()        

    elif (options.disable == True):
        if (vhost.is_enabled(name)):
            vhost.disable(name)
            vhost.apache_reload()
        else:
            print '* Vhost %s is not enabled' % name
        
        """ test for existance """
    elif (options.test == True):
        if (vhost.is_exists() == True):
            print '* Vhost exists'
        else:
            print '* Vhost does not exists'
    
        """ invalid arguments passed """
    else:
        parser.print_help()
        
""" backup operation """
def backup(vhost, ini):
    print '* Virtual host %s will be backuped' % vhost.use_name
    
    backup_dir = ini.get_option('backup', 'backup_dir')
    proj_backup_dir = backup_dir + '/' + vhost.use_name
    
    if (os.path.exists(backup_dir) == False):
        print '* Backup directory does not exists in %s. Creating...' % backup_dir
        os.makedirs(proj_backup_dir)
        
    if (os.path.exists(proj_backup_dir) == False):
        print '* Creating directory for current project'
        os.makedirs(proj_backup_dir)
    
    print '* Loading installation data'
    ini = config.config(dir_dumps + vhost.use_name + '.ini')
    
    if (ini.get_option('install_data', 'mysql_created') == '1'):
        print '* Backuping database'
        vhost.backup_database(proj_backup_dir)
    
    print '* Backuping data'
    vhost.backup_data(proj_backup_dir)

    vhost.backup_host(proj_backup_dir)
    
    print '* Backuping metadata'
    vhost.backup_meta(proj_backup_dir)
        
if __name__ == "__main__":
    main()
