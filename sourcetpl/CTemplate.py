
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

from . import base, FileTemplate, package

HEADER_EXTENSION = '.h'
SOURCE_EXTENSION = '.c'

COMMENT = '''
/*
 * Description:
 *
 * Author: $FULL_AUTHOR_NAME
 * Created at: $DATE
 * Project: $PROJECT_BIN_NAME
 *
 * Copyright (c) $YEAR All rights reserved
 */

'''

COMMENT_SINGLE = '''
/*
 * Description:
 *
 * Author: $FULL_AUTHOR_NAME
 * Created at: $DATE
 * Project: $SINGLE_FILE_PROJECT_NAME
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
TARGET = $(OUTPUTDIR)/$PROJECT_BIN_NAME

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

APP_MAKEFILE_PACKAGE = '''.PHONY: outputdirs package_version

CC = $COMPILER

machine = $(shell uname -m)

ifeq ($(machine), x86_64)
    ARCH_DIR = x86_64
else
    ARCH_DIR = i686
endif

OUTPUTDIR = ../bin/$(ARCH_DIR)
TARGET = $(OUTPUTDIR)/$PROJECT_BIN_NAME

INCLUDEDIR = -I../include

CFLAGS = -Wall -Wextra -O0 -ggdb $(INCLUDEDIR)

LIBDIR = -L/usr/local/lib
LIBS =

C_FILES := $(wildcard *.c)
OBJS = $(C_FILES:.c=.o)

$(TARGET): outputdirs package_version $(OBJS)
	$(CC) -o $(TARGET) $(OBJS) $(LIBDIR) $(LIBS)

