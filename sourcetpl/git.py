
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
Functions to handle git files to include in a project.
"""

import os
from string import Template

from . import utils, FileTemplate, license
from .languages import C, python

README = '''# ${PROJECT_NAME}
A brief description of the project.

'''

class Git(object):
    def __init__(self, root_dir, args, project_vars):
        self._args = args
        self._project_vars = project_vars

        if self._args.package is True:
            root_dir = os.path.dirname(root_dir)

        self._files = FileTemplate.FileTemplateInfo(root_dir)
        self._prepare_files()


    def _prepare_files(self):
        """
        Adds all required files to a git project.
        """
        gitignore = {
            utils.C_LANGUAGE: C.GITIGNORE,
            utils.PYTHON_LANGUAGE: python.GITIGNORE
        }.get(self._args.language)

        # .gitignore
        self._files.add('.gitignore', '',
                        body=Template(gitignore)\
                                .safe_substitute(self._project_vars))

        # README.md
        self._files.add('README.md', '',
                        body=Template(README)\
                                .safe_substitute(self._project_vars))

        # LICENSE
        if self._args.license is not None:
            license_cnt = license.license(self._args.license)

            if license_cnt is not None:
                self._files.add('LICENSE', '',
                                body=Template(license_cnt)\
                                    .safe_substitute(self._project_vars))


    def create(self):
        self._files.save_all()



