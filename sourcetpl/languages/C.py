
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
Specific data for C language.
"""

HEADER_EXTENSION = '.h'
SOURCE_EXTENSION = '.c'

HEAD = '''
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

GITIGNORE = '''# Object files
*.o
*.ko
*.obj
*.elf

# Precompiled Headers
*.gch
*.pch

# Libraries
*.lib
*.a
*.la
*.lo

# Shared objects (inc. Windows DLLs)
*.dll
*.so
*.so.*
*.dylib

# Executables
*.exe
*.out
*.app
*.i*86
*.x86_64
*.hex

# Debug files
*.dSYM/
*.log

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