clean:
	rm -rf $(OBJS) $(TARGET) *~ ../include/*~

purge: clean $(TARGET)

outputdirs: $(OUTPUTDIR)
$(OUTPUTDIR):
	mkdir -p $(OUTPUTDIR)

PACKAGE_CONF=package/package.conf
PACKAGE_VERSION_NAME=package_version
PACKAGE_VERSION=../../$(PACKAGE_VERSION_NAME).h
package_version: $(PACKAGE_VERSION)
$(PACKAGE_VERSION):
	$(shell (cd ../../ && source-tpl -t header -n $(PACKAGE_VERSION_NAME) \\
	    -c "^#define MAJOR_VERSION	`cfget -C $(PACKAGE_CONF) version/major`\\
	    ^#define MINOR_VERSION	`cfget -C $(PACKAGE_CONF) version/minor` \\
	    ^#define RELEASE		`cfget -C $(PACKAGE_CONF) version/release` \\
	    ^#define BETA		`cfget -C $(PACKAGE_CONF) version/beta`^"))
'''

LIB_MAKEFILE = '''.PHONY: shared static clean dest_clean install outputdirs

CC = $COMPILER
AR = ar

ARCH_TEST := $(shell uname -m)

ifeq ($(ARCH_TEST), x86_64)
    ARCH = x86_64
else
    ARCH = i686
endif

MAJOR_VERSION := $(shell command grep MAJOR_VERSION ../include/${PROJECT_NAME}.h | awk '{print $$$4}')
MINOR_VERSION := $(shell command  grep MINOR_VERSION ../include/${PROJECT_NAME}.h | awk '{print $$$4}')
RELEASE := $(shell command grep RELEASE ../include/${PROJECT_NAME}.h | awk '{print $$$4}')

USR_DIR = /usr/local/lib
PREFIX = ${PROJECT_NAME}
LIBNAME = $(PREFIX).so
SONAME = $(LIBNAME)
SHARED_LIBNAME := $(LIBNAME).$(MAJOR_VERSION).$(MINOR_VERSION).$(RELEASE)
STATIC_LIBNAME := $(PREFIX).a

OUTPUTDIR = ../bin/$(ARCH)
TARGET_SHARED := $(OUTPUTDIR)/$(SHARED_LIBNAME)
TARGET_STATIC := $(OUTPUTDIR)/$(STATIC_LIBNAME)

INCLUDEDIR = -I../include
CFLAGS = -Wall -Wextra -fPIC -ggdb -O0 -g3 -fvisibility=hidden \\
        -D${PROJECT_NAME_UPPER}_COMPILE -D_GNU_SOURCE $(INCLUDEDIR)

LIBDIR =
LIBS =

VPATH = ../include:.

C_FILES := $(wildcard *.c)
OBJS = $(C_FILES:.c=.o)

shared: outputdirs $(OBJS)
	$(CC) -shared -Wl,-soname,$(SONAME),--version-script,$(PREFIX).sym -o $(TARGET_SHARED) $(OBJS) $(LIBDIR) $(LIBS)

static: outputdirs $(OBJS)
	$(AR) -sr $(TARGET_STATIC) $(OBJS)

clean:
	rm -rf $(OBJS) $(TARGET_SHARED) $(TARGET_STATIC) *~ ../include/*~

dest_clean:
	rm -f $(USR_DIR)/$(LIBNAME)*

install:
	cp -f $(TARGET_SHARED) $(USR_DIR)
	rm -rf $(USR_DIR)/$(LIBNAME) $(USR_DIR)/$(SONAME)
	ln -s $(USR_DIR)/$(SHARED_LIBNAME) $(USR_DIR)/$(LIBNAME)
	ln -s $(USR_DIR)/$(SHARED_LIBNAME) $(USR_DIR)/$(SONAME)

outputdirs: $(OUTPUTDIR)
$(OUTPUTDIR):
	mkdir -p $(OUTPUTDIR)

'''

LIB_MAKEFILE_PACKAGE = '''.PHONY: shared static clean dest_clean install outputdirs package_version

CC = $COMPILER
AR = ar

ARCH_TEST := $(shell uname -m)

ifeq ($(ARCH_TEST), x86_64)
    ARCH = x86_64
else
    ARCH = i686
endif

MAJOR_VERSION = $(shell command grep MAJOR_VERSION ../../package_version.h | awk '{print $$$3}')
MINOR_VERSION = $(shell command  grep MINOR_VERSION ../../package_version.h | awk '{print $$$3}')
RELEASE = $(shell command grep RELEASE ../../package_version.h | awk '{print $$$3}')

USR_DIR = /usr/local/lib
PREFIX = ${PROJECT_NAME}
LIBNAME = $(PREFIX).so
SONAME = $(LIBNAME)
SHARED_LIBNAME = $(LIBNAME).$(MAJOR_VERSION).$(MINOR_VERSION).$(RELEASE)
STATIC_LIBNAME = $(PREFIX).a

OUTPUTDIR = ../bin/$(ARCH)
TARGET_SHARED = $(OUTPUTDIR)/$(SHARED_LIBNAME)
TARGET_STATIC = $(OUTPUTDIR)/$(STATIC_LIBNAME)

INCLUDEDIR = -I../include
CFLAGS = -Wall -Wextra -fPIC -ggdb -O0 -g3 -fvisibility=hidden \\
        -D${PROJECT_NAME_UPPER}_COMPILE -D_GNU_SOURCE $(INCLUDEDIR)

LIBDIR =
LIBS =

VPATH = ../include:.

C_FILES := $(wildcard *.c)
OBJS = $(C_FILES:.c=.o)

shared: outputdirs package_version $(OBJS)
	$(CC) -shared -Wl,-soname,$(SONAME),--version-script,$(PREFIX).sym -o $(TARGET_SHARED) $(OBJS) $(LIBDIR) $(LIBS)

static: outputdirs package_version $(OBJS)
	$(AR) -sr $(TARGET_STATIC) $(OBJS)

clean:
	rm -rf $(OBJS) $(TARGET_SHARED) $(TARGET_STATIC) *~ ../include/*~

dest_clean:
	rm -f $(USR_DIR)/$(LIBNAME)*

install:
	cp -f $(TARGET_SHARED) $(USR_DIR)
	rm -rf $(USR_DIR)/$(LIBNAME) $(USR_DIR)/$(SONAME)
	ln -s $(USR_DIR)/$(SHARED_LIBNAME) $(USR_DIR)/$(LIBNAME)
	ln -s $(USR_DIR)/$(SHARED_LIBNAME) $(USR_DIR)/$(SONAME)

outputdirs: $(OUTPUTDIR)
$(OUTPUTDIR):
	mkdir -p $(OUTPUTDIR)

PACKAGE_CONF=package/package.conf
PACKAGE_VERSION_NAME=package_version
PACKAGE_VERSION=../../$(PACKAGE_VERSION_NAME).h
package_version: $(PACKAGE_VERSION)
$(PACKAGE_VERSION):
	$(shell (cd ../../ && source-tpl -t header -n $(PACKAGE_VERSION_NAME) \\
	    -c "^#define MAJOR_VERSION	`cfget -C $(PACKAGE_CONF) version/major`\\
	    ^#define MINOR_VERSION	`cfget -C $(PACKAGE_CONF) version/minor` \\
	    ^#define RELEASE		`cfget -C $(PACKAGE_CONF) version/release` \\
	    ^#define BETA		`cfget -C $(PACKAGE_CONF) version/beta`^"))

'''

LIBSYM = '''${PROJECT_NAME_UPPER}_0.1 {
    global:
        *;
    local:
        *;
};
'''

DEF_HEADER = '''
#define MAJOR_VERSION   0
#define MINOR_VERSION   1
#define RELEASE         1
#define BETA            true
'''

PACKAGE_DEF_HEADER = '''
/*
 * Package version: major, minor and release.
 */
