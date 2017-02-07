
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

import os
import time

class FileTemplateInfo(object):
    """
    A class to store information about files which will be created to a
    specific project requested by the user.

    Each file is stored in an internal dictionary, containing properties so
    we can correctly create them later.

    :param root_pathname: The path where all files will be saved.
    """
    def __init__(self, root_pathname):
        self._templates = dict()
        self._root_pathname = root_pathname


    def files(self):
        """
        Returns the templates stored inside.
        """
        return dict(self._templates)


    def add(self, filename, base_path, head=None, tail=None, body=None,
            source=False, header=False, executable=False):
        """
        Adds a file to our internal dictionary, informing if it's a source or a
        header file or none of them, its path do be written and if it's a
        executable file.

        It's optional to inform the file extension, but if none is used and in
        the end we want one, we have to add an extra property: 'extension'.
        """
        self._templates[filename] = {
            'head': head,
            'body': body,
            'tail': tail,
            'source': source,
            'header': header,
            'executable': executable,
            'path': base_path
        }


    def set_property(self, filename, property_name, property_data):
        self._templates[filename][property_name] = property_data


    def _save_one(self, filename, properties):
        """
        Save a file to the respective project.
        """
        extension = properties.get('extension')

        if extension and '.' not in filename:
            if '.' not in extension:
                filename += '.'

            filename += extension

        if len(self._root_pathname):
            subdir = properties.get('path')
            full_path = self._root_pathname + '/' + subdir
            full_filename = full_path + '/' + filename

            if not os.access(full_path, os.F_OK):
                os.makedirs(full_path)
        else:
            full_filename = filename

        with open(full_filename, 'w') as fd:
            try:
                fd.write(properties.get('head'))
            except:
                pass

            try:
                fd.write(properties.get('body'))
            except:
                pass

            try:
                fd.write(properties.get('tail'))
            except:
                pass

        if properties.get('executable', False) is True:
            os.system('chmod +x %s' % full_filename)


    def save_all(self):
        """
        Save all internal files to their respective names and content. All files
        will use self._root_pathname as its base directory.
        """
        for filename in self._templates.keys():
            self._save_one(filename, self._templates.get(filename))



