package pdf

/*
#cgo LDFLAGS: -lwkhtmltox
#include <stdio.h>
#include <stdlib.h>
#include <wkhtmltox/pdf.h>

typedef const char cchar;
typedef const int cint;

extern void converterWarningCb(wkhtmltopdf_converter* converter, cchar* msg);
extern void converterErrorCb(wkhtmltopdf_converter* converter, cchar* msg);
extern void converterPhaseChangedCb(wkhtmltopdf_converter* converter);
extern void converterProgressChangedCb(wkhtmltopdf_converter* converter, cint progress);
extern void converterFinishedCb(wkhtmltopdf_converter* converter, cint status);

static inline void converter_initialize_callbacks(wkhtmltopdf_converter* c) {
	wkhtmltopdf_set_warning_callback(c, converterWarningCb);
	wkhtmltopdf_set_error_callback(c, converterErrorCb);
	wkhtmltopdf_set_phase_changed_callback(c, converterPhaseChangedCb);
	wkhtmltopdf_set_progress_changed_callback(c, converterProgressChangedCb);
	wkhtmltopdf_set_finished_callback(c, converterFinishedCb);
}
*/
import "C"
import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unsafe"
)

// Colorspace represents the color mode of the output document content.
type Colorspace string

// Colorspace values.
const (
	Color     Colorspace = "Color"
	Grayscale Colorspace = "Grayscale"
)

// Orientation represents the orientation of the output document pages.
type Orientation string

// Page orientation values.
const (
	Portrait  Orientation = "Portrait"
	Landscape Orientation = "Landscape"
)

// PaperSize represents the size of the output document pages.
type PaperSize string

// Paper size values.
const (
	A0        PaperSize = "A0"        // 841 x 1189 mm
	A1        PaperSize = "A1"        // 594 x 841 mm
	A2        PaperSize = "A2"        // 420 x 594 mm
	A3        PaperSize = "A3"        // 297 x 420 mm
	A4        PaperSize = "A4"        // 210 x 297 mm
	A5        PaperSize = "A5"        // 148 x 210 mm
	A6        PaperSize = "A6"        // 105 x 148 mm
	A7        PaperSize = "A7"        // 74 x 105 mm
	A8        PaperSize = "A8"        // 52 x 74 mm
	A9        PaperSize = "A9"        // 37 x 52 mm
	B0        PaperSize = "B0"        // 1000 x 1414 mm
	B1        PaperSize = "B1"        // 707 x 1000 mm
	B2        PaperSize = "B2"        // 500 x 707 mm
	B3        PaperSize = "B3"        // 353 x 500 mm
	B4        PaperSize = "B4"        // 250 x 353 mm
	B5        PaperSize = "B5"        // 176 x 250 mm
	B6        PaperSize = "B6"        // 125 x 176 mm
	B7        PaperSize = "B7"        // 88 x 125 mm
	B8        PaperSize = "B8"        // 62 x 88 mm
	B9        PaperSize = "B9"        // 33 x 62 mm
	B10       PaperSize = "B10"       // 31 x 44 mm
	C5E       PaperSize = "C5E"       // 163 x 229 mm
	Comm10E   PaperSize = "Comm10E"   // 105 x 241 mm
	DLE       PaperSize = "DLE"       // 110 x 220 mm
	Executive PaperSize = "Executive" // 190.5 x 254 mm
	Folio     PaperSize = "Folio"     // 210 x 330 mm
	Ledger    PaperSize = "Ledger"    // 431.8 x 279.4 mm
	Legal     PaperSize = "Legal"     // 215.9 x 355.6 mm
	Letter    PaperSize = "Letter"    // 215.9 x 279.4 mm
	Tabloid   PaperSize = "Tabloid"   // 279.4 x 431.8 mm
)

