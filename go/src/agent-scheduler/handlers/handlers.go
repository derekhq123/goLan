package handlers

import (
	"fmt"
	"net/http"
)

const (
	version = "v1.0.0"
)

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, version)
}
