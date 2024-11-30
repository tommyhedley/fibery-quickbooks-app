package main

import (
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, struct{}{})
}

func TransformHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, struct{}{})
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, struct{}{})
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, struct{}{})
}
