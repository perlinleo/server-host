package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

// func setupResponse(w *http.ResponseWriter, req *http.Request) {
// 	(*w).Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
// 	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
// 	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Set-Cookie")
// }

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// setupResponse(&w, r)
	// if (*r).Method == "OPTIONS" {
	// 	return
	// }

	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main() {
	mux := mux.NewRouter()

	spa := spaHandler{staticPath: "static", indexPath: "index.html"}
	mux.PathPrefix("/").Handler(spa)

	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"http://127.0.0.1:8080"},
	// 	AllowCredentials: true,
	// })
	// handler := c.Handler(mux)

	// headersOk := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin", "X-Requested-With", "Authorization"})
	// originsOk := handlers.AllowedOrigins([]string{"*"})
	// methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	srv := &http.Server{
		Handler: mux,
		// Handler: handler,
		// Handler:      handlers.CORS(headersOk, originsOk, methodsOk)(mux),
		Addr:         ":443",
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Fatal(srv.ListenAndServeTLS("../drip-ilyagu.com+4.pem", "../drip-ilyagu.com+4-key.pem"))
}
