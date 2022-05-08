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

type result struct {
	text []byte
	err  error
}

func (a *application) handleUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	start := time.Now()
	defer func(t time.Time) {
		log.Printf("since: client: %s - %v\n", time.Since(t), r.RemoteAddr)
	}(start)

	handle, err := tesseract.TessBaseAPICreateWithMonitor("tur+eng")
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.ParseMultipartForm(10 << 20)

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

	ch := make(chan result)
	go func() {
		defer close(ch)
		fmt.Printf("fsize in handler: %v\n", len(f))
		texts, err := handle.ProcessImageMem(f)
		if err != nil {
			ch <- result{text: nil, err: err}
			return
		}

		log.Printf("path: %q\n", r.URL.Path)
		d := data{Content: texts}
		jsonResp, err := json.Marshal(d)
		if err != nil {
			ch <- result{text: nil, err: err}
			return
		}

		ch <- result{text: jsonResp, err: err}
	}()

	select {
	case <-ctx.Done():
		log.Println("context done")
		handle.SetCancelFunc()
		<-ch // on cancellation, wait for ProcessImageMem to finish
		handle.End()
		handle.Delete()
	case r := <-ch:
		if r.err != nil {
			log.Println("error ProcessImageMem: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(r.text)
	}
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
	addr := flag.String("addr", ":8080", "http network address")
	app := &application{
		// handle: handle,
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	log.Printf("starting server on %s\n", srv.Addr)
	log.Fatalln(srv.ListenAndServe())
}
