//
// Copyright (C) 2017 Rodrigo Freitas
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, write to the Free Software Foundation, Inc.,
// 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
//
package templates

import (
	"os"
	"text/template"

	"source-template/pkg/base"
)

const pluginScriptContent = `
jerminus -j {{.ProjectName}}.jtf -N
`

const packageBuildScriptContent = `
arch=""
mode="debug"
package="{{.ProjectName}}"

usage()
{
    echo "Usage: build-package.sh [OPTIONS]"
    echo "Script to build the current project as a debian file."
    echo
    echo "Options:"
    echo -e " -h\tShows this help screen."
    echo -e " -R\tCompiles the application in release mode (debug is default)."
    echo
}

validate_arch()
{
    if [ "$arch" != "386" -a "$arch" != "amd64" ]; then
        echo -1
    else
        echo 0
    fi
}

rust_compile()
{
    if [ "$mode" = "release" ]; then
        (cd ../$package && cargo build --release || exit -1)
    else
        (cd ../$package && cargo build || exit -1)
    fi

    if [ $? != 0 ]; then
        return -1
    fi

    return 0
}

go_compile()
{
    (cd ../$package/cmd/$package && GOARCH=$arch go build || exit -1)

    if [ $? != 0 ]; then
        return -1
    fi

    return 0
}

c_compile()
{
    if [ ! -d ../$package/build ]; then
        mkdir ../$package/build
        (cd ../$package/build && cmake ..)
    fi

    (cd ../$package/build && make || exit -1)

    if [ $? != 0 ]; then
        return -1
    fi

    return 0
}

compile()
{
    echo "Compiling..."

    if [ -e ../$package/CMakeLists.txt ]; then
        c_compile
    elif [ -e ../$package/Cargo.toml ]; then
        rust_compile
    else
        go_compile
    fi
}

package_version()
{
    echo "Get package version here"
}

package_release()
{
    echo "Get package release here"
}

copy_package_core_files()
{
    echo "Copy the package core files to the package structure"
}

build_package()
{
    local tmpdir="$package-release"
    local version=$(package_version)
    local release=$(package_release)
    local filename=$package-$version-$release-$arch.deb
    local depends=""

    echo "Copying internal package files..."
    mkdir -p $tmpdir/{opt/$package,DEBIAN,etc/systemd/system}
    copy_package_core_files

    # Copy package and misc files
    cp default/p* $tmpdir/DEBIAN
    cp ../misc/*.service $tmpdir/etc/systemd/system

    cat << CONTROL >> $tmpdir/DEBIAN/control
Package: $package
Priority: optional
Version: $version-$release
Architecture: $arch
Depends: $depends
Maintainer: {{.Author}}
Description:
CONTROL

    echo "Building package $filename"
    fakeroot dpkg-deb -Zgzip -b $tmpdir $filename

    rm -rf $tmpdir
}

while getopts ha:R: opts; do
    case $opts in
        h)
            usage
            exit 1
            ;;

        a)
            arch=$OPTARG
            ;;

        R)
            mode=$OPTARG
            ;;

        ?)
            exit -1
            ;;
    esac
done

if [ -z "$arch" -o $(validate_arch) != 0 ]; then
    echo "Unsupported '$arch' architecture!"
    exit -1
fi

# compile
compile
ret=$?

if [ $ret != 0 ]; then
    exit -1
fi

# build the package
build_package

`

type BashFile struct {
	content string
	base.FileOptions
	ContentData
}

func (s BashFile) Header(file *os.File) {
}

func (s BashFile) HeaderComment(file *os.File) {
	file.WriteString("#!/bin/bash\n")
	tpl, err := BashSourceHeader()

	if err != nil {
		return
	}

	tpl.Execute(file, s.ContentData)
}

func (s BashFile) Footer(file *os.File) {
	file.WriteString("\nexit 0\n")
}

func (s BashFile) Content(file *os.File) {
	tmpTpl := template.New("script")
	tpl, err := tmpTpl.Parse(s.content)

	if err != nil {
		return
	}

	tpl.Execute(file, s.ContentData)
}

func NewBash(options base.FileOptions) base.FileTemplate {
	var content string
	bname, _ := extractFilename(options.Name, options.ProjectType)

	if options.ProjectType == base.XantePluginProject {
		if bname == options.ProjectName {
			content = pluginScriptContent
		}
	}

	if options.PackageProject {
		if bname == "build-package" {
			content = packageBuildScriptContent
		}
	}

	return &BashFile{
		FileOptions: options,
		content:     content,
		ContentData: GetContentData(options),
	}
}
