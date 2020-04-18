package pdf

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// ConvertOptions encapsulates wkhtmltopdf's converter and object options
type ConvertOptions struct {
	URL              string            `json:"url"`
	ConverterOptions map[string]string `json:"converterOptions"`
	ObjectOptions    map[string]string `json:"objectOptions"`
}

// ConvertRequestChannel receives request for convertion
var ConvertRequestChannel = make(chan ConvertOptions)

// ConvertResponseChannel delivers converted content
var ConvertResponseChannel = make(chan []byte)

// StopConvertLoopChannel informs stop signal
var StopConvertLoopChannel = make(chan bool)
var stopConvertLoop = false

// StartConvertLoop is the main thread loop listen to ConvertRequestChannel
// and feeding ConvertResponseChannel
func StartConvertLoop() {
	log.Println("Starting convert loop")

	Init()
	defer Destroy()

	go func() {
		<-StopConvertLoopChannel
		log.Println("Received StopConvertLoop signal")
		stopConvertLoop = true
	}()

	for !stopConvertLoop {
		log.Println("Waiting for convertion request...")
		options := <-ConvertRequestChannel
		log.Println("Received a convertion request for:", options.URL)

		content, err := convert(options)
		if err != nil {
			log.Println("Failed to convert:", options.URL)
			ConvertResponseChannel <- nil
		}

		log.Println("Sending PDF content:", options.URL)
		ConvertResponseChannel <- content
	}

	log.Println("Convert loop is over")
}

// StopConvertLoop sends a StopConvertLoop signal
func StopConvertLoop() {
	log.Println("Sending StopConvertLoop signal")
	StopConvertLoopChannel <- true
}

// ConvertPostHandler converts HTML to PDF based on payload options
func ConvertPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respondWithText(w, r, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	log.Println("--- BEGIN REQUEST ---")

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var options ConvertOptions
	if err := decoder.Decode(&options); err != nil {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			body = make([]byte, 0, 0)
		}
		respondWithText(w, r, http.StatusBadRequest, string(body))
		return
	}

	log.Println("Requesting for convert:", options.URL)
	ConvertRequestChannel <- options
	content := <-ConvertResponseChannel
	log.Println("Received response for:", options.URL)

	if content == nil {
		respondWithText(w, r, http.StatusInternalServerError, "Failed to convert file")
		return
	}
	respondWithPDF(w, r, http.StatusOK, content)

	log.Println("--- END REQUEST ---")
}

func respondWithText(w http.ResponseWriter, r *http.Request, statusCode int, payload string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(payload))
}

func respondWithPDF(w http.ResponseWriter, r *http.Request, statusCode int, payload []byte) {
	w.WriteHeader(statusCode)
	w.Header().Add("content-type", "application/pdf")
	w.Write(payload)
}

func convert(opt ConvertOptions) ([]byte, error) {
	// Create object from url
	object, err := NewObject(opt.URL)
	if err != nil {
		log.Println("Could not create object for", opt.URL)
		return nil, err
	}
	log.Println("Object URL:", opt.URL)

	// Add object options
	for k, v := range opt.ObjectOptions {
		object.SetOption(k, v)
	}

	// Create converter
	converter, err := NewConverter()
	if err != nil {
		log.Println("Could not create converter for", opt.URL)
		return nil, err
	}
	defer converter.Destroy()

	// Add created object to the converter
	converter.Add(object)

	// Add converter options
	for k, v := range opt.ConverterOptions {
		converter.SetOption(k, v)
	}

	// Convert the objects and get the output PDF document
	output := new(bytes.Buffer)
	err = converter.Run(output)
	if err != nil {
		log.Println("Could not convert object to PDF:", opt.URL)
		return nil, err
	}
	raw := output.Bytes()
	log.Println("PDF", len(raw), "bytes of size:", opt.URL)

	return raw, nil
}
