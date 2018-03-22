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

const libContent = `cmake_minimum_required(VERSION 2.8)
project({{.ProjectName}})

# Options
option(DEBUG "Enable/Disable debug library" ON)
option(SHARED "Enable/Disable the shared library version" ON)

include_directories(include)
include_directories("include/api")
include_directories("include/internal")

if(CMAKE_C_COMPILER_VERSION VERSION_GREATER 5)
    add_definitions(-fgnu89-inline)
endif()

if(DEBUG)
    set(CMAKE_BUILD_TYPE Debug)
else(DEBUG)
    set(CMAKE_BUILD_TYPE Release)
endif(DEBUG)

add_definitions("-Wall -Wextra -fPIC")
add_definitions("-DLIB{{.ProjectNameUpper}}_COMPILE -D_GNU_SOURCE")

file(GLOB SOURCES "src/*.c")

set(SOURCE
    ${SOURCES})

set(VERSION_SCRIPT
    ${CMAKE_CURRENT_SOURCE_DIR}/misc/lib${PROJECT_NAME}.sym)

set(LIBRARY_HEADER
    ${CMAKE_CURRENT_SOURCE_DIR}/include/lib${PROJECT_NAME}.h)

execute_process(COMMAND grep MAJOR_VERSION ${LIBRARY_HEADER}
    COMMAND awk "{print $4}"
    COMMAND tr "\n" " "
    COMMAND sed "s/ //"
    OUTPUT_VARIABLE MAJOR_VERSION)

execute_process(COMMAND grep MINOR_VERSION ${LIBRARY_HEADER}
    COMMAND awk "{print $4}"
    COMMAND tr "\n" " "
    COMMAND sed "s/ //"
    OUTPUT_VARIABLE MINOR_VERSION)

execute_process(COMMAND grep RELEASE ${LIBRARY_HEADER}
    COMMAND awk "{print $4}"
    COMMAND tr "\n" " "
    COMMAND sed "s/ //"
    OUTPUT_VARIABLE RELEASE)

set(DESTINATION_BIN_DIR "/usr/local/lib")
set(DESTINATION_HEADER_DIR "/usr/local/include")

link_directories(${DESTINATION_BIN_DIR})

if(SHARED)
    add_library(${PROJECT_NAME} SHARED ${SOURCE})
    target_link_libraries(${PROJECT_NAME} collections)
    set(LIB_VERSION ${MAJOR_VERSION}.${MINOR_VERSION}.${RELEASE})
    set_target_properties(${PROJECT_NAME} PROPERTIES VERSION ${LIB_VERSION}
        SOVERSION ${MAJOR_VERSION})

    set_target_properties(${PROJECT_NAME} PROPERTIES
                          LINK_FLAGS "-Wl,--version-script,${VERSION_SCRIPT}")

    set_target_properties(${PROJECT_NAME} PROPERTIES
                          SUFFIX .so.${MAJOR_VERSION}.${MINOR_VERSION}.${RELEASE})
else(SHARED)
    add_library(${PROJECT_NAME} STATIC ${SOURCE})
endif(SHARED)

install(TARGETS ${PROJECT_NAME} DESTINATION ${DESTINATION_BIN_DIR})
install(FILES ${LIBRARY_HEADER} DESTINATION ${DESTINATION_HEADER_DIR}/${PROJECT_NAME})
install(DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/include/api DESTINATION ${DESTINATION_HEADER_DIR}/${PROJECT_NAME})
`

const appContent = `project({{.ProjectName}})
cmake_minimum_required(VERSION 2.8)

# Options
option(DEBUG "Enable/Disable debug version" ON)

include_directories(include)
include_directories("/usr/local/include")

if(CMAKE_C_COMPILER_VERSION VERSION_GREATER 5)
    add_definitions(-fgnu89-inline)
endif()

add_definitions("-Wall -Wextra -O0")

if(DEBUG)
    add_definitions("-ggdb")
endif(DEBUG)

file(GLOB SOURCES "src/*c")
add_executable(${PROJECT_NAME} ${SOURCES})

link_directories("/usr/local/lib")
target_link_libraries(${PROJECT_NAME} {{.LibcollectionsLinker}})
`

const pluginCMakeContent = `project({{.ProjectName}})
cmake_minimum_required(VERSION 2.8)

# Options
option(DEBUG "Enable/Disable debug version" ON)

include_directories(include)
include_directories("/usr/local/include")

if(CMAKE_C_COMPILER_VERSION VERSION_GREATER 5)
    add_definitions(-fgnu89-inline)
endif()

add_definitions("-Wall -Wextra -O0 -fPIC -fvisibility=hidden -D_GNU_SOURCE")

if(DEBUG)
    add_definitions("-ggdb -g3")
endif(DEBUG)

file(GLOB SOURCES "src/*c")

link_directories("/usr/local/lib")
add_library(${PROJECT_NAME} SHARED ${SOURCES})
target_link_libraries(${PROJECT_NAME} xante collections)
set_target_properties(${PROJECT_NAME} PROPERTIES
                      LINK_FLAGS "-Wl,-soname,${PROJECT_NAME}.so")

set_target_properties(${PROJECT_NAME} PROPERTIES SUFFIX .so)
set_target_properties(${PROJECT_NAME} PROPERTIES PREFIX "")
`

const goPluginMakefile = `
.PHONY: clean install purge

TARGET = {{.ProjectName}}.so

$(TARGET): plugin.go
	go build -o $(TARGET) -buildmode=c-shared plugin.go

clean:
	rm -f $(TARGET)

purge: clean $(TARGET)

install:
	cp -f $(TARGET) /usr/local/lib
`

type Makefile struct {
	Options base.FileOptions
	ContentData
}

func (m Makefile) Header(file *os.File) {
	// nothing here
}

func (m Makefile) HeaderComment(file *os.File) {
	// nothing here
}

func (m Makefile) Footer(file *os.File) {
	// nothing here
}

func (m Makefile) Content(file *os.File) {
	var content string
	tpl := template.New("cmake")

	if m.Options.ProjectType == base.LibraryProject {
		content = libContent
	} else if m.Options.ProjectType == base.XantePluginProject {
		if m.Options.Language == base.GoLanguage {
			content = goPluginMakefile
		} else {
			content = pluginCMakeContent
		}
	} else {
		content = appContent
	}

	tpl, err := tpl.Parse(content)

	if err != nil {
		return
	}

	tpl.Execute(file, m.ContentData)
}

func NewMakefile(options base.FileOptions) base.FileTemplate {
	contentData := GetContentData(options)

	if options.LibcollectionsFeatures {
		contentData.LibcollectionsLinker = "collections"
	}

	return &Makefile{
		Options:     options,
		ContentData: contentData,
	}
}
