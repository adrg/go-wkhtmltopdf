package main

import (
	"bytes"
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

func startServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path.
		if r.Method != http.MethodGet || r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// Get URL to convert.
		urls, ok := r.URL.Query()["url"]
		if !ok || len(urls) != 1 {
			http.Error(w, "invalid request query", http.StatusBadRequest)
		}
		url := urls[0]

		// Convert the page at the specified URL to PDF.
		out := bytes.NewBuffer(nil)
		if err := callFunc(func() error {
			// Create object from URL.
			object, err := pdf.NewObject(string(url))
			if err != nil {
				return err
			}

			// Create converter.
			converter, err := pdf.NewConverter()
			if err != nil {
				log.Fatal(err)
			}
			defer converter.Destroy()

			// Add object to the converter.
			converter.Add(object)
			converter.Title = url
			converter.PaperSize = pdf.A4

			// Perform the conversion.
			return converter.Run(out) // Must be called on the main thread.
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
