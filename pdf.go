/*
Package pdf implements wkhtmltopdf Go bindings. It can be used to convert HTML documents to PDF files.
The package does not use the wkhtmltopdf binary. Instead, it uses the wkhtmltox library directly.

Example
	package main

	import (
		"log"
		"os"

		pdf "github.com/adrg/go-wkhtmltopdf"
	)

	func main() {
		pdf.Init()
		defer pdf.Destroy()

		// Create object from file.
		object, err := pdf.NewObject("sample1.html")
		if err != nil {
			log.Fatal(err)
		}
		object.Header.ContentCenter = "[title]"
		object.Header.DisplaySeparator = true

		// Create object from URL.
		object2, err := pdf.NewObject("https://google.com")
		if err != nil {
			log.Fatal(err)
		}
		object.Footer.ContentLeft = "[date]"
		object.Footer.ContentCenter = "Sample footer information"
		object.Footer.ContentRight = "[page]"
		object.Footer.DisplaySeparator = true

		// Create object from reader.
		inFile, err := os.Open("sample2.html")
		if err != nil {
			log.Fatal(err)
		}
		defer inFile.Close()

		object3, err := pdf.NewObjectFromReader(inFile)
		if err != nil {
			log.Fatal(err)
		}
		object3.Zoom = 1.5
		object3.TOC.Title = "Table of Contents"

		// Create converter.
		converter, err := pdf.NewConverter()
		if err != nil {
			log.Fatal(err)
		}
		defer converter.Destroy()

		// Add created objects to the converter.
		converter.Add(object)
		converter.Add(object2)
		converter.Add(object3)

		// Set converter options.
		converter.Title = "Sample document"
		converter.PaperSize = pdf.A4
		converter.Orientation = pdf.Landscape
		converter.MarginTop = "1cm"
		converter.MarginBottom = "1cm"
		converter.MarginLeft = "10mm"
		converter.MarginRight = "10mm"

		// Convert objects and save the output PDF document.
		outFile, err := os.Create("out.pdf")
		if err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()

		if err := converter.Run(outFile); err != nil {
			log.Fatal(err)
		}
	}

For more information see http://wkhtmltopdf.org/usage/wkhtmltopdf.txt
*/
package pdf

/*
#cgo LDFLAGS: -lwkhtmltox
#include <stdlib.h>
#include <wkhtmltox/pdf.h>
*/
import "C"
import "errors"

var registry *objectRegistry

// Init initializes the library, allocating all necessary resources.
func Init() error {
	if C.wkhtmltopdf_init(0) != 1 {
		return errors.New("could not initialize library")
	}

	registry = newObjectRegistry()
	return nil
}

// Version returns the version of the library.
func Version() string {
	return C.GoString(C.wkhtmltopdf_version())
}

// HasPatchedQT returns true if the library is built against the
// wkhtmltopdf version of QT.
func HasPatchedQT() bool {
	return C.wkhtmltopdf_extended_qt() != 0
}

// Destroy releases all the resources used by the library.
func Destroy() {
	registry = nil
	C.wkhtmltopdf_deinit()
}
