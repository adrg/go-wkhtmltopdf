package main

import (
	"io/ioutil"
	"log"
	"os"

	pdf "github.com/leandrosilva/go-wkhtmltopdf" // <- This is what I use on my machine
	// pdf "github.com/adrg/go-wkhtmltopdf"         <- You may want to use this instead
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
