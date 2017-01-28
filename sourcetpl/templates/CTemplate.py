
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

from . import base
from . import FileTemplate

HEADER_EXTENSION = '.h'
SOURCE_EXTENSION = '.c'

COMMENT = '''
/*
 * Description:
 *
 * Author: $FULL_AUTHOR_NAME
 * Created at: $DATE
 * Project: $PROJECT_NAME
 *
 * Copyright (c) $YEAR All rights reserved
 */

'''

APP_MAKEFILE = '''.PHONY: outputdirs

CC = $COMPILER

machine = $(shell uname -m)

ifeq ($(machine), x86_64)
    ARCH_DIR = x86_64
else
    ARCH_DIR = i686
endif

OUTPUTDIR = ../bin/$(ARCH_DIR)
TARGET = $(OUTPUTDIR)/$PROJECT_NAME

INCLUDEDIR = -I../include

CFLAGS = -Wall -Wextra -O0 -ggdb $(INCLUDEDIR)

LIBDIR = -L/usr/local/lib
LIBS =

C_FILES := $(wildcard *.c)
OBJS = $(C_FILES:.c=.o)

$(TARGET): outputdirs $(OBJS)
	$(CC) -o $(TARGET) $(OBJS) $(LIBDIR) $(LIBS)

clean:
	rm -rf $(OBJS) $(TARGET) *~ ../include/*~

purge: clean $(TARGET)
outputdirs: $(OUTPUTDIR)
$(OUTPUTDIR):
	mkdir -p $(OUTPUTDIR)
'''

LIB_MAKEFILE = '''
CC = $COMPILER

PCT_VERSION = $(shell grep -w VERSION ../include/${PROJECT_NAME}_internal.h | awk '{print $$$4}' | cut -d \\" -f 2)
MAJOR_VERSION = $(shell grep MAJOR_VERSION ../include/${PROJECT_NAME}_internal.h | awk '{print $$$4}')
MINOR_VERSION = $(shell grep MINOR_VERSION ../include/${PROJECT_NAME}_internal.h | awk '{print $$$4}')

USR_DIR = /usr/local/lib
LIBNAME = $PROJECT_NAME.so
SONAME = $(LIBNAME).$(PCT_VERSION).$(MAJOR_VERSION)
DEST_LIBNAME = $(SONAME).$(MINOR_VERSION)

TARGET = ../bin/$(SONAME)
TARGET_DEST = ../bin/$(DEST_LIBNAME)

INCLUDEDIR = -I../include

LIBDIR =
LIBS =

CFLAGS = -Wall -Wextra -fPIC -ggdb -O0 -g3 -fvisibility=hidden \\
        -D${PROJECT_NAME_UPPER}_COMPILE -D_GNU_SOURCE $(INCLUDEDIR)

VPATH = ../include:.

OBJS = 		\\
	common.o

$(TARGET): $(OBJS)
	rm -f ../bin/$(LIBNAME)*
	$(CC) -shared -Wl,-soname,$(SONAME) -o $(TARGET_DEST) $(OBJS) $(LIBDIR) $(LIBS)

clean:
	rm -rf $(OBJS) $(TARGET) $(TARGET_DEST) *~ ../include/*~

dest_clean:
	rm -f $(USR_DIR)/$(LIBNAME)*

purge: clean $(TARGET)

install:
	cp -f $(TARGET_DEST) $(USR_DIR)
	rm -rf $(USR_DIR)/$(LIBNAME) $(USR_DIR)/$(SONAME)
	ln -s $(USR_DIR)/$(DEST_LIBNAME) $(USR_DIR)/$(LIBNAME)
	ln -s $(USR_DIR)/$(DEST_LIBNAME) $(USR_DIR)/$(SONAME)
'''

DEF_HEADER = '''
#define MAJOR           0
#define MINOR           1
#define RELEASE         1
#define BETA            true
#define BUILD           0
'''

PACKAGE_DEF_HEADER = '''
/*
 * Package version: major, minor and release.
 */
#include "../../package_version.h"

#define BUILD           0
'''

