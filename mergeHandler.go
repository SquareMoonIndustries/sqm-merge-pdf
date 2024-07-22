package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func init() {
	routes = append(routes, Route{"mergeHandler", "POST", "/merge", mergeHandler})
}

func mergeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("secret") != settings.Secret {
		http.Error(w, "Security lockdown sector 4", http.StatusForbidden)
		logger.Info("forbidden")
		return
	}

	// Parse the multipart form data
	// 32Mb in memory
	err := r.ParseMultipartForm(32 << 20) // 32Mb
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the pdf form field array
	multipartFormData := r.MultipartForm
	filename := r.FormValue("filename")
	if filename == "" {
		filename = "tmp.pdf"
	}
	//logger.Info(filename)
	handleMultipartFormData(multipartFormData, w, filename)
}

func handleMultipartFormData(multipartFormData *multipart.Form, w http.ResponseWriter, filename string) {
	// Create a slice to store io.ReadSeeker instances
	var fileReaders []io.ReadSeeker

	for _, v := range multipartFormData.File["pdf"] {
		// Open the file
		file, err := v.Open()
		if err != nil {
			//fmt.Println("Error opening file:", err)
			//continue
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Info("Error opening file: ", err)
		}
		defer file.Close()

		// Read the first 512 bytes to detect the file type
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil && err != io.EOF {
			//fmt.Println("Error reading file:", err)
			//continue
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error(err)
			return
		}

		// Detect the content type
		contentType := http.DetectContentType(buffer)
		if contentType != "application/pdf" {
			http.Error(w, "Not a pdf file was sent", http.StatusBadRequest)
			logger.Info("Not a PDF file was sent")
			return
		}
		//fmt.Println("Detected content type:", contentType)

		// Reset the file reader to the beginning
		file.Seek(0, io.SeekStart)

		// Read the file contents into a byte slice
		content, err := io.ReadAll(file)
		if err != nil {
			//fmt.Println("Error reading file:", err)
			//continue
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error(err)
			return
		}

		// Create a bytes.Reader from the byte slice and append it to the slice
		fileReaders = append(fileReaders, bytes.NewReader(content))

		//fmt.Println(v.Filename, ":", v.Size)
	}

	// Define the PDF processing configuration
	config := model.NewDefaultConfiguration()

	// Add headers
	w.Header().Add("content-type", "application/pdf")
	w.Header().Add("Content-Disposition", "attachment; filename=\""+filename+"\"")

	// Merge the files
	err := api.MergeRaw(fileReaders, w, false, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		//fmt.Println("Error merging PDF files:", err)
		logger.Info(err)
		return
	}
}
