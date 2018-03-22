// The C file source type implementation.
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
	"fmt"
	"os"
	"text/template"

	"source-template/pkg/base"
)

type SourceFile struct {
	filename string
	content  string
	options  base.FileOptions
	ContentData
}

func (s SourceFile) Header(file *os.File) {
	var cnt string

	switch s.options.Language {
	case base.CLanguage:
		// if we're creating a project, probably will have an include directive here
		if s.options.ProjectType == base.LibraryProject {
			cnt = fmt.Sprintf("\n#include \"lib%[1]s.h\"\n", s.options.ProjectName)
		} else if s.options.ProjectType == base.XantePluginProject {
			cnt = fmt.Sprintf("\n#include \"plugin.h\"\n")
		} else {
			// XXX: Do we need this include in a single source file?
			cnt = fmt.Sprintf("\n#include \"%[1]s.h\"\n", s.options.ProjectName)
		}

	case base.GoLanguage:
		cnt = fmt.Sprintf("\npackage main\n")
	}

	file.WriteString(cnt)
}

// TODO: Add go support
func (s SourceFile) HeaderComment(file *os.File) {
	tpl, err := SourceHeader(s.options.Language)

	if err != nil {
		return
	}

	tpl.Execute(file, s.ContentData)
}

func (s SourceFile) Footer(file *os.File) {
	// nothing here
}

func (s SourceFile) Content(file *os.File) {
	tmpTpl := template.New("source")
	tpl, err := tmpTpl.Parse(s.content)

	if err != nil {
		return
	}

	tpl.Execute(file, s.ContentData)
}

const mainContent = `
static void usage(void)
{
    printf("Usage: %s [OPTIONS]\n", APP_NAME);
    printf("A brief description.\n\n");
    printf("Options:\n\n");
    printf("  -h\tShows this help screen.\n");
    printf("  -v\tShows current {{.ProjectName}} version.\n");
    printf("\n");
}

static void version(void)
{
    printf("%s - Version %d.%d.%d %s\n", APP_NAME, MAJOR_VERSION, MINOR_VERSION,
           RELEASE, (BETA == true) ? "beta" : "");
}

int main(int argc, char **argv)
{
	const char *opt = "hv\0";
	int option;

	do {
		option = getopt(argc, argv, opt);

		switch (option) {
			case 'h':
				usage();
				return 1;

			case 'v':
				version();
				return 1;

			case '?':
				return -1;
		}
	} while (option != -1);

	return 0;
}
`

func pluginContent(options base.FileOptions) string {
	if options.ProjectOptions.Language == base.CLanguage {
		return `
/*
 *
 * Plugin information
 *
 */
CL_PLUGIN_SET_INFO(
    "{{.ProjectName}}",
    "0.1.1",
    "{{.Author}}",
    "description"
)

/*
 *
 * Startup and shutdown
 *
 */

CL_PLUGIN_INIT()
{
    return 0;
}

CL_PLUGIN_UNINIT()
{
}

/*
 *
 * Main libxante events
 *
 */

CL_PLUGIN_FUNCTION(int, xapl_init)
{
    xante_event_arg_t *xante_args;

	cl_plugin_argument_pointer(args, "xpp_args", (void **)&xante_args);
    return 0;
}

CL_PLUGIN_FUNCTION(void, xapl_uninit)
{
    xante_event_arg_t *xante_args;

	cl_plugin_argument_pointer(args, "xpp_args", (void **)&xante_args);
}

CL_PLUGIN_FUNCTION(void, xapl_config_load)
{
    xante_event_arg_t *xante_args;

	cl_plugin_argument_pointer(args, "xpp_args", (void **)&xante_args);
}

CL_PLUGIN_FUNCTION(void, xapl_config_unload)
{
    xante_event_arg_t *xante_args;

	cl_plugin_argument_pointer(args, "xpp_args", (void **)&xante_args);
}

CL_PLUGIN_FUNCTION(int, xapl_changes_saved)
{
    xante_event_arg_t *xante_args;

	cl_plugin_argument_pointer(args, "xpp_args", (void **)&xante_args);
    return 0;
}
`
	} else {
		return `
import "C"
import (
	"unsafe"

	"collections/pkg/collections"
	"xante/pkg/xante"
)

//export plugin_name
func plugin_name() *C.char {
	return C.CString("{{.ProjectName}}")
}

//export plugin_version
func plugin_version() *C.char {
	return C.CString("0.1.1")
}

//export plugin_author
func plugin_author() *C.char {
	return C.CString("{{.Author}}")
}

//export plugin_description
func plugin_description() *C.char {
	return C.CString("description")
}

//
// Startup and shutdown
//

//export plugin_init
func plugin_init() int {
	return 0
}

//export plugin_uninit
func plugin_uninit() {
}

//
// Libxante main events
//

//export xapl_init
func xapl_init(args unsafe.Pointer) int {
	return 0
}

//export xapl_uninit
func xapl_uninit(args unsafe.Pointer) {
}

//export xapl_config_load
func xapl_config_load(args unsafe.Pointer) {
}

//export xapl_config_unload
func xapl_config_unload(args unsafe.Pointer) {
}

//export xapl_changes_saved
func xapl_changes_saved(args unsafe.Pointer) int {
	return 0
}

func main() {
	//
	// We need the main function otherwise the ELF shared object created will be
	// in the wrong format, as an AR (archive) file and not an ELF shared object.
	//
}
`
	}
}

func NewSource(options base.FileOptions) base.FileTemplate {
	var content string
	bname, _ := extractFilename(options.Name, options.ProjectType)
	contentData := GetContentData(options)

	// here we build what will be the file content based on its name (basename)
	if bname == "main" {
		content = mainContent
	} else if bname == "error" {
		content = errorContent(Source, options)
	} else if bname == "plugin" {
		content = pluginContent(options)
	}

	return &SourceFile{
		options:     options,
		filename:    bname,
		content:     content,
		ContentData: contentData,
	}
}
