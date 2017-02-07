
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
The C language project creation.
"""

import os

from string import Template

from . import base, FileTemplate, package, utils, log, license
from .languages import C

class CTemplate(base.BaseTemplate):
    def __init__(self, root_dir, args, project_vars):
        super(CTemplate, self).__init__()
        self._args = args
        self._project_vars = project_vars

        # Store the real info from each source/header which will be created.
        self._files = FileTemplate.FileTemplateInfo(root_dir)
        self._prepare_project_files()


    def _library_header(self, filename):
        """
        Decides which file content will be used for a header file if we're
        creating a library.
        """
        if self._args.project_name in filename:
            if self._args.package is True:
                return Template(C.PACKAGE_DEF_HEADER)\
                            .safe_substitute(self._project_vars)
            else:
                body = C.MAIN_LIBRARY_HEADER
        else:
            body = C.MISC_LIBRARY_HEADER

        body += {
            'error.h': C.LIB_ERROR_HEADER
        }.get(filename, '')

        return Template(body).safe_substitute(self._project_vars)


    def _application_header(self, filename):
        """
        Decides which file content will be used for an application header or a
        single header file if we're creating an application.
        """
        content = ''

        if '_def' in filename:
            if self._args.package is True or package.is_dir():
                content = \
                    Template(C.PACKAGE_DEF_HEADER)\
                        .safe_substitute(self._project_vars)
            else:
                content = \
                    Template(C.DEF_HEADER).safe_substitute(self._project_vars)
        else:
            if self._args.content is not None:
                content = Template(self._args.content.replace('^', '\n'))\
                                .safe_substitute(self._project_vars)
            elif self._args.project_name == filename and \
                    self._args.project_type not in (utils.PTYPE_SOURCE, \
                                                    utils.PTYPE_HEADER):
                content = Template(C.MAIN_HEADER)\
                                .safe_substitute(self._project_vars)

        return content


    def _header_files_to_include(self):
        """
        Creates a list of #include directives to be inserted into one of the
        project header files.
        """
        include = ''

        if self._args.package is True and \
                self._args.project_type == utils.PTYPE_APPLICATION:
            # Empty
            return ''

        for key, value in self._files.files().iteritems():
            if value.get('header'):
                include += '\n#include "%s"' % key

        if len(include):
            include += '\n'

        return include


    def _get_header_content(self, filename):
        """
        Returns a content to insert in header files.
        """
        self._project_vars['HEADER_FILES'] = self._header_files_to_include();
        self._project_vars['FILENAME'] = utils.assemble_filename(filename, '.h')
        self._project_vars['FILENAME_UPPER'] = \
                utils.assemble_filename(filename, '.h').replace('.', '_')\
                                .replace('-', '_').upper()

        get_header = {
            utils.PTYPE_LIBRARY: self._library_header,
        }.get(self._args.project_type, self._application_header)

        return Template(C.HEADER_HEAD + get_header(filename) + C.HEADER_TAIL)\
                .safe_substitute(self._project_vars)


    def _add_source(self, filename, path='src', extension=C.SOURCE_EXTENSION,
                    comment=True):
        """
        Adds a file into the internal FileTemplate object.
        """
        content = None
        source = False
        header = False

        if self._args.license is None:
            c_head = Template(C.HEAD).safe_substitute(self._project_vars)
        else:
            c_head = Template(C.HEAD_LICENSE)\
                        .safe_substitute(self._project_vars) %\
                        license.license_block(self._args.license,
                                              self._project_vars,
                                              comment_char=' *')

        if extension == C.HEADER_EXTENSION:
            content = self._get_header_content(filename)
            header = True
        else:
            source = True

            # Did we receive a content from the command line?
            if self._args.content is not None:
                content = Template(self._args.content.replace('^', '\n'))\
                                .safe_substitute(self._project_vars)
            else:
                content = {
                    'main': Template(C.MAIN_SOURCE)\
                                .safe_substitute(self._project_vars),
                    'error': Template(C.LIB_ERROR_SOURCE)\
                                .safe_substitute(self._project_vars)
                }.get(filename)

        self._files.add(filename, path, head=c_head, body=content,
                        source=source, header=header)

        self._files.set_property(filename, 'extension', extension)


    def _add_makefile(self):
        """
        Adds a Makefile to the project.
        """
        if self._args.package is True:
            mcontent = {
                utils.PTYPE_APPLICATION: C.APP_MAKEFILE_PACKAGE,
                utils.PTYPE_LIBRARY: C.LIB_MAKEFILE_PACKAGE
            }.get(self._args.project_type)
        else:
            mcontent = {
                utils.PTYPE_APPLICATION: C.APP_MAKEFILE,
                utils.PTYPE_LIBRARY: C.LIB_MAKEFILE
            }.get(self._args.project_type)

        self._files.add('Makefile', 'src',
                        body=Template(mcontent)\
                                 .safe_substitute(self._project_vars))


    def _prepare_project_files(self):
        """
        Prepare all project files (sources and headers).
        """
        app_name = self._args.project_name.lower()

        # Is just a single file?
        if self._args.project_type in (utils.PTYPE_SOURCE, utils.PTYPE_HEADER):
            extension = {
                utils.PTYPE_SOURCE: C.SOURCE_EXTENSION,
                utils.PTYPE_HEADER: C.HEADER_EXTENSION
            }.get(self._args.project_type)

            self._add_source(self._args.project_name, '', extension, False)
            return
        else:
            self._add_makefile()

        if self._args.project_type in (utils.PTYPE_APPLICATION, \
                utils.PTYPE_LIBCOLLECTION_APP):
            self._add_source(app_name, 'include', C.HEADER_EXTENSION)
            self._add_source('main')

            for suffix in ['_prt', '_def', '_struct']:
                self._add_source(app_name + suffix, 'include',
                                 C.HEADER_EXTENSION)

            if self._args.project_type == utils.PTYPE_LIBCOLLECTION_APP:
                for filename in ['log', 'config', 'core']:
                    self._add_source(filename)

        if self._args.project_type == utils.PTYPE_LIBRARY:
            self._add_source('utils')
            self._add_source('error')
            self._add_source('utils.h', 'include', C.HEADER_EXTENSION)
            self._add_source('error.h', 'include', C.HEADER_EXTENSION)

            # We must add this file as the last one, since we'll need to
            # add a few includes inside it.
            self._add_source(self._args.prefix + app_name, 'include',
                             C.HEADER_EXTENSION)

            self._files.add(self._args.prefix + app_name + '.sym', 'src',
                            body=Template(C.LIBSYM)\
                                     .safe_substitute(self._project_vars))

        for filename in self._args.sources:
            self._add_source(filename)

        for filename in self._args.headers:
            self._add_source(filename, 'include', C.HEADER_EXTENSION)


    def create(self):
        self._files.save_all()


    def info(self):
        # TODO: print project description
        pass



