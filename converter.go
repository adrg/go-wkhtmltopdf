package pdf

/*
#cgo LDFLAGS: -lwkhtmltox
#include <stdio.h>
#include <stdlib.h>
#include <wkhtmltox/pdf.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

type Converter struct {
	converter *C.wkhtmltopdf_converter
	settings  *C.wkhtmltopdf_global_settings

	objects []*Object
}

func NewConverter() *Converter {
	settings := C.wkhtmltopdf_create_global_settings()
	converter := C.wkhtmltopdf_create_converter(settings)

	return &Converter{converter: converter, settings: settings}
}

func (c *Converter) AddObject(object *Object) {
	C.wkhtmltopdf_add_object(c.converter, object.settings, nil)
	c.objects = append(c.objects, object)
}

func (c *Converter) SetOption(name, value string) error {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	v := C.CString(value)
	defer C.free(unsafe.Pointer(v))

	if C.wkhtmltopdf_set_global_setting(c.settings, n, v) != 1 {
		return errors.New("Could not set option")
	}

	return nil
}

func (c *Converter) Convert() ([]byte, error) {
	if len(c.objects) == 0 {
		return nil, errors.New("You must add at least one object to convert")
	}
	if err := c.SetOption("out", ""); err != nil {
		return nil, errors.New("Could not set output option")
	}

	if C.wkhtmltopdf_convert(c.converter) != 1 {
		return nil, errors.New("Conversion failed")
	}

	var output *C.uchar
	size := C.wkhtmltopdf_get_output(c.converter, &output)
	if size == 0 {
		return nil, errors.New("Could not retrieve converted object")
	}

	return C.GoBytes(unsafe.Pointer(output), C.int(size)), nil
}

func (c *Converter) Destroy() {
	for _, o := range c.objects {
		o.destroy()
	}

	C.wkhtmltopdf_destroy_converter(c.converter)
}