#include "../../package_version.h"

#define BUILD           0
'''

MAIN_HEADER = '''
#include "${PROJECT_NAME}_def.h"
#include "${PROJECT_NAME}_struct.h"
#include "${PROJECT_NAME}_prt.h"
'''

LIB_HEADER = '''
#ifdef ${PROJECT_NAME_UPPER}_COMPILE
# define MAJOR_VERSION  0
# define MINOR_VERSION  1
# define RELEASE        1
#endif
'''

class CTemplate(base.BaseTemplate):
    def __init__(self, args, project_vars):
        super(CTemplate, self).__init__()
        self._args = args
        self._project_vars = project_vars

        # Store the real info from each source/header which will be created.
        self._files = FileTemplate.FileTemplateInfo()
        self._prepare_project_files()


    def _library_header(self, _filename):
        """
        Decides which file content will be used for a library header.
        """
        content = ''

        if self._args.package is False:
            content = Template(LIB_HEADER).safe_substitute(self._project_vars)
        else:
            content = Template(PACKAGE_DEF_HEADER)\
                            .safe_substitute(self._project_vars)

        return content


    def _application_header(self, filename):
        """
        Decides which file content will be used for an application header or a
        single header file.
        """
        content = ''

        if '_def' in filename:
            if self._args.package is True or package.is_dir():
                content = \
                    Template(PACKAGE_DEF_HEADER)\
                        .safe_substitute(self._project_vars)
            else:
                content = \
                    Template(DEF_HEADER).safe_substitute(self._project_vars)
        else:
            if self._args.content is not None:
                content = Template(self._args.content.replace('^', '\n'))\
                                .safe_substitute(self._project_vars)
            elif self._args.project_name == filename and \
                    self._args.project_type not in (base.PTYPE_SOURCE, \
                                                    base.PTYPE_HEADER):
                content = Template(MAIN_HEADER)\
                                .safe_substitute(self._project_vars)

        return content


    def _get_header_content(self, filename):
        """
        Returns a content to insert in header files. Internally, we check if
        a content is needed to be inserted together.
        """
        content = ''
        upper_filename = filename.replace('.', '_').replace('-', '_').upper()

        # Do we have any content to add into the file?
        get_header = {
            base.PTYPE_LIBRARY: self._library_header,
        }.get(self._args.project_type, self._application_header)

        content = get_header(filename)

        return '''#ifndef _%s_H
#define _%s_H     1
%s
#endif

''' % (upper_filename, upper_filename, content)


    def _add_file(self, filename, path='src', extension=SOURCE_EXTENSION,
                  comment=True):
        """
        Adds a file into the internal FileTemplate object. Here we also
        """
        content = None

        if extension == HEADER_EXTENSION:
            content = self._get_header_content(filename)
        else:
            if self._args.content is not None:
                content = Template(self._args.content.replace('^', '\n'))\
                                .safe_substitute(self._project_vars)

        self._files.add(filename, path, content)
        self._files.set_property(filename, 'extension', extension)
        self._files.set_property(filename, 'comment', comment)


    def _add_makefile(self):
        """
        Adds a Makefile to the project.
        """
        if self._args.package is True:
            mcontent = {
                base.PTYPE_APPLICATION: APP_MAKEFILE_PACKAGE,
                base.PTYPE_LIBRARY: LIB_MAKEFILE_PACKAGE
            }.get(self._args.project_type)
        else:
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

            self._add_file(self._args.project_name, '', extension, False)
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

            self._files.add(self._args.prefix + app_name + '.sym', 'src',
                            Template(LIBSYM).safe_substitute(self._project_vars))

        for filename in self._args.sources:
            self._add_file(filename)

        for filename in self._args.headers:
            self._add_file(filename, 'include', HEADER_EXTENSION)


    def _create_single_file(self, filename, root_dir):
        """
        Creates the file with name @filename. Every one will have a comment
        block in the beginning.
        """
        file_data = self._files.properties(filename)
        extension = file_data.get('extension', '')
        path = file_data.get('path')
        content = file_data.get('data')
        comment = file_data.get('comment', False)
        comment_cnt = None

        if extension not in filename:
            filename += extension

        if len(root_dir):
            pathname = root_dir + '/' + path + '/' + filename
        else:
            pathname = filename

        # If we don't have a comment as True previously adjusted we assume that
        # is a single file creation and we use other comment block.
        if comment is True:
            comment_cnt = Template(COMMENT).safe_substitute(self._project_vars)
        else:
            comment_cnt = Template(COMMENT_SINGLE)\
                            .safe_substitute(self._project_vars)

        with open(pathname, 'w') as out_fd:
            if comment_cnt:
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

        for directory in subdirs:
            try:
                os.mkdir(root_dirname + '/' + directory)
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



