package main

import (
	"net/http"
)

func handleIndex(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("This is API backend for Fluffy Acorn. Have fun."))
}
