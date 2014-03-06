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
class hosts():
    def __init__(self, filename):
        self.filename = filename
    
    def read(self, mode = 'a+'):
        return open(self.filename, mode)
    
    def add_line(self, line):
        os.system('echo "%s" >> %s' % (line, self.filename))

    def remove_line(self, msg):
        os.system('sed "/%s/d" %s > /tmp/tmp-hosts' % (msg, self.filename))
        os.system('mv -f /tmp/tmp-hosts %s' % self.filename)