class CTemplate(base.BaseTemplate):
    def __init__(self, args, project_vars):
        self._args = args
        self._project_vars = project_vars

        # Store the real info from each source/header which will be created.
        self._files = FileTemplate.FileTemplateInfo()
        self._prepare_project_files()


    def _get_header_content(self, filename):
        content = ''
        upper_filename = filename.replace('.', '_').upper()

        if '_def' in filename:
            if self._args.package is True:
                content = \
                    Template(PACKAGE_DEF_HEADER)\
                        .safe_substitute(self._project_vars)
            else:
                content = \
                    Template(DEF_HEADER).safe_substitute(self._project_vars)
        else:
            if self._args.content is not None:
                content = Template(self._args.content)\
                                .safe_substitute(self._project_vars)

        return '''#ifndef _%s_H
#define _%s_H     1
%s
#endif

''' % (upper_filename, upper_filename, content)


    def _add_file(self, filename, path='src', extension=SOURCE_EXTENSION):
        """
        Adds a file into the internal FileTemplate object. Here we also
        """
        content = None

        if extension == HEADER_EXTENSION:
            content = self._get_header_content(filename)
        else:
            if self._args.content is not None:
                content = Template(self._args.content)\
                                .safe_substitute(self._project_vars)

        self._files.add(filename, path, content)
        self._files.set_property(filename, 'extension', extension)
        self._files.set_property(filename, 'comment', True)


    def _add_makefile(self):
        """
        Adds a Makefile to the project.
        """
        mcontent = {
            base.PTYPE_APPLICATION: APP_MAKEFILE,
            base.PTYPE_LIBRARY: LIB_MAKEFILE
        }.get(self._args.project_type)

        self._files.add('Makefile', 'src',
                        Template(mcontent).safe_substitute(self._project_vars))


    def _prepare_project_files(self):
        """
        Prepare all project files (sources and headers).
        """
        app_name = self._args.project_name.lower()

        # Is just a single file?
        if self._args.project_type in (base.PTYPE_SOURCE, base.PTYPE_HEADER):
            extension = {
                base.PTYPE_SOURCE: SOURCE_EXTENSION,
                base.PTYPE_HEADER: HEADER_EXTENSION
            }.get(self._args.project_type)

            self._add_file(self._args.project_name, '', extension)
            return
        else:
            self._add_makefile()

        if self._args.project_type in (base.PTYPE_APPLICATION, \
                base.PTYPE_LIBCOLLECTION_APP):
            self._add_file(app_name, 'include', HEADER_EXTENSION)
            self._add_file('main')

            for suffix in ['_prt', '_def', '_struct']:
                self._add_file(app_name + suffix, 'include', HEADER_EXTENSION)

            if self._args.project_type == base.PTYPE_LIBCOLLECTION_APP:
                for filename in ['log', 'config', 'core']:
                    self._add_file(filename)

        if self._args.project_type == base.PTYPE_LIBRARY:
            self._add_file('utils')
            self._add_file(self._args.prefix + app_name, 'include',
                           HEADER_EXTENSION)

        for filename in self._args.sources:
            self._add_file(filename)

        for filename in self._args.headers:
            self._add_file(filename, 'include', HEADER_EXTENSION)


    def _create_single_file(self, filename, root_dir):
        file_data = self._files.properties(filename)
        extension = file_data.get('extension', '')
        path = file_data.get('path')
        content = file_data.get('data')
        comment = file_data.get('comment', False)

        if extension not in filename:
            filename += extension

        if len(root_dir):
            pathname = root_dir + '/' + path + '/' + filename
        else:
            pathname = filename

        if comment is True:
            comment_cnt = Template(COMMENT).safe_substitute(self._project_vars)

        with open(pathname, 'w') as out_fd:
            if comment is True:
                out_fd.write(comment_cnt)

            if content is not None:
                out_fd.write(content)


    def _create_directory_structure(self):
        """
        Creates all projects required directories.
        """
        root_dirname = self._args.prefix + \
                self._args.project_name.replace('_', '-')

        if self._args.package is True:
            root_dirname = self._args.root_dir + '/' + root_dirname

        try:
            os.mkdir(root_dirname)
        except OSError:
            raise

        subdirs = ['src', 'include', 'bin', 'po']

        if self._args.project_type != base.PTYPE_LIBRARY:
            subdirs.append('doc')

        for d in subdirs:
            try:
                os.mkdir(root_dirname + '/' + d)
            except OSError:
                raise

        return root_dirname


    def create(self):
        if self._args.project_type not in (base.PTYPE_SOURCE, base.PTYPE_HEADER):
            try:
                root_dir = self._create_directory_structure()
            except:
                raise
        else:
            root_dir = ''

        for filename in self._files.filenames():
            self._create_single_file(filename, root_dir)


    def info(self):
        # TODO: print project description
        pass



