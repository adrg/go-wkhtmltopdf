package main

/*
Test payloads:

{
    "url": "http://www.google.com",
	"converterOptions": {
	    "margin.left": "10mm",
	    "margin.right": "10mm",
	    "margin.top": "10mm",
	    "margin.bottom": "10mm"
	},
	"objectOptions": {
	}
}

{
    "url": "sample1.html",
	"converterOptions": {
	    "margin.left": "10mm",
	    "margin.right": "10mm",
	    "margin.top": "10mm",
	    "margin.bottom": "10mm"
	},
	"objectOptions": {
		"load.windowStatus": "ready"
	}
}
*/

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	pdf "github.com/leandrosilva/go-wkhtmltopdf" // <- This is what I use on my machine
	// pdf "github.com/adrg/go-wkhtmltopdf"         <- You may want to use this instead
)

func startHTTPServer() error {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", pdf.ConvertPostHandler).Methods("POST")

	httpPort := "7070"
	log.Println("HTTP Server listening on", httpPort)

	server := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        muxRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 2,
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func main() {
	go startHTTPServer()
	pdf.StartConvertLoop()
}
