package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/maliceio/malice/commands"
)

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

// Index root route
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

// MaliceScan scan a file
func MaliceScan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	file := vars["file"]

	if r.Body == nil {
		http.Error(w, "Body must be set", http.StatusBadRequest)
		return
	}

	// I didn't find what kind of body would give an error, but let's write the error code anyway.
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If S3 has any problem with our request, it will return an error. It could also be a 500 error if S3 itself has a problem.
	path := mux.Params(r).ByName("path")
	err = c.bucket.Put(path, content, r.Header.Get("Content-Type"), s3.PublicRead)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url := c.bucket.URL(path)
	w.Header().Set("Location", url)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"url":"%s"}`, url)
}

// MaliceLookUp lookup hash/url in intel sources
func MaliceLookUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	hashOrURL := vars["hashOrURL"]

	err := commands.APILookUp(hashOrURL)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	// w.Write([]byte(outPut))
}
