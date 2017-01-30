
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

from . import languages
from templates import FileTemplate

PREFIX = 'package'

BUILD_PACKAGE = '''#!/bin/bash

#
# Description:
#
# Author: $FULL_AUTHOR_NAME
# Created at: $DATE
# Project: $PROJECT_NAME
#
# Copyright (c) $YEAR All rights reserved
#

package_conf="../package.conf"
package_tmp_dir=tmpbuild
arch=i686

# Compile all internal applications
compile_applications()
{
    local modules=`cfget -C $package_conf package/modules`

    for app in ${modules[@]}; do
        echo "(cd ../../$app/src && make clean && make)"
    done
}

# Prepare the current environment to build the package
prepare_to_build()
{
    rm -f *.deb
    mkdir -p $package_tmp_dir/{DEBIAN,usr/bin,etc/{cron.d,init.d}}
}

# Copy all necessary files
copy_necessary_files()
{
    local modules=`cfget -C $package_conf package/modules`

    prepare_to_build

    # Copy all binaries
    for app in ${modules[@]}; do
        cp -f ../../$app/bin/$arch/* $package_tmp_dir/usr/bin || \\
            echo "Error while copying binary from module '$app'."
    done

    # Copy misc scripts
    local dest_cron=`basename ../misc/*_cron _cron`
    cp -f ../misc/*_cron $package_tmp_dir/etc/cron.d/$dest_cron || \\
        echo "Error copying to file '$dest_cron'."

    local dest_init=`basename ../misc/*_initd _initd`
    cp -f ../misc/*_initd $package_tmp_dir/etc/init.d/$dest_init || \\
        echo "Error copying to file '$dest_init'."

    # Copy all debian scripts
    for script in postinst postrm preinst prerm; do
        cp -f ../debian/$script $package_tmp_dir/DEBIAN || \\
            echo "Error copying file '$script'."
    done
}

clear_temporary_files()
{
    rm -rf $package_tmp_dir
}

# Build the package
build_package()
{
    # Build all package info, such as version, revision, etc
    local package=`cfget -C $package_conf package/name`
    local major=`cfget -C $package_conf version/major`
    local minor=`cfget -C $package_conf version/minor`
    local release=`cfget -C $package_conf version/release`
    local beta_release=`cfget -C $package_conf version/beta`
    local version=$major.$minor.$release

    local depends=''
    local maintainer=''
    local description=''

    # Create CONTROL file
    cat << CONTROL >> $package_tmp_dir/DEBIAN/control
Package: $package
Priority: optional
Version: $version
Architecture: $arch
Depends: $depends
Maintainer: $maintainer
Description: $description
CONTROL

    # Build the package
    if [ $beta_release == true ]; then
        beta="-Beta"
    else
        beta=""
    fi

    deb_filename=$package-$version.deb
    fakeroot dpkg-deb -Zgzip -b $package_tmp_dir $deb_filename

    clear_temporary_files
}

usage()
{
    echo "Usage: ./build-package [OPTIONS]"
    echo
    echo "Options"
    echo -e " -h\\t\\tShows this help screen."
    echo -e " -a [arch]\\tDefines the package destination architecture."
    echo

    exit 1
}

while getopts "ha:" OPTION; do
    case $OPTION in
        h)
            usage
            ;;
        a)
            arch=$OPTARG
            ;;
        \?)
            exit 1
            ;;
    esac
done

compile_applications
copy_necessary_files
build_package

exit 0
'''

CLEAN_PACKAGE = '''#!/bin/bash

package_dir=../../

# Remove older versions
rm -rf *.deb

for arq in $package_dir*/src; do
    echo "Cleaning source directory: $arq"
    (cd $arq && make clean)
done

exit 0
'''

DEB_SCRIPTS = '''#!/bin/bash

exit 0
'''

CRON = '''SHELL=/bin/sh
PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin

*/1 * * * *    root    /etc/init.d/$PROJECT_NAME.sh status || /etc/init.d/$PROJECT_NAME.sh start
'''

INITD = '''#!/bin/sh

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

PACKAGE_CONF = '''# Main package informations.
# The package modules must be separated by spaces.
[package]
name=$PROJECT_NAME
modules=$PROJECT_NAME

# Package version: major.minor.release
# Example: 0.1.1
[version]
major=0
minor=1
release=1
beta=true
'''

class Package(object):
    def __init__(self, args, project_vars):
        self._args = args
        self._project_vars = project_vars
        self._root_dir = PREFIX + '-' + \
                self._args.prefix + self._args.project_name.replace('_', '-')

        self._files = FileTemplate.FileTemplateInfo()
        self._prepare_package_files()


    def current_dir(self):
        """
        Returns package current root directory.
        """
        return self._root_dir


    def _prepare_package_files(self):
        """
        Adds all required files to a package.
        """
        prefix = self._args.project_name.replace('-', '_')
        build_package = {
            languages.C_LANGUAGE: \
                Template(BUILD_PACKAGE).safe_substitute(self._project_vars)
        }.get(self._args.language)

        files = [
            # debian scripts
            ('postinst', True, 'debian',
                Template(DEB_SCRIPTS).safe_substitute(self._project_vars)),

            ('postrm', True, 'debian',
                Template(DEB_SCRIPTS).safe_substitute(self._project_vars)),

            ('preinst', True, 'debian',
                Template(DEB_SCRIPTS).safe_substitute(self._project_vars)),

            ('prerm', True, 'debian',
                Template(DEB_SCRIPTS).safe_substitute(self._project_vars)),

            # build-package
            ('build-package', True, 'mount', build_package),

            # clean-package
            ('clean-package', True, 'mount',
                Template(CLEAN_PACKAGE).safe_substitute(self._project_vars)),

            # cron
            (prefix + '_cron', False, 'misc',
                Template(CRON).safe_substitute(self._project_vars)),

            # initd
            (prefix + '_initd', True, 'misc',
                Template(INITD).safe_substitute(self._project_vars)),

            # package.conf
            ('package.conf', False, '',
                Template(PACKAGE_CONF).safe_substitute(self._project_vars))
        ]

        for script in files:
            self._files.add(script[0], script[2], data=script[3],
                            executable=script[1])


    def _create_directories(self):
        """
        Create the package required directories.
        """
        subdirs = [
            ['package', ['debian', 'mount', 'misc']]
        ]

        os.mkdir(self._root_dir)

        for directory in subdirs:
            os.mkdir(self._root_dir + '/' + directory[0])

            for subdir in directory[1]:
                os.mkdir(self._root_dir + '/' + directory[0] + '/' + subdir)


    def _create_files(self):
        for filename in self._files.filenames():
            file_data = self._files.properties(filename)

            pathname = self._root_dir + '/package/' + file_data.get('path') \
                    + '/' + filename

            with open(pathname, 'w') as out_fd:
                out_fd.write(file_data.get('data'))

            if file_data.get('chmod') is True:
                os.system('chmod +x %s' % pathname)


    def create(self):
        self._create_directories()
        self._create_files()



