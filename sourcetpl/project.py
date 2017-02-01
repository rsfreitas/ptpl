
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
Project template creation.
"""

import time
import collections
import os
import glob

from . import package, languages, CTemplate, base

def supported_projects():
    """
    Gets all supported project formats.

    :return Returns a dictionary with all supported project as keys and a brief
            description of each one.
    """
    return {
        base.PTYPE_SOURCE : 'Indicates the creation of a single source file.',
        base.PTYPE_HEADER : 'Indicates the creation of a single header file.',
        base.PTYPE_APPLICATION :
            '''Indicates the creation of a directory with the following
            \t\tstructure: $name/{include,src}, containing template files for
            \t\ta single application (with a main function).''',

        base.PTYPE_LIBRARY:
            '''Indicates the creation of a directory to hold a library
            \t\tdevelopment project, with a specific Makefile.''',

        base.PTYPE_LIBCOLLECTION_APP :
            '''Indicates the creation of an application using
            \t\tlibcollections as its base.''',

        base.PTYPE_LIBCOLLECTION_C_PLUGIN :
            '''Indicates the creation of a libcollections' plugin (C)''',

        base.PTYPE_LIBCOLLECTION_PY_PLUGIN :
            '''Indicates the creation of a libcollections' plugin (Python)''',

        base.PTYPE_LIBCOLLECTION_JAVA_PLUGIN :
            '''Indicates the creation of a libcollections' plugin (Java)'''
    }



def supported_projects_description():
    """
    Returns a string containing a formatted output of all supported projects.
    """
    data = collections.OrderedDict(sorted(supported_projects().items()))
    description = ''

    for key, value in data.iteritems():
        description += '%-25s - %s\n' % (key, value)

    return description



def supported_languages():
    """
    Gets all supported languages.

    :return Returns a list of all supported languages.
    """
    return [languages.C_LANGUAGE]



def _is_project_dir():
    """
    Checks if the current directory belongs to a project (an application or
    a library).

    :return Returns the project name if the directory is from a project or
            None otherwise.
    """
    pwd = os.getcwd()

    for path in ['/../include/*.h', '/*.c']:
        if len(glob.glob(pwd + path)):
            return os.path.basename(os.path.dirname(pwd))

    return None



class Template(object):
    """
    Class to create the project, using all options from the user.

    :param args: All arguments received from command line.
    """
    def __init__(self, args):
        self._args = args
        self._args.prefix = self._project_prefix()
        self._args.root_dir = ''

        self._common_vars = {
            'DATE': time.strftime('%c'),
            'YEAR': time.strftime('%Y'),
            'FULL_AUTHOR_NAME': self._args.author,
            'COMPILER': self._args.compiler,
            'PROJECT_NAME': self._args.prefix + \
                    self._args.project_name.replace('-', '_'),
            'PROJECT_NAME_UPPER': self._args.prefix.upper() + \
                    self._args.project_name.upper().replace('-', '_'),
            'PROJECT_BIN_NAME': self._args.prefix + \
                    self._args.project_name.replace('_', '-')
        }

        # Disable package flag if we're creating a single file and get the
        # project name to use in the template.
        if self._args.project_type in (base.PTYPE_SOURCE, base.PTYPE_HEADER):
            self._args.package = False
            single_project_name = _is_project_dir()

            if single_project_name:
                self._common_vars['SINGLE_FILE_PROJECT_NAME'] = \
                        single_project_name

        if self._args.package is True:
            self._package = package.Package(args, self._common_vars)
            self._args.root_dir = self._package.current_dir()

        self._template = {
            languages.C_LANGUAGE: CTemplate.CTemplate(self._args,
                                                      self._common_vars)
        }.get(self._args.language)

        # TODO: Download the code license


    def _project_prefix(self):
        """
        Gets the project prefix name according to the command line options.
        """
        return {
            base.PTYPE_LIBRARY: 'lib'
        }.get(self._args.project_type, '')


    def create(self):
        """
        Create our template. If a package is required, we must create the
        templates inside its directory.
        """
        if self._args.package is True:
            self._package.create()

        self._template.create()


    def info(self):
        self._template.info()



