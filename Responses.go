package main

import "net/http"

func AccessDenied(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
}
