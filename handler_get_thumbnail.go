package main

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerThumbnailGet(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid video ID", err)
		return
	}

	entries, err := os.ReadDir(cfg.assetsRoot)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error reading assets directory", err)
		return
	}

	prefix := videoID.String() + "."
	var foundPath string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, prefix) {
			foundPath = filepath.Join(cfg.assetsRoot, name)
			break
		}
	}

	if foundPath == "" {
		respondWithError(w, http.StatusNotFound, "Thumbnail not found", nil)
		return
	}

	data, err := os.ReadFile(foundPath)
	if err != nil {
		if os.IsNotExist(err) {
			respondWithError(w, http.StatusNotFound, "Thumbnail not found", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error reading thumbnail", err)
		return
	}

	ext := filepath.Ext(foundPath)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(data)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error writing response", err)
		return
	}
}
