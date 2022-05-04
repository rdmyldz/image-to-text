package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/rdmyldz/i2t/tesseract"
)

type data struct {
	Content []string `json:"content,omitempty"`
}

const staticFiles = "../frontend/build/static/"
const indexPath = "../frontend/build/index.html"

var ts = template.Must(template.New("home").ParseFiles(indexPath))

type application struct {
	handle *tesseract.TessBaseAPI
}

func (a *application) handleHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("path: %q\n", r.URL.Path)
	err := ts.ExecuteTemplate(w, "index.html", nil)
	log.Printf("error: %v", err)
}

func (a *application) handleUpload(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func(t time.Time) {
		log.Printf("since: client: %s - %v\n", time.Since(t), r.RemoteAddr)
	}(start)

	handle, err := tesseract.TessBaseAPICreate("tur+eng")
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "POST, GET ")
	// w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization ")

	r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	f, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("error reading r.Body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("fsize in handler: %v\n", len(f))
	texts, err := handle.ProcessImageMem(f)
	if err != nil {
		log.Println("error ProcessImageMem: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("path: %q\n", r.URL.Path)
	d := data{Content: texts}
	jsonResp, err := json.Marshal(d)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func (app *application) routes() http.Handler {

	router := mux.NewRouter()
	router.HandleFunc("/upload", app.handleUpload).Methods("POST", "GET")
	router.HandleFunc("/", app.handleHome).Methods("GET")

	fileServer := http.FileServer(http.Dir(staticFiles))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	return router
}

// TODO: implement a worker pool
// for now, carried to create the tesseract instance into handleUpload
// since this way, we can't convert images to texts concurrently
func main() {
	// handle, err := tesseract.TessBaseAPICreate("tur+eng")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	addr := flag.String("addr", ":8080", "http network address")
	app := &application{
		// handle: handle,
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}
	// router :=mux.NewRouter()
	// router.HandleFunc("/", app.handleHome).Methods("GET")
	log.Printf("starting server on %s\n", srv.Addr)
	log.Fatalln(srv.ListenAndServe())
}
