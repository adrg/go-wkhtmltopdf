<h1 align="center">
  <div>
    <img src="https://raw.githubusercontent.com/adrg/adrg.github.io/master/assets/projects/go-wkhtmltopdf/logo.svg" alt="go-wkhtmltopdf logo"/>
  </div>
</h1>

<h3 align="center">Go bindings and high-level HTML to PDF conversion interface.</h3>

<p align="center">
    <a href="https://github.com/adrg/go-wkhtmltopdf/actions/workflows/tests.yml">
        <img alt="Tests status" src="https://github.com/adrg/go-wkhtmltopdf/actions/workflows/tests.yml/badge.svg">
    </a>
    <a href="https://pkg.go.dev/github.com/adrg/go-wkhtmltopdf">
        <img alt="pkg.go.dev documentation" src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white">
    </a>
    <a href="https://opensource.org/licenses/MIT" rel="nofollow">
        <img alt="MIT license" src="https://img.shields.io/github/license/adrg/go-wkhtmltopdf"/>
    </a>
    <a href="https://goreportcard.com/report/github.com/adrg/go-wkhtmltopdf">
        <img alt="Go report card" src="https://goreportcard.com/badge/github.com/adrg/go-wkhtmltopdf" />
    </a>
    <a href="https://discord.gg/Jd63kBf">
        <img alt="Discord channel" src="https://img.shields.io/discord/767381740427542588?label=discord" />
    </a>
    <a href="https://github.com/adrg/go-wkhtmltopdf/issues">
        <img alt="GitHub issues" src="https://img.shields.io/github/issues/adrg/go-wkhtmltopdf">
    </a>
    <a href="https://ko-fi.com/T6T72WATK">
        <img alt="Buy me a coffee" src="https://img.shields.io/static/v1.svg?label=%20&message=Buy%20me%20a%20coffee&color=579fbf&logo=buy%20me%20a%20coffee&logoColor=white"/>
    </a>
</p>

Implements [wkhtmltopdf](https://wkhtmltopdf.org) Go bindings. It can be used to convert HTML documents to PDF files.
The package does not use the `wkhtmltopdf` binary. Instead, it uses the `wkhtmltox` library directly.

Full documentation can be found at https://pkg.go.dev/github.com/adrg/go-wkhtmltopdf.

**Examples**

* [Basic usage](examples/basic-usage/main.go)
* [Converter callbacks](examples/converter-callbacks/main.go)
* [Convert HTML document based on JSON input](examples/json-input/main.go)
* [Basic web page to PDF conversion server](examples/http-server)
* [Configurable web page to PDF conversion server](examples/http-server-advanced)

> Note: The `HTML` to `PDF` conversion (calls to the `Converter.Run` method) must be performed on the main thread.
> This is a limitation of the `wkhtmltox` library. Please see the `HTTP` server [example](examples/http-server)
> for more information.

## Prerequisites

In order to use the package, `wkhtmltox` must be installed. Installation packages
for multiple operating systems can be found at [https://builds.wkhtmltopdf.org](https://wkhtmltopdf.org/downloads.html).

Please see the wiki pages of this project for detailed installation instructions.
- [Install on Linux](https://github.com/adrg/go-wkhtmltopdf/wiki/Install-on-Linux)
- [Install on Windows](https://github.com/adrg/go-wkhtmltopdf/wiki/Install-on-Windows)

> Note: `wkhtmltox` does not seem to be actively maintained. Please see the [project status](https://wkhtmltopdf.org/status.html) for more information, recommendations and future plans.

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
	// Initialize library.
	if err := pdf.Init(); err != nil {
		log.Fatal(err)
	}
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
	object2.Footer.ContentLeft = "[date]"
	object2.Footer.ContentCenter = "Sample footer information"
	object2.Footer.ContentRight = "[page]"
	object2.Footer.DisplaySeparator = true

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

	// Create output file.
	outFile, err := os.Create("out.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	// Run converter. Due to a limitation of the `wkhtmltox` library, the
	// conversion must be performed on the main thread.
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
See [CONTRIBUTING.MD](CONTRIBUTING.md).

**Contributors**:
[adrg](https://github.com/adrg),
[leandrosilva](https://github.com/leandrosilva),
[MicahParks](https://github.com/MicahParks).

## References

For more information see the [wkhtmltopdf documentation](https://wkhtmltopdf.org/usage/wkhtmltopdf.txt)
and the [wkhtmltox documentation](https://wkhtmltopdf.org/libwkhtmltox).

## License

Copyright (c) 2016 Adrian-George Bostan.

This project is licensed under the [MIT license](https://opensource.org/licenses/MIT).
See [LICENSE](LICENSE) for more details.