// ConverterOpts defines a set of options to be used in the conversion process.
type ConverterOpts struct {
	// The paper size of the output document.
	// E.g.: A4.
	PaperSize PaperSize `json:"paperSize" yaml:"paperSize"`

	// The width of the output document.
	// E.g.: "4cm".
	Width string `json:"width" yaml:"width"`

	// The height of the output document.
	// E.g. "12in".
	Height string `json:"height" yaml:"height"`

	// The orientation of the output document.
	// E.g.: Portrait.
	Orientation Orientation `json:"orientation" yaml:"orientation"`

	// The color mode of the output document.
	// E.g.: Color.
	Colorspace Colorspace `json:"colorspace" yaml:"colorspace"`

	// DPI of the output document.
	// E.g.: 96.
	DPI uint64 `json:"dpi" yaml:"dpi"`

	// A number added to all page numbers when rendering headers, footers and
	// tables of contents.
	PageOffset int64 `json:"pageOffset" yaml:"pageOffset"`

	// Copies of the converted documents to be included in the output document.
	// E.g.: 1.
	Copies uint64 `json:"copies" yaml:"copies"`

	// Specifies whether copies should be collated.
	Collate bool `json:"collate" yaml:"collate"`

	// The title of the output document.
	Title string `json:"title" yaml:"title"`

	// Specifies whether outlines should be generated for the output document.
	GenerateOutline bool `json:"generateOutline" yaml:"generateOutline"`

	// The maximum number of nesting levels in outlines.
	// E.g.: 4.
	OutlineDepth uint64 `json:"outlineDepth" yaml:"outlineDepth"`

	// A location to write an XML representation of the generated outlines.
	OutlineDumpPath string `json:"outlineDumpPath" yaml:"outlineDumpPath"`

	// Specifies whether the conversion process should use lossless compression.
	UseCompression bool `json:"useCompression" yaml:"useCompression"`

	// Size of the top margin. (e.g. "2cm")
	// E.g.: "1cm".
	MarginTop string `json:"marginTop" yaml:"marginTop"`

	// Size of the bottom margin. (e.g. "2cm")
	// E.g.: "1cm".
	MarginBottom string `json:"marginBottom" yaml:"marginBottom"`

	// Size of the left margin. (e.g. "2cm")
	// E.g.: "10mm".
	MarginLeft string `json:"marginLeft" yaml:"marginLeft"`

	// Size of the right margin. (e.g. "2cm")
	// E.g.: "10mm".
	MarginRight string `json:"marginRight" yaml:"marginRight"`

	// The maximum number of DPI for the images in the output document.
	// E.g.: 600.
	ImageDPI uint64 `json:"imageDPI" yaml:"imageDPI"`

	// The compression factor to use for the JPEG images in the output document.
	// E.g.: 100 (range 0-100).
	ImageQuality uint64 `json:"imageQuality" yaml:"imageQuality"`

	// Path of the file used to load and store cookies for web objects.
	CookieJarPath string `json:"cookieJarPath" yaml:"cookieJarPath"`
}

// NewConverterOpts returns a new instance of converter options, configured
// using sensible defaults.
//
//   Defaults options:
//
//   PaperSize:       A4
//   Orientation:     Portrait
//   Colorspace:      Color
//   DPI:             96
//   Copies:          1
//   Collate:         true
//   GenerateOutline: true
//   UseCompression:  true
//   MarginLeft:      "10mm"
//   MarginRight:     "10mm"
//   ImageDPI:        600
//   ImageQuality:    100
func NewConverterOpts() *ConverterOpts {
	return &ConverterOpts{
		PaperSize:       A4,
		Orientation:     Portrait,
		Colorspace:      Color,
		DPI:             96,
		Copies:          1,
		Collate:         true,
		GenerateOutline: true,
		UseCompression:  true,
		MarginLeft:      "10mm",
		MarginRight:     "10mm",
		ImageDPI:        600,
		ImageQuality:    100,
	}
}

// Converter represents an HTML to PDF converter. The contained options are
// applied to all converted objects.
type Converter struct {
	*ConverterOpts
	converter *C.wkhtmltopdf_converter
	settings  *C.wkhtmltopdf_global_settings
	objects   []*Object
	phases    []string

	// Warning is called when a warning is issued in the conversion process.
	Warning func(msg string)

	// Error is called when an error is encountered in the conversion process.
	Error func(msg string)

	// PhaseChanged is called when the conversion phase changes.
	PhaseChanged func(phaseIndex int)

	// ProgressChanged is called when the conversion progress changes.
	// The progress is reported for each conversion phase.
	ProgressChanged func(progressPercent int)

	// Finished is called when the conversion process ends.
	Finished func(success bool)
}

