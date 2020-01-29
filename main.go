package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	dir = "access"
)

func main() {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("/images/"))))
	http.HandleFunc("/upload", uploadHandler)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("assets"))))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// POST takes the uploaded file(s) and saves it to disk.
	case "POST":
		// parse the multipart form in the request
		err := r.ParseMultipartForm(100000)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// get a ref to the parsed multipart form
		m := r.MultipartForm

		// get the *fileheaders
		files := m.File["file"]
		for i, fi := range files {

			// for each fileheader, get a handle to the actual file
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// create destination file making sure the path is writeable.
			// dst, err := os.Create(files[i].Filename)
			dst, err := os.Create(filepath.Join(dir, filepath.Base(fi.Filename)))
			defer dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// copy the uploaded file to the destination file
			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			_, _ = w.Write([]byte(fi.Filename))
		}
		// display success message.
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
