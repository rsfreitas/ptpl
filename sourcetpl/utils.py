
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
Utility functions.
"""

import re
import commands
import collections

# Supported projects
PTYPE_SOURCE = 'source'
PTYPE_HEADER = 'header'
PTYPE_APPLICATION = 'application'
PTYPE_LIBRARY = 'library'
PTYPE_LIBCOLLECTION_APP = 'libcollection-app'
PTYPE_LIBCOLLECTION_C_PLUGIN = 'libcollection-c-plugin'
PTYPE_LIBCOLLECTION_PY_PLUGIN = 'libcollection-py-plugin'
PTYPE_LIBCOLLECTION_JAVA_PLUGIN = 'libcollection-java-plugin'

# Supported programming languages
C_LANGUAGE = 'C'
PYTHON_LANGUAGE = 'Python'

def git_author_name():
    """
    Gets the author name from a git environment.

    :return Returns the author name.
    """
    return commands.getoutput('git config user.name')



def multiple_split(entry):
    """
    Splits a string by any of the delimiters.

    :return Returns a list with the splitted strings.
    """
    if entry is None or len(entry) == 0:
        return list()

    return re.split(';|,|\|| ', entry)



def supported_projects():
    """
    Gets all supported project formats.

    :return Returns a dictionary with all supported project as keys and a brief
            description of each one.
    """
    return {
        PTYPE_SOURCE : 'Indicates the creation of a single source file.',
        PTYPE_HEADER : 'Indicates the creation of a single header file.',
        PTYPE_APPLICATION :
            '''Indicates the creation of a directory with the following
            \t\tstructure: $name/{include,src}, containing template files for
            \t\ta single application (with a main function).''',

        PTYPE_LIBRARY:
            '''Indicates the creation of a directory to hold a library
            \t\tdevelopment project, with a specific Makefile.''',

        PTYPE_LIBCOLLECTION_APP :
            '''Indicates the creation of an application using
            \t\tlibcollections as its base.''',

        PTYPE_LIBCOLLECTION_C_PLUGIN :
            '''Indicates the creation of a libcollections' plugin (C)''',

        PTYPE_LIBCOLLECTION_PY_PLUGIN :
            '''Indicates the creation of a libcollections' plugin (Python)''',

        PTYPE_LIBCOLLECTION_JAVA_PLUGIN :
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
    return [C_LANGUAGE, PYTHON_LANGUAGE]



