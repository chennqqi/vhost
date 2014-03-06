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

from argparse import ArgumentParser

def get_parser():
    parser = ArgumentParser(
        description='Simple tool for basic managing of apache\'s virtual hosts',
        epilog = 'Bug reports send to alex.oleshkevich@gmail.com'
    )

    parser.add_argument(
        "-c", "--create",
        action="store_true",  default=False,
        dest="create", help="create a new virtual host"
    )
    
    parser.add_argument(
        "-e", "--enable",
        action="store_true",  default=False,
        dest="enable", help="enable a virtual host"
    )

    parser.add_argument(
        "-d", "--disable",
        action="store_true",  default=False,
        dest="disable", help="disable a virtual host"
    )


    parser.add_argument(
        "-p", "--purge",
        action="store_true",  default=False,
        dest="purge", help="remove vhost"
    )

    parser.add_argument(
        "-m", "--mysql",
        action="store_true",  default=False,
        dest="database", help="create a new database with the same name as vhost"
    )

    parser.add_argument(
        "-t", "--test",
        action="store_true",  default=False,
        dest="test", help="test if vhost exists"
    )

    parser.add_argument(
        "-b", "--backup",
        action="store_true",  default=False,
        dest="backup", help="backup vhost"
    )

    parser.add_argument(
        "-r", "--restore",
        action="store_true",  default=False,
        dest="restore", help="restore backup of vhost"
    )
    
    parser.add_argument(
        "-n", "--noindex",
        action="store_false",  default=None,
        dest="noindex", help="don't create index.htm"
    )
    
    parser.add_argument(
        "-f", "--remove-data",
        action="store_true",  default=None,
        dest="remove_data", help="also remove www data as well as vhost"
    )
    
    parser.add_argument(
        "-z", "--in-public",
        action="store_true",  default=None,
        dest="remove_data", help="create vhost in subdirectory of DOCUMENT ROOT"
    )
    
    parser.add_argument(
        "--db-create-user",
        action="store_true",  default=None,
        dest="db_create_user", help="create a database user"
    )
    
    parser.add_argument(
        "name",
        action="store",  default=False,
        help="vhost name"
    )

    return parser