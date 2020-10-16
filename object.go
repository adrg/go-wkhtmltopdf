package pdf

/*
#cgo LDFLAGS: -lwkhtmltox
#include <stdlib.h>
#include <wkhtmltox/pdf.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unsafe"
)

// ErrorAction defines actions to take in case of object load failure.
type ErrorAction string

// Error action values.
const (
	ActionAbort  ErrorAction = "abort"
	ActionIgnore ErrorAction = "ignore"
	ActionSkip   ErrorAction = "skip"
)

// TOC contains settings related to the table of contents of an object.
type TOC struct {
	// Specifies whether dotted lines should be used for the line of items
	// of the TOC.
	UseDottedLines bool

	// The title used for the table of contents.
	// E.g.: "Table of Contents".
	Title string

	// Specifies whether the TOC items should contain links to the content.
	GenerateForwardLinks bool

	// Specifies whether the content should contain links to the TOC.
	GenerateBackLinks bool

	// The indentation used for the TOC nesting levels.
	// E.g.: "1em".
	Indentation string

	// Scaling factor for each nesting level of the TOC.
	// E.g.: 1.
	FontScale float64
}

// Header contains settings related to the headers and footers of an object.
type Header struct {
	// The system font name to use for headers/footers.
	// E.g.: "Arial".
	Font string

	// The font size to use for headers/footers.
	// E.g.: 12.
	FontSize uint64

	// Content to print on each of the available regions of the header/footer.
	// Substitution variables that can be used in the content fields:
	//  - [page]       The number of the current page.
	//  - [frompage]   The number of the first page.
	//  - [topage]     The number of the last page.
	//  - [webpage]    The URL of the source page.
	//  - [section]    The name of the current section.
	//  - [subsection] The name of the current subsection.
	//  - [date]       The current date in system local format.
	//  - [isodate]    The current date in ISO 8601 extended format.
	//  - [time]       The current time in system local format.
	//  - [title]      The title of the of the current page object.
	//  - [doctitle]   The title of the output document.
	//  - [sitepage]   The number of the page in the currently converted site.
	//  - [sitepages]  The number of pages in the current site being converted.
	// e.g.: object.Footer.ContentRight = "[page]"
	ContentLeft   string
	ContentCenter string
	ContentRight  string

	// Specifies whether a line separator should be printed for headers/footers.
	DisplaySeparator bool

	// The amount of space between the header/footer and the content.
	// E.g.: 0.
	Spacing float64

	// Location of a user defined HTML document to be used as the header/footer.
	CustomLocation string
}

// ObjectOpts defines a set of options to be used in the conversion process.
type ObjectOpts struct {
	// Specifies the location of the HTML document. Can be a file path or a URL.
	Location string

	// Specifies whether external links in the HTML document should be converted
	// to external PDF links.
	UseExternalLinks bool

	// Specifies whether internal links in the HTML document should be converted
	// into PDF references.
	UseLocalLinks bool

	// Specifies whether HTML forms should be converted into PDF forms.
	ProduceForms bool

	// Specifies whether the sections from the HTML document are included in
	// outlines and TOCs.
	IncludeInOutline bool

	// Specifies whether the page count of the HTML document participates in
	// the counter used for tables of contents, headers and footers.
	CountPages bool

	// Contains settings for the TOC of the object.
	TOC TOC

	// Contains settings for the header of the object.
	Header Header

	// Contains settings for the footer of the object.
	Footer Header

	// The username to use when logging in to a website.
	Username string

	// The password to use when logging in to a website.
	Password string

	// The amount of milliseconds to wait after page load, before
	// executing JS scripts.
	// E.g.: 300.
	JavascriptDelay uint64

	// Specifies the `window.status` value to wait for, before
	// rendering the page.
	// E.g.: "ready".
	WindowStatus string

	// Zoom factor to use for the document content.
	// E.g.: 1.
	Zoom float64

	// Specifies whether local file access is blocked.
	BlockLocalFileAccess bool

	// Specifies whether slow JS scripts should be stopped.
	StopSlowScripts bool

	// Specifies a course of action when an HTML document fails to load.
	// E.g.: ActionAbort.
	ErrorAction ErrorAction

	// The name of a proxy to use when loading the HTML document.
	Proxy string

	// Specifies whether the background of the HTML document is preserved.
	PrintBackground bool

	// Specifies whether the images in the HTML document are loaded.
	LoadImages bool

	// Specifies whether Javascript should be executed.
	EnableJavascript bool

	// Specifies whether to use intelligent shrinkng in order to fit more
	// content on a page.
	UseSmartShrinking bool

	// The minimum font size allowed for rendering content.
	MinFontSize uint64

	// The text encoding to use if the HTML document does not specify one.
	// E.g.: "utf-8".
	DefaultEncoding string

	// Specifies whether the content should be rendered using the print media
	// type instead of the screen media type.
	UsePrintMediaType bool

	// The location of a user defined stylesheet to use when converting
	// the HTML document.
	UserStylesheetLocation string

	// Specifies whether NS plugins should be enabled.
	EnablePlugins bool
}

// NewObjectOpts returns a new instance of object options, configured
// using sensible defaults.
//
//   Defaults options:
//
//   UseExternalLinks:  true
//   UseLocalLinks:     true
//   IncludeInOutline:  true
//   CountPages:        true
//   JavascriptDelay:   300
//   Zoom:              1
//   StopSlowScripts:   true
//   ErrorAction:       ActionAbort
//   PrintBackground:   true
//   LoadImages:        true
//   EnableJavascript:  true
//   UseSmartShrinking: true
//   DefaultEncoding:   "utf-8"
//   TOC:
//   	UseDottedLines:       true
//   	Title:                "Table of Contents"
//   	GenerateForwardLinks: true
//   	GenerateBackLinks:    true
//   	Indentation:          "1em"
//   	FontScale:            1
//   Header:
//   	Font:     "Arial"
//   	FontSize: 12
//   Footer:
//   	Font:     "Arial"
//   	FontSize: 12,
func NewObjectOpts() *ObjectOpts {
	return &ObjectOpts{
		UseExternalLinks:  true,
		UseLocalLinks:     true,
		ProduceForms:      true,
		IncludeInOutline:  true,
		CountPages:        true,
		JavascriptDelay:   300,
		Zoom:              1,
		StopSlowScripts:   true,
		ErrorAction:       ActionAbort,
		PrintBackground:   true,
		LoadImages:        true,
		EnableJavascript:  true,
		UseSmartShrinking: true,
		DefaultEncoding:   "utf-8",
		TOC: TOC{
			UseDottedLines:       true,
			Title:                "Table of Contents",
			GenerateForwardLinks: true,
			GenerateBackLinks:    true,
			Indentation:          "1em",
			FontScale:            1,
		},
		Header: Header{
			Font:     "Arial",
			FontSize: 12,
		},
		Footer: Header{
			Font:     "Arial",
			FontSize: 12,
		},
	}
}

// Object represents an HTML document. The contained options are applied only
// to the current object.
type Object struct {
	*ObjectOpts
	settings  *C.wkhtmltopdf_object_settings
	temporary bool
}

// NewObject returns a new object instance from the document at the specified
// location. The location can be a file path or a URL. The object is configured
// using sensible defaults. See NewObjectOpts for the default options.
func NewObject(location string) (*Object, error) {
	return newObject(location, false, nil)
}

// NewObjectWithOpts returns a new object instance from the document at the
// specified location. The location can be a file path or a URL. The object is
// configured using the specified options. If no options are provided, sensible
// defaults are used. See NewObjectOpts for the default options.
func NewObjectWithOpts(opts *ObjectOpts) (*Object, error) {
	return newObject("", false, opts)
}

// NewObjectFromReader creates a new object from the specified reader.
// The object is configured using sensible defaults. See NewObjectOpts for
// the default options.
func NewObjectFromReader(r io.Reader) (*Object, error) {
	file, err := ioutil.TempFile("", "pdf-")
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(file, r); err != nil {
		return nil, err
	}

	tempLocation := file.Name()
	if err := file.Close(); err != nil {
		return nil, err
	}

	location := fmt.Sprintf("%s.html", tempLocation)
	if err := os.Rename(tempLocation, location); err != nil {
		return nil, err
	}

	return newObject(location, true, nil)
}

func newObject(location string, temp bool, opts *ObjectOpts) (*Object, error) {
	if opts == nil {
		opts = NewObjectOpts()
	}
	if location != "" {
		opts.Location = location
	}
	if opts.Location == "" {
		return nil, errors.New("must provide HTML document location")
	}

	settings := C.wkhtmltopdf_create_object_settings()
	if settings == nil {
		return nil, errors.New("could not create object settings")
	}

	return &Object{
		ObjectOpts: opts,
		settings:   settings,
		temporary:  temp,
	}, nil
}

// Destroy releases all resources used by the object.
func (o *Object) Destroy() {
	// Remove temporary file.
	if o.temporary && o.Location != "" {
		os.Remove(o.Location)
		o.Location = ""
		o.temporary = false
	}

	// Destroy settings.
	if o.settings != nil {
		C.wkhtmltopdf_destroy_object_settings(o.settings)
		o.settings = nil
	}
}

func (o *Object) setOption(name, value string) error {
	if name = strings.TrimSpace(name); name == "" {
		return errors.New("object option name cannot be empty")
	}

	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	v := C.CString(value)
	defer C.free(unsafe.Pointer(v))

	if C.wkhtmltopdf_set_object_setting(o.settings, n, v) != 1 {
		return fmt.Errorf("could not set object option: %s", name)
	}

	return nil
}

func (o *Object) setOptions() error {
	if o.settings == nil {
		return errors.New("cannot use uninitialized or destroyed object")
	}

	setter := o.setOption
	opts := []*setOp{
		// General options.
		newSetOp("page", o.Location, optTypeString, setter, true),
		newSetOp("useExternalLinks", o.UseExternalLinks, optTypeBool, setter, true),
		newSetOp("useLocalLinks", o.UseLocalLinks, optTypeBool, setter, true),
		newSetOp("produceForms", o.ProduceForms, optTypeBool, setter, true),
		newSetOp("includeInOutline", o.IncludeInOutline, optTypeBool, setter, true),
		newSetOp("pagesCount", o.CountPages, optTypeBool, setter, true),

		// TOC options.
		newSetOp("toc.useDottedLines", o.TOC.UseDottedLines, optTypeBool, setter, true),
		newSetOp("toc.captionText", o.TOC.Title, optTypeString, setter, true),
		newSetOp("toc.forwardLinks", o.TOC.GenerateForwardLinks, optTypeBool, setter, true),
		newSetOp("toc.backLinks", o.TOC.GenerateBackLinks, optTypeBool, setter, true),
		newSetOp("toc.indentation", o.TOC.Indentation, optTypeString, setter, false),
		newSetOp("toc.fontScale", o.TOC.FontScale, optTypeFloat, setter, false),

		// Header options.
		newSetOp("header.fontName", o.Header.Font, optTypeString, setter, false),
		newSetOp("header.fontSize", o.Header.FontSize, optTypeUint, setter, false),
		newSetOp("header.left", o.Header.ContentLeft, optTypeString, setter, true),
		newSetOp("header.center", o.Header.ContentCenter, optTypeString, setter, true),
		newSetOp("header.right", o.Header.ContentRight, optTypeString, setter, true),
		newSetOp("header.line", o.Header.DisplaySeparator, optTypeBool, setter, true),
		newSetOp("header.spacing", o.Header.Spacing, optTypeFloat, setter, true),
		newSetOp("header.htmlUrl", o.Header.CustomLocation, optTypeString, setter, true),

		// Footer options.
		newSetOp("footer.fontName", o.Footer.Font, optTypeString, setter, false),
		newSetOp("footer.fontSize", o.Footer.FontSize, optTypeUint, setter, false),
		newSetOp("footer.left", o.Footer.ContentLeft, optTypeString, setter, true),
		newSetOp("footer.center", o.Footer.ContentCenter, optTypeString, setter, true),
		newSetOp("footer.right", o.Footer.ContentRight, optTypeString, setter, true),
		newSetOp("footer.line", o.Footer.DisplaySeparator, optTypeBool, setter, true),
		newSetOp("footer.spacing", o.Footer.Spacing, optTypeFloat, setter, true),
		newSetOp("footer.htmlUrl", o.Footer.CustomLocation, optTypeString, setter, true),

		// Load options.
		newSetOp("load.username", o.Username, optTypeString, setter, false),
		newSetOp("load.password", o.Password, optTypeString, setter, false),
		newSetOp("load.jsdelay", o.JavascriptDelay, optTypeUint, setter, false),
		newSetOp("load.windowStatus", o.WindowStatus, optTypeString, setter, false),
		newSetOp("load.zoomFactor", o.Zoom, optTypeFloat, setter, false),
		newSetOp("load.blockLocalFileAccess", o.BlockLocalFileAccess, optTypeBool, setter, true),
		newSetOp("load.stopSlowScripts", o.StopSlowScripts, optTypeBool, setter, true),
		newSetOp("load.loadErrorHandling", string(o.ErrorAction), optTypeString, setter, false),
		newSetOp("load.proxy", o.Proxy, optTypeString, setter, false),

		// Web options.
		newSetOp("web.background", o.PrintBackground, optTypeBool, setter, true),
		newSetOp("web.loadImages", o.LoadImages, optTypeBool, setter, true),
		newSetOp("web.enableJavascript", o.EnableJavascript, optTypeBool, setter, true),
		newSetOp("web.enableIntelligentShrinking", o.UseSmartShrinking, optTypeBool, setter, true),
		newSetOp("web.minimumFontSize", o.MinFontSize, optTypeUint, setter, false),
		newSetOp("web.defaultEncoding", o.DefaultEncoding, optTypeString, setter, false),
		newSetOp("web.printMediaType", o.UsePrintMediaType, optTypeBool, setter, true),
		newSetOp("web.userStyleSheet", o.UserStylesheetLocation, optTypeString, setter, true),
		newSetOp("web.enablePlugins", o.EnablePlugins, optTypeBool, setter, true),
	}

	for _, opt := range opts {
		if err := opt.execute(); err != nil {
			return err
		}
	}

	return nil
}
