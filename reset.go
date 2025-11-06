package main

import (
	"net/http"
	"os"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err := cfg.DB.DeleteAllUsers(r.Context())
	if err != nil {
		http.Error(w, "Could not reset database", http.StatusInternalServerError)
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
