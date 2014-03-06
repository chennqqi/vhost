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

from ConfigParser import RawConfigParser
 
class config():
    def __init__(self, file):
        self.config = RawConfigParser()
        self.config.readfp(open(file))
        
    def get_option(self, section, key):
        return self.config.get(section, key)
        
    def set_option(self, section, key, value):
        return self.config.set(section, key, value)