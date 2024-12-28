package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	pdf "github.com/adrg/go-wkhtmltopdf"
)

func init() {
	// Set main function to run on the main thread.
	runtime.LockOSThread()
}

var run = make(chan func())

func main() {
	// Initialize library.
	if err := pdf.Init(); err != nil {
		log.Fatal(err)
	}
	defer pdf.Destroy()

	// Start HTTP server on another Go routine.
	go startServer()

	// Listen for functions that need to run on the main thread.
	var quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	for {
		select {
		case f := <-run:
			f()
		case <-quit:
			log.Println("shutting down")
			return
		}
	}
}

// callFunc calls the provided function on the main thread.
func callFunc(f func() error) error {
	err := make(chan error)
	run <- func() {
		err <- f()
	}
	return <-err
}

type requestData struct {
	ConverterOpts *pdf.ConverterOpts `json:"converterOpts"`
	ObjectOpts    *pdf.ObjectOpts    `json:"objectOpts"`
}

func startServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path.
		if r.Method != http.MethodPost || r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// Set default options. Any option fields specified in the request
		// body will overwrite the defaults.
		data := &requestData{
			ConverterOpts: pdf.NewConverterOpts(),
			ObjectOpts:    pdf.NewObjectOpts(),
		}

		// Decode request body.
		if err := json.NewDecoder(r.Body).Decode(data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		// Convert the page at the specified URL to PDF.
		out := bytes.NewBuffer(nil)
		if err := callFunc(func() error {
			// Create object with options.
			object, err := pdf.NewObjectWithOpts(data.ObjectOpts)
			if err != nil {
				return err
			}

			// Create converter with options.
			converter, err := pdf.NewConverterWithOpts(data.ConverterOpts)
			if err != nil {
				log.Fatal(err)
			}
			defer converter.Destroy()

			// Add object to the converter.
			converter.Add(object)

			// Run converter. Due to a limitation of the `wkhtmltox` library,
			// the conversion must be performed on the main thread.
			return converter.Run(out)
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// Serve converted file.
		w.Header().Set("Content-Disposition", "attachment; filename=download.pdf")
		w.Header().Set("Content-Type", "application/pdf")
		if _, err := io.Copy(w, out); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
