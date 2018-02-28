package templates

import (
	//	"fmt"
	"os"
	//	"path/filepath"
	//	"strings"
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
    target_link_libraries(${PROJECT_NAME} collections sqlite3 crypto dialog ncursesw)
    set_target_properties(${PROJECT_NAME} PROPERTIES
                          LINK_FLAGS "-Wl,-soname,${PROJECT_NAME}.so,--version-script,${VERSION_SCRIPT}")

    set_target_properties(${PROJECT_NAME} PROPERTIES
                          SUFFIX .so.${MAJOR_VERSION}.${MINOR_VERSION}.${RELEASE})
else(SHARED)
    add_library(${PROJECT_NAME} STATIC ${SOURCE})
endif(SHARED)

install(TARGETS ${PROJECT_NAME} DESTINATION ${DESTINATION_BIN_DIR})
install(FILES ${LIBRARY_HEADER} DESTINATION ${DESTINATION_HEADER_DIR}/${PROJECT_NAME})
install(DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/include/api DESTINATION ${DESTINATION_HEADER_DIR}/${PROJECT_NAME})

# If we're dealing with a shared version, we create its symbolic link
if(SHARED)
    set(RELEASE_NAME lib${PROJECT_NAME}.so.${MAJOR_VERSION}.${MINOR_VERSION}.${RELEASE})
    set(LINK_NAME lib${PROJECT_NAME}.so)
    install(CODE "execute_process(COMMAND ln -sf ${DESTINATION_BIN_DIR}/${RELEASE_NAME} ${DESTINATION_BIN_DIR}/${LINK_NAME})")
endif(SHARED)
`

type Makefile struct {
	Options base.FileOptions
	ContentData
}

func (m Makefile) Header(file *os.File) {
}

func (m Makefile) HeaderComment(file *os.File) {
}

func (m Makefile) Footer(file *os.File) {
}

func (m Makefile) Content(file *os.File) {
	var content string
	tpl := template.New("cmake")

	if m.Options.ProjectType == base.LibraryProject {
		content = libContent
	}

	tpl, err := tpl.Parse(content)

	if err != nil {
		return
	}

	tpl.Execute(file, m.ContentData)
}

func NewMakefile(options base.FileOptions) base.FileTemplate {
	return &Makefile{
		Options:     options,
		ContentData: GetContentData(options),
	}
}
