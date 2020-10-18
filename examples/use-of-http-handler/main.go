package main

/*
Test payloads:

{
	"converterOpts": {
	    "marginLeft": "10mm",
	    "marginRight": "10mm",
	    "marginTop": "10mm",
	    "marginBottom": "10mm",
        "outlineDepth": 0
	},
	"objectOpts": {
        "location": "http://www.google.com"
	}
}

{
	"converterOpts": {
	    "marginLeft": "10mm",
	    "marginRight": "10mm",
	    "marginTop": "10mm",
	    "marginBottom": "10mm",
        "outlineDepth": 0
	},
	"objectOpts": {
        "location": "sample1.html",
		"windowStatus": "ready"
	}
}*/

import (
	"log"
	"net/http"
	"runtime"
	"time"

	pdf "github.com/adrg/go-wkhtmltopdf"
	"github.com/gorilla/mux"
)

func init() {
	runtime.LockOSThread()
}

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
