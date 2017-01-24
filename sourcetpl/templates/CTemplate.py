
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

comment = '''
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

app_makefile = '''.PHONY: outputdirs

CC = $COMPILER

machine = $(shell uname -m)

ifeq ($(machine), x86_64)
    ARCH_DIR = x86_64
else
    ARCH_DIR = x86
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

lib_makefile = '''
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


class CTemplate(base.BaseTemplate):
    def __init__(self, args, project_vars, prefix):
        self._args = args
        self._project_vars = project_vars
        self._prefix = prefix
        print 'ola'


    def __make_header_content(self, filename):
        u = filename.replace('.', '_').upper()
        s = '''
#ifndef _%s
#define _%s     1

#endif

''' % (u, u)

        return s


    def __create_single_file(self, project_type='', filename='', dest_dir=''):
        """
        Saves the template of a single source/header file.
        """
        if len(filename) == 0:
            filename = self._args.project_name

        if len(project_type) == 0:
            project_type = self._args.project_type

        if project_type == base.PTYPE_SOURCE:
            if ".c" not in filename:
                filename += '.c'

            content = ''
        elif project_type == base.PTYPE_HEADER:
            if ".h" not in filename:
                filename += '.h'

            content = self.__make_header_content(filename)

        if len(dest_dir) != 0:
            dest_dir += '/'

        dest_dir += filename
        output = Template(comment).safe_substitute(self._project_vars)

        # TODO: replace with with
        fd = open(dest_dir, 'w')
        fd.write(output)

        if len(content) != 0:
            fd.write(content)

        fd.close()


    def __create_directory_structure(self):
        root_dirname = self._prefix + \
                self._args.project_name.replace('_', '-')

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

        return root_dirname, subdirs


    def __create_makefile(self, dest_dir=''):
        mcontent = {
            base.PTYPE_APPLICATION: app_makefile,
            base.PTYPE_LIBRARY: lib_makefile
        }.get(self._args.project_type)

        output = Template(mcontent).safe_substitute(self._project_vars)

        if len(dest_dir) != 0:
            dest_dir += '/'

        dest_dir += 'Makefile'

        # TODO: replace with with
        fd = open(dest_dir, 'w')
        fd.write(output)
        fd.close()


    def create(self):
        if self._args.project_type in (base.PTYPE_SOURCE, base.PTYPE_HEADER):
            self.__create_single_file()
            return

        sources = self._args.sources
        headers = self._args.headers

        try:
            dirs = self.__create_directory_structure()
        except OSError as e:
            print e
            return -1

        src_dir = dirs[0] + '/src'
        header_dir = dirs[0] + '/include'

        for f in sources:
            self.__create_single_file(base.PTYPE_SOURCE, f, src_dir)

        for f in headers:
            self.__create_single_file(base.PTYPE_HEADER, f, header_dir)

        self.__create_makefile(src_dir)

        return 0


    def info(self):
        # TODO: print project description
#        if self._args.quiet is False:
#            print self.options

        pass



