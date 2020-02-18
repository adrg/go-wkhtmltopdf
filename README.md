go-wkhtmltopdf
==============
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/adrg/go-wkhtmltopdf)
[![License: MIT](https://img.shields.io/badge/license-mit-red.svg?style=flat-square)](https://opensource.org/licenses/mit)
[![Go Report Card](https://goreportcard.com/badge/github.com/adrg/go-wkhtmltopdf)](https://goreportcard.com/report/github.com/adrg/go-wkhtmltopdf)

Implements [wkhtmltopdf](https://wkhtmltopdf.org) Go bindings. It can be used to convert HTML documents to PDF files.
The package does not use the wkhtmltopdf binary. Instead, it uses the wkhtmltox library directly.

Full documentation can be found at: https://godoc.org/github.com/adrg/go-wkhtmltopdf

## Requirements

In order to use the package, wkhtmltox must be installed. Installation packages
for multiple operating systems can be found at
[https://builds.wkhtmltopdf.org](https://builds.wkhtmltopdf.org).

On Debian based distributions, use dpkg to install the downloaded installation package.
```
sudo dpkg -i wkhtmltox.deb
sudo ldconfig
```

## Installation
    go get github.com/adrg/go-wkhtmltopdf

## Usage

```go
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
```

## Stargazers over time

[![Stargazers over time](https://starchart.cc/adrg/go-wkhtmltopdf.svg)](https://starchart.cc/adrg/go-wkhtmltopdf)

## Contributing

Contributions in the form of pull requests, issues or just general feedback,
are always welcome.
See [CONTRIBUTING.MD](https://github.com/adrg/go-wkhtmltopdf/blob/master/CONTRIBUTING.md).

## References

For more information see the [wkhtmltopdf documentation](https://wkhtmltopdf.org/usage/wkhtmltopdf.txt)
and the [wkhtmltox documentation](https://wkhtmltopdf.org/libwkhtmltox).

## License

Copyright (c) 2016 Adrian-George Bostan.

This project is licensed under the [MIT license](https://opensource.org/licenses/MIT).
See [LICENSE](https://github.com/adrg/go-wkhtmltopdf/blob/master/LICENSE) for more details.
