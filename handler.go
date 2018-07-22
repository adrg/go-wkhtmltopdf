package pdf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var mutex = &sync.Mutex{}

// Options encapsulates wkhtmltopdf's converter and object options
type Options struct {
	URL              string            `json:"url"`
	ConverterOptions map[string]string `json:"converterOptions"`
	ObjectOptions    map[string]string `json:"objectOptions"`
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

	var options Options
	if err := decoder.Decode(&options); err != nil {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			body = make([]byte, 0, 0)
		}
		respondWithText(w, r, http.StatusBadRequest, string(body))
		return
	}

	log.Println("Waiting lock for", options.URL)
	mutex.Lock()
	log.Println("Locked")
	content, err := convert(options)
	log.Println("Unlocking...")
	mutex.Unlock()

	if err != nil {
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

func convert(opt Options) ([]byte, error) {
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
	converter := NewConverter()
	defer converter.Destroy()

	// Add created object to the converter
	converter.AddObject(object)

	// Add converter options
	for k, v := range opt.ConverterOptions {
		converter.SetOption(k, v)
	}

	// Convert the objects and get the output PDF document
	output, err := converter.Convert()
	if err != nil {
		log.Println("Could not convert object to PDF:", opt.URL)
		return nil, err
	}
	log.Println("PDF", len(output), "bytes of size:", opt.URL)

	return output, nil
}
