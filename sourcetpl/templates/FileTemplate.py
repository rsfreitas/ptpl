
#
# Copyright (C) 2015 Rodrigo Freitas
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 2 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License along
# with this program; if not, write to the Free Software Foundation, Inc.,
# 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
#

"""
A module to handle our FileTemplate..
"""

class FileTemplateInfo(object):
    """
    A class to store information about files.

    Each file is stored in an internal dictionary, containing properties so
    we can correctly create them later.
    """
    def __init__(self, files=[]):
        self._templates = dict()
        self._load_files(files)
        self._index = -1


    def _load_files(self, files):
        for filename in files:
            self._templates[filename] = { 'empty': True }


    def add(self, filename, path, data=None, executable=False):
        self._templates[filename] = {
            'empty': False,
            'path': path,
            'chmod': executable
        }

        if data is not None:
            self._templates[filename]['data'] = data


    def set_property(self, filename, property_name, property_data):
        self._templates[filename][property_name] = property_data


    def filenames(self):
        return self._templates.keys()


    def properties(self, filename):
        return self._templates.get(filename)



