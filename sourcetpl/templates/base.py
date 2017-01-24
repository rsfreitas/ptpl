
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
Template base abstract class.
"""

from abc import ABCMeta, abstractmethod

# Supported projects
PTYPE_SOURCE = 'source'
PTYPE_HEADER = 'header'
PTYPE_APPLICATION = 'application'
PTYPE_LIBRARY = 'library'
PTYPE_LIBCOLLECTION_APP = 'libcollection-app'
PTYPE_LIBCOLLECTION_C_PLUGIN = 'libcollection-c-plugin'
PTYPE_LIBCOLLECTION_PY_PLUGIN = 'libcollection-py-plugin'
PTYPE_LIBCOLLECTION_JAVA_PLUGIN = 'libcollection-java-plugin'

class BaseTemplate(object):
    """
    Abstract class for all supported languages
    """
    __metaclass__ = ABCMeta

    @abstractmethod
    def create(self):
        """
        Method to create the required project.
        """
        pass


    @abstractmethod
    def info(self):
        """
        Method to show all info of the project in a user friendly format.
        """
        pass



