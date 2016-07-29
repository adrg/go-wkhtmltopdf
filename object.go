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
	"unsafe"
)

type Object struct {
	settings *C.wkhtmltopdf_object_settings

	filename  string
	temporary bool
}

func newObject(filename string, temporary bool) (*Object, error) {
	settings := C.wkhtmltopdf_create_object_settings()
	o := &Object{settings: settings, filename: filename, temporary: temporary}

	if err := o.SetOption("page", filename); err != nil {
		return nil, err
	}

	return o, nil
}

func NewObject(filename string) (*Object, error) {
	return newObject(filename, false)
}

func NewObjectFromReader(r io.Reader) (*Object, error) {
	file, err := ioutil.TempFile("", "pdf-")
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(file, r); err != nil {
		return nil, err
	}

	tempFilename := file.Name()
	if err := file.Close(); err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("%s.html", tempFilename)
	if err := os.Rename(tempFilename, filename); err != nil {
		return nil, err
	}

	return newObject(filename, true)
}

func (o *Object) SetOption(name, value string) error {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	v := C.CString(value)
	defer C.free(unsafe.Pointer(v))

	if C.wkhtmltopdf_set_object_setting(o.settings, n, v) != 1 {
		return errors.New("Could not set option")
	}

	return nil
}

func (o *Object) destroy() {
	if !o.temporary {
		return
	}

	os.Remove(o.filename)
}
