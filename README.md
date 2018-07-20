go-wkhtmltopdf
==============
[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/adrg/go-wkhtmltopdf)
[![License: MIT](http://img.shields.io/badge/license-mit-red.svg?style=flat-square)](http://opensource.org/licenses/mit)

Implements [wkhtmltopdf](http://wkhtmltopdf.org) Go bindings. It can be used to convert HTML documents to PDFs.
The package does not use the wkhtmltopdf binary. Instead, it uses the wkhtmltox library directly.

Full documentation can be found at: http://godoc.org/github.com/adrg/go-wkhtmltopdf

## Installation
    go get github.com/adrg/go-wkhtmltopdf

Alternatively, you may want to clone this repository if you're running a OS other than Windows or a more up to date version of wkhtmltopdf, since the wkhtmltox library shipped here is actually a DLL, version 0.12.4.

In this case, we've got a How To waiting you down this page. Keep going.

## Usage

We encourage you to check examples folder out to build and run this very example. You gotta enjoy it yourself.

```go
package main

import (
	"fmt"
	"log"
	"os"

	pdf "github.com/adrg/go-wkhtmltopdf"
)

func main() {
	pdf.Init()
	defer pdf.Destroy()

	// Create object from file
	object, err := pdf.NewObject("sample1.html")
	if err != nil {
		log.Fatal(err)
	}
	object.SetOption("header.center", "This is the header of the first page")
	object.SetOption("footer.right", "[page]")
	object.SetOption("load.windowStatus", "ready")

	// Create object from url
	object2, err := pdf.NewObject("https://google.com")
	if err != nil {
		log.Fatal(err)
	}
	object2.SetOption("footer.right", "[page]")

	// Create object from reader
	file, err := os.Open("sample2.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	object3, err := pdf.NewObjectFromReader(file)
	if err != nil {
		log.Fatal(err)
	}
	object3.SetOption("footer.right", "[page]")

	// Create converter
	converter := pdf.NewConverter()
	defer converter.Destroy()

	// Add created objects to the converter
	converter.AddObject(object)
	converter.AddObject(object2)
	converter.AddObject(object3)

	// Add converter options
	converter.SetOption("documentTitle", "Sample document")
	converter.SetOption("margin.left", "10mm")
	converter.SetOption("margin.right", "10mm")
	converter.SetOption("margin.top", "10mm")
	converter.SetOption("margin.bottom", "10mm")

	// Convert the objects and get the output PDF document
	output, err := converter.Convert()
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("example.pdf", output, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
```

## How to install, build and run the shipped example

Even though this is a very simple process, we've got a Makefile to help us do it over and over without much typing.

### 1: Clone this repo
	git clone https://github.com/adrg/go-wkhtmltopdf.git

Great! Thanks for cloning. Now go ahead and change to the project directory.

	cd go-wkhtmltopdf

### 2: Bring in your own wkhtmltox library

If you want to bring in your wkhtmltox library, just copy the files to the **./wkhtmltox** directory and your are good to go.

### 3: Setup Makefile
	PKG_SRC_PATH := $(GOPATH)/src/github.com/adrg/go-wkhtmltopdf

What you want to do is to set this variable to your actual $GOPATH and package source path as well. This depends on the repository you cloned.

### 4: Install this package
	make install

This is going to copy this directory to your $GOPATH source structure and **go install** it, making it available to your own programs to use it.

### 5: Build the example
	make build-example

As result of this command you going to get a **run-example.exe** file in the example directory.

### 6: Finally, run the example
	make run-example

Instead of this command, you can simply change to **./example** directory and fire **./run-example.exe**, which is going to result in the **example.pdf** file.

Voilà!

## Conversion options

### Converter options
```
- size.paperSize    Paper size of the output document (e.g. A4)
- size.width        Width of the output document (e.g.  4cm)
- size.height       Height of the output document (e.g. 12in)
- orientation       Orientation of the output document (values: Landscape, Portrait)
- colorMode         Color mode to use (values: Color, Grayscale)
- resolution        Most likely has no effect
- dpi               DPI to use for printing (e.g. 80)
- pageOffset        Offset added to page numbers when printing headers, footers and tables of contents
- copies            Copies of object to print (e.g. 2)
- collate           Specifies if the copies should be collated (values: true, false)
- outline           Specifies if an outline should be generated (values: true, false)
- outlineDepth      The maximum depth level of the outline (e.g. 4)
- dumpOutline       Dump an XML representation of the outline to the specified file
- documentTitle     Title of the output document
- useCompression    Use lossless compression for the output document (values: true, false)
- margin.top        Size of the top margin (e.g. 2cm)
- margin.bottom     Size of the bottom margin (e.g. 3in)
- margin.left       Size of the left margin (e.g. 4mm)
- margin.right      Size of the right margin (e.g. 2cm)
- imageDPI          Maximum DPI value to use for images in the output document
- imageQuality      JPEG compression factor to use when producing the output document (e.g. 92)
- load.cookieJar    Path of file used to load and store cookies
```

### Object options

##### Load options
```
- load.username                Username to use for logging into a website (e.g. bart)
- load.password                Password to use for logging into a website (e.g. elbarto)
- load.jsdelay                 Milliseconds to wait after page load before print start (e.g. 1200)
- load.windowStatus            Wait until window.status is equal to this string before rendering page
- load.zoomFactor              Zoom of the content (e.g. 2.2)
- load.blockLocalFileAccess    Block local files from accessing other local files (values: true, false)
- load.stopSlowScript          Stop slow running javascript (values: true, false)
- load.debugJavascript         Pass JS warnings and errors to the warning callback (values: true, false)
- load.loadErrorHandling       Action to take on object conversion failure (values: abort, skip, ignore)
- load.proxy                   Proxy to use when loading the object
```

##### Header options
```
- header.fontSize    Font size to use for the header (e.g. 13)
- header.fontName    Name of the font to use for the header (e.g. verdana)
- header.left        Text to print in the left part of the header
- header.center      Text to print in the center part of the header
- header.right       Text to print in the right part of the header
- header.line        Specifies whether a line is printed under the header (values: true, false)
- header.spacing     Amount of space to put between the header and the content (e.g. 1.8)
- header.htmlUrl     URL for a HTML document to use for the header
```

##### Footer options
```
- footer.fontSize    Font size to use for the footer (e.g. 13)
- footer.fontName    Name of the font to use for the footer (e.g. verdana)
- footer.left        Text to print in the left part of the footer
- footer.center      Text to print in the center part of the footer
- footer.right       Text to print in the right part of the footer
- footer.line        Specifies whether a line is printed above the footer (values: true, false)
- footer.spacing     Amount of space to put between the footer and the content (e.g. 1.8)
- footer.htmlUrl     URL for a HTML document to use for the footer
```

##### Web page options
```
- web.background                    Specifies if the background is printed (values: true, false)
- web.loadImages                    Specifies if images are loaded (values: true, false)
- web.enableJavascript              Specifies if Javascript is enabled (values: true, false)
- web.enableIntelligentShrinking    Enable smart shrinking to fit more content (values: true, false)
- web.minimumFontSize               Minimum font size allowed (e.g. 9)
- web.printMediaType                Use the print media type (values: true, false)
- web.defaultEncoding               Specifies the default document encoding (e.g. utf-8)
- web.userStyleSheet                URL or path to a user specified style sheet
- web.enablePlugins                 Enable NS plugins (values: true, false)
```

For more information see the [wkhtmltopdf documentation](http://wkhtmltopdf.org/usage/wkhtmltopdf.txt)

## License
Copyright (c) 2016 Adrian-George Bostan.

This project is licensed under the [MIT license](http://opensource.org/licenses/MIT). See LICENSE for more details.
