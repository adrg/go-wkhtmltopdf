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

Alternatively, you may want to clone this repository if you're running a OS other than Windows or a more up to date version of wkhtmltopdf, since the wkhtmltox library shipped here is actually a DLL, version 0.12.4.

In this case, we've got a **How To** waiting for you down this page. Keep going.

## Usage

We encourage you to check **examples** folder out to build and run this very example (i.e. ex1). You gotta enjoy it yourself, don't you?

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

## How to install, build and run those shipped examples

Even though this is a very simple process, we've got a Makefile to help us do it over and over without much typing.

### 1: Clone this repo

	git clone https://github.com/adrg/go-wkhtmltopdf.git

Great! Thanks for cloning. Now go ahead and change to the project directory.

	cd go-wkhtmltopdf

### 2: Bring in your own wkhtmltox library

If you want to bring in your wkhtmltox library, just copy the files to the **./wkhtmltox** directory and you are good to go. Otherwise, skip this step.

Also you might want (or have) to provide environment variable `CGO_LDFLAGS` depending on your operational system, etc. Just keep that in mind. And if you need some help, please, refer to [cgo doc](https://golang.org/cmd/cgo/).

### 3: Setup Makefile

	export PKG_SRC_PATH=$GOPATH/src/github.com/adrg/go-wkhtmltopdf

What you want to do is to set this variable to your actual $GOPATH and package source path as well. This depends on the repository you cloned.

### 4: Install this package

	make install

This is going to copy this directory to your $GOPATH source structure and **go install** it, making it available to your own programs to use it.

### 5: Build the example 1

	cd ./examples/ex1
	make build

As result of this command you going to get a **ex1.exe** file in ths very directory.

### 6: Finally, run it

	make run

Instead of this command, you can simply fire **./ex1.exe**, which is going to result in the **example.pdf** file.

Voilà!

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
