package main

import (
	"io"
	"net/http"
	"os"
)

func LogoHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("./logo.svg")
	if err != nil {
		http.Error(w, "Unable to open SVG file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	svgData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read SVG file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")

	w.Write(svgData)
}
