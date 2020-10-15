package main

import (
	"bytes"
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
	object.Footer.ContentCenter = "This is the header of the first page"
	object.Footer.ContentRight = "[page]"
	object.WindowStatus = "ready"

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
	object3.Footer.ContentLeft = "[date]"
	object3.Footer.ContentCenter = "Sample footer information 3"
	object3.Footer.ContentRight = "[page]"

	// Create converter
	converter, err := pdf.NewConverter()
	if err != nil {
		log.Fatal(err)
	}
	defer converter.Destroy()

	// Add created objects to the converter
	converter.Add(object)
	converter.Add(object2)
	converter.Add(object3)

	// Add converter options
	converter.Title = "Sample document"
	converter.MarginLeft = "10mm"
	converter.MarginRight = "10mm"
	converter.MarginTop = "10mm"
	converter.MarginBottom = "10mm"

	// Convert the objects and get the output PDF document
	output := new(bytes.Buffer)
	err = converter.Run(output)
	if err != nil {
		log.Fatal(err)
	}
	raw := output.Bytes()

	err = ioutil.WriteFile("ex1.pdf", raw, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
