// File serving, STREAMING files

// api package used for all api dir files
package api

// imports
import (
	"net/http" // http handling
	"os" // os file handling
	"path/filepath" // filepath for file paths
	"homeserver/internal/db" // database handling for logging
)

// Function to stream files to client, w = ResponseWriter r = Request
func StreamFileHandler(w http.ResponseWriter, r *http.Request, basePath string) {
	// Get file path from query parameter
	subPath := r.URL.Query().Get("path")
	if subPath == "" {
		http.Error(w, "No file path given", http.StatusBadRequest) // Bad request, no path provided
		return
	}

	fullPath := filepath.Join(basePath, subPath)

	// Open file in read only with os.Open, which gives a *os.File(os.Open does) and uses io.ReadSeeker for reading and seeking for videos. Equivalent to var file *os.File, file,, err = os.Open(path)
	file, err := os.Open(fullPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound) // File not found
		return
	}

	// Defer to close opened file
	defer file.Close()

	// Use .Stat to get file metadata
	info, err := file.Stat()
	if err != nil {
		http.Error(w, "Couldn't get file info", http.StatusInternalServerError)
		return
	}

	// Log streaming activity
	_ = db.LogActivity("stream", fullPath, r.Header.Get("X-User-ID"), r.Header.Get("X-Device-ID"), r.RemoteAddr)

	// Serve content with http.ServeContent from net/http, auto handles content-type header, content-length, range requests (bytes x to y where user is rather than whole file), and more. 
	// w = ResponseWriter, r = Request, info.Name() = file name for content-disposition header, info.ModTime() = file mod time for caching, file = the file itself opened with os.Open as an io.ReadSeeker for streaming
	http.ServeContent(w, r, info.Name(), info.ModTime(), file)
}
