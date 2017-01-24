
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
Functions to control a package creation.
"""

import os

from string import Template

PREFIX = 'package'

build_package = '''
'''

clean_package = '''#!/bin/bash

package_dir=../../

# Remove older versions
rm -rf *.pkg

for arq in $package_dir*/src; do
    echo "Cleaning source directory: <$arq>\n"
    (cd $arq && make clean)
done

exit 0
'''

deb_scripts = '''#!/bin/bash

exit 0
'''

cron = '''SHELL=/bin/sh
PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin

*/1 * * * *    root    /etc/init.d/$PROJECT_NAME.sh status || /etc/init.d/$PROJECT_NAME.sh start
'''

initd = '''#!/bin/sh

. /lib/lsb/init-functions

case "$1" in
    start)
        log_begin_msg "Starting $PROJECT_NAME: "

        if start-stop-daemon --start --quiet --exec /usr/local/bin/$PROJECT_NAME; then
            log_end_msg 0
        else
            log_end_msg 1
        fi
        ;;

    stop)
        log_begin_msg "Shutting down $PROJECT_NAME: "

        if start-stop-daemon --stop --quiet --exec /usr/local/bin/$PROJECT_NAME; then
            log_end_msg 0
        else
            log_end_msg 1
        fi
        ;;

    status)
        if [ -s /var/run/$PROJECT_NAME.pid ]; then
            if kill -0 `cat /var/run/$PROJECT_NAME.pid` 2>/dev/null; then
                log_success_msg "$PROJECT_NAME esta sendo executado"
                exit 0
            else
                log_failure_msg "/var/run/$PROJECT_NAME.pid exists but $PROJECT_NAME is not running"
                exit 1
            fi
        else
            log_success_msg "$PROJECT_NAME is not running"
            exit 3
        fi
        ;;

    restart)
        $0 stop
        sleep 5
        $0 start
        ;;

    reload)
        log_begin_msg "Restarting $PROJECT_NAME: "
        start-stop-daemon --stop --signal 10 --exec /usr/local/bin/$PROJECT_NAME || log_end_msg 1
        log_end_msg 0
        ;;

    *)
        log_begin_msg "Usage: %s (start|stop|status|restart|reload)" "$0"
        exit 1
esac

exit 0
'''

class Package(object):
    def __init__(self, args, project_vars):
        self._args = args
        self._project_vars = project_vars
        self._root_dir = PREFIX + '-' + \
                self._args.prefix + self._args.project_name.replace('_', '-')


    def current_dir(self):
        """
        Returns package current root directory.
        """
        return self._root_dir


    def _create_directories(self):
        """
        Creates the package structure directories.
        """
        subdirs = [
            ['package', ['debian', 'mount', 'misc']]
        ]

        os.mkdir(self._root_dir)

        for directory in subdirs:
            os.mkdir(self._root_dir + '/' + directory[0])

            for subdir in directory[1]:
                os.mkdir(self._root_dir + '/' + directory[0] + '/' + subdir)


    def _create_debian_scripts(self):
        """
        Creates the debian package required scripts.
        """
        output = Template(deb_scripts).safe_substitute(self._project_vars)

        for filename in ['postinst', 'preinst', 'postrm', 'prerm']:
            pathname = self._root_dir + '/package/debian/' + filename

            with open(pathname, 'w') as out_fd:
                out_fd.write(output)

            os.system('chmod +x %s' % pathname)


    def _create_package_misc_scripts(self):
        """
        Creates a bunch of scripts for the package.
        """
        # Creates the cron script
        suffixes = ['_cron', '_initd']
        output = [
            Template(cron).safe_substitute(self._project_vars),
            Template(initd).safe_substitute(self._project_vars)
        ]

        path = self._root_dir + '/package/misc/' + \
                self._args.project_name.replace('-', '_')

        for data, suffix in zip(output, suffixes):
            with open(path + suffix, 'w') as out_fd:
                out_fd.write(data)

        os.system('chmod +x %s' % path + '_initd')


    def _create_package_handle_scripts(self):
        """
        Create scripts to build the package.
        """
        names = ['buid-package', 'clean-package']
        output = [
            Template(build_package).safe_substitute(self._project_vars),
            Template(clean_package).safe_substitute(self._project_vars)
        ]

        path = self._root_dir + '/package/mount/'

        for data, name in zip(output, names):
            with open(path + name, 'w') as out_fd:
                out_fd.write(data)

            os.system('chmod +x %s' % path + name)


    def create(self):
        self._create_directories()
        self._create_debian_scripts()
        self._create_package_misc_scripts()
        self._create_package_handle_scripts()