// NewConverter returns a new converter instance, configured using sensible
// defaults. See NewConverterOpts for the default options.
func NewConverter() (*Converter, error) {
	return NewConverterWithOpts(nil)
}

// NewConverterWithOpts returns a new converter instance, configured using the
// specified options. If no options are provided, sensible defaults are used.
// See NewConverterOpts for the default options.
func NewConverterWithOpts(opts *ConverterOpts) (*Converter, error) {
	if opts == nil {
		opts = NewConverterOpts()
	}

	// Create converter settings.
	settings := C.wkhtmltopdf_create_global_settings()
	if settings == nil {
		return nil, errors.New("could not create converter settings")
	}

	// Create converter.
	cConverter := C.wkhtmltopdf_create_converter(settings)
	if cConverter == nil {
		return nil, errors.New("could not create converter")
	}

	converter := &Converter{
		ConverterOpts: opts,
		converter:     cConverter,
		settings:      settings,
	}

	// Initialize converter callbacks.
	C.converter_initialize_callbacks(cConverter)

	// Retrieve conversion phases.
	phaseCount := int(C.wkhtmltopdf_phase_count(cConverter))

	converter.phases = make([]string, phaseCount)
	for i := 0; i < phaseCount; i++ {
		converter.phases[i] = C.GoString(C.wkhtmltopdf_phase_description(cConverter, C.int(i)))
	}

	// Add converter to object registry.
	registry.add(objectID(cConverter), converter)

	return converter, nil
}

// Add appends the specified object to the list of objects to be converted.
func (c *Converter) Add(object *Object) {
	c.objects = append(c.objects, object)
}

// Run performs the conversion and copies the output to the provided writer.
func (c *Converter) Run(w io.Writer) error {
	if c.converter == nil {
		return errors.New("cannot use uninitialized or destroyed converter")
	}
	if w == nil {
		return errors.New("the provided writer cannot be nil")
	}

	// Set converter and object options.
	if len(c.objects) == 0 {
		return errors.New("must add at least one object to convert")
	}
	if err := c.setOptions(); err != nil {
		return err
	}

	// Convert objects.
	if C.wkhtmltopdf_convert(c.converter) != 1 {
		return errors.New("could not convert the added objects")
	}

	// Get conversion output buffer.
	var output *C.uchar
	size := C.wkhtmltopdf_get_output(c.converter, &output)
	if size == 0 {
		return errors.New("could not retrieve the converted file")
	}

	// Copy output to the provided writer.
	buf := bytes.NewBuffer(C.GoBytes(unsafe.Pointer(output), C.int(size)))
	if _, err := io.Copy(w, buf); err != nil {
		return err
	}

	return nil
}

// Destroy releases all resources used by the converter.
func (c *Converter) Destroy() {
	// Destroy settings.
	if c.settings != nil {
		C.wkhtmltopdf_destroy_global_settings(c.settings)
		c.settings = nil
	}

	// Destroy converter.
	if c.converter != nil {
		registry.remove(objectID(c.converter))
		C.wkhtmltopdf_destroy_converter(c.converter)
		c.converter = nil
	}

	// Destroy converter objects.
	for _, o := range c.objects {
		o.Destroy()
	}
	c.objects = nil
}

// Phases returns the list of phases undergone in the conversion process.
func (c *Converter) Phases() []string {
	return c.phases
}

// PhaseDescription returns the description of the phase with the specified
// index. If the phase index is invalid, the method returns an empty string.
func (c *Converter) PhaseDescription(phaseIndex int) string {
	if phaseIndex < 0 || phaseIndex >= len(c.phases) {
		return ""
	}

	return c.phases[phaseIndex]
}

// CurrentPhaseIndex returns the index of the current conversion phase.
func (c *Converter) CurrentPhaseIndex() int {
	return int(C.wkhtmltopdf_current_phase(c.converter))
}

