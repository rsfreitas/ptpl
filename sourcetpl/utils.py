
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



