
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

from .templates import *

# Supported programming languages
C_LANGUAGE = 'C'

def supported_projects():
    """
    Gets all supported project formats.

    :return Returns a list of all supported projects.
    """
    return [base.PTYPE_SOURCE, base.PTYPE_HEADER, base.PTYPE_APPLICATION,
            base.PTYPE_LIBRARY]



def supported_languages():
    """
    Gets all supported languages.

    :return Returns a list of all supported languages.
    """
    return [C_LANGUAGE]



class Template(object):
    """
    Class to create the project, using all options from the user.

    :param args: All arguments received from command line.
    """
    def __init__(self, args):
        self._args = args
        prefix = self._project_prefix()

        self._common_vars = {
            'DATE': time.strftime('%c'),
            'YEAR': time.strftime('%Y'),
            'FULL_AUTHOR_NAME': self._args.author,
            'COMPILER': self._args.compiler,
            'PROJECT_NAME': prefix + \
                    self._args.project_name.replace('-', '_'),
            'PROJECT_NAME_UPPER': prefix.upper() + \
                    self._args.project_name.upper().replace('-', '_')
        }

        self._template = {
            C_LANGUAGE: CTemplate.CTemplate(self._args, self._common_vars,
                                            prefix)
        }.get(self._args.language)


    def _project_prefix(self):
        """
        Gets the project prefix name according to the command line options.
        """
        return {
            base.PTYPE_LIBRARY: 'lib'
        }.get(self._args.project_type, '')


    def create(self):
        self._template.create()


    def info(self):
        self._template.info()