func (c *Converter) setOption(name, value string) error {
	if name = strings.TrimSpace(name); name == "" {
		return errors.New("converter option name cannot be empty")
	}

	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	v := C.CString(value)
	defer C.free(unsafe.Pointer(v))

	if C.wkhtmltopdf_set_global_setting(c.settings, n, v) != 1 {
		return fmt.Errorf("could not set converter option `%s` to `%s`", name, value)
	}

	return nil
}

func (c *Converter) setOptions() error {
	setter := c.setOption
	opts := []*setOp{
		newSetOp("size.pageSize", string(c.PaperSize), optTypeString, setter, false),
		newSetOp("size.width", c.Width, optTypeString, setter, false),
		newSetOp("size.height", c.Height, optTypeString, setter, false),
		newSetOp("orientation", string(c.Orientation), optTypeString, setter, false),
		newSetOp("colorMode", string(c.Colorspace), optTypeString, setter, false),
		newSetOp("dpi", c.DPI, optTypeUint, setter, false),
		newSetOp("pageOffset", c.PageOffset, optTypeInt, setter, true),
		newSetOp("copies", c.Copies, optTypeUint, setter, false),
		newSetOp("collate", c.Collate, optTypeBool, setter, true),
		newSetOp("outline", c.GenerateOutline, optTypeBool, setter, true),
		newSetOp("outlineDepth", c.OutlineDepth, optTypeUint, setter, false),
		newSetOp("dumpOutline", c.OutlineDumpPath, optTypeString, setter, true),
		newSetOp("documentTitle", c.Title, optTypeString, setter, true),
		newSetOp("useCompression", c.UseCompression, optTypeBool, setter, true),
		newSetOp("margin.top", c.MarginTop, optTypeString, setter, false),
		newSetOp("margin.bottom", c.MarginBottom, optTypeString, setter, false),
		newSetOp("margin.left", c.MarginLeft, optTypeString, setter, false),
		newSetOp("margin.right", c.MarginRight, optTypeString, setter, false),
		newSetOp("imageDPI", c.ImageDPI, optTypeUint, setter, false),
		newSetOp("imageQuality", c.ImageQuality, optTypeUint, setter, false),
		newSetOp("load.cookieJar", c.CookieJarPath, optTypeString, setter, true),
		newSetOp("out", "", optTypeString, setter, true),
	}

	for _, opt := range opts {
		if err := opt.execute(); err != nil {
			return err
		}
	}

	// Set object options.
	for _, o := range c.objects {
		if err := o.setOptions(); err != nil {
			return err
		}

		C.wkhtmltopdf_add_object(c.converter, o.settings, nil)
	}

	return nil
}

//export converterWarningCb
func converterWarningCb(cConverter *C.wkhtmltopdf_converter, msg *C.cchar) {
	converter := getConverterByID(objectID(cConverter))
	if converter != nil && converter.Warning != nil {
		converter.Warning(C.GoString(msg))
	}
}

//export converterErrorCb
func converterErrorCb(cConverter *C.wkhtmltopdf_converter, msg *C.cchar) {
	converter := getConverterByID(objectID(cConverter))
	if converter != nil && converter.Error != nil {
		converter.Error(C.GoString(msg))
	}
}

//export converterPhaseChangedCb
func converterPhaseChangedCb(cConverter *C.wkhtmltopdf_converter) {
	converter := getConverterByID(objectID(cConverter))
	if converter != nil && converter.PhaseChanged != nil {
		converter.PhaseChanged(converter.CurrentPhaseIndex())
	}
}

//export converterProgressChangedCb
func converterProgressChangedCb(cConverter *C.wkhtmltopdf_converter, progress C.int) {
	converter := getConverterByID(objectID(cConverter))
	if converter != nil && converter.ProgressChanged != nil {
		converter.ProgressChanged(int(progress))
	}
}

//export converterFinishedCb
func converterFinishedCb(cConverter *C.wkhtmltopdf_converter, status C.int) {
	converter := getConverterByID(objectID(cConverter))
	if converter != nil && converter.Finished != nil {
		converter.Finished(status == 1)
	}
}

func getConverterByID(id objectID) *Converter {
	object, ok := registry.get(id)
	if !ok {
		return nil
	}

	converter, _ := object.(*Converter)
	return converter
}
