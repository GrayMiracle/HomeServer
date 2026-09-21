// Backend file handling

// api package used by all files in the api dir
package api


// imports
import (
	"encoding/json" // JSON encoding/decoding
	"net/http" // HTTP handling
	"homeserver/internal/files" // File handling package
	"os" // System for file handling
	"io" // Input Output for uploads
	"path/filepath" // Upload filepaths
	"homeserver/internal/db" // Database handling for logging
)



// ALL API Handler functions



// Function to list files from mentioned dir. ResponseWriter w is to write response back to client, Request r is the incoming request from client
func ListFilesHandler (w http.ResponseWriter, r *http.Request, basePath string) {
	// Only provide files from specified backend dir and nothing else
	subPath:= r.URL.Query().Get("path") // Get path from URL query parameters. ex. http://localhost:3000/files?path=C:\Users\Username\Documents

	// assign fullpath, if no subpath use that directly
	var fullPath string
	if subPath == "" {
		fullPath = basePath
	} else { // if there is subpath join with basepath
		fullPath = filepath.Join(basePath, subPath)
	}

	items, err := files.ListDirectory(fullPath, basePath) // Use imported ListDirectory from files package using path as argument
	if err != nil {
		http.Error(w, "Couldn't read dir", http.StatusInternalServerError) // If error break here
		return
	}

	w.Header().Set("Content-Type", "application/json") // JSON response header
	json.NewEncoder(w).Encode(items) //Encode items to JSON for global compatibility and send as response through w as a ResponseWriter
}

// function to download files
func DownloadFileHandler(w http.ResponseWriter, r *http.Request, basePath string) {
	subPath := r.URL.Query().Get("path") // Get subpath for file to download

	if subPath == "" {
		http.Error(w, "No file to download selected", http.StatusBadRequest) // Bad request, no path provided
		return
	}

	fullPath := filepath.Join(basePath, subPath)

	intent := r.URL.Query().Get("intent")
	eventType := "download"
	if intent == "view" {
		eventType = "view"
	}

	// Log download
	_ = db.LogActivity(eventType, fullPath, r.Header.Get("X-User-ID"), r.Header.Get("X-Device-ID"), r.RemoteAddr)

	http.ServeFile(w, r, fullPath) // Transfer file to client using builtin fn ServeFile, which handles file reading and response writing
}

// function to delete files
func DeleteFileHandler(w http.ResponseWriter, r *http.Request, basePath string) {
	// Endpoint restriction
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request", http.StatusMethodNotAllowed) // Cannot delete from other requests
		return
	}

	// Get subpath from URL query. ex. \file.txt
	subPath := r.URL.Query().Get("path")
	if subPath == "" {
		http.Error(w, "No path given", http.StatusBadRequest)
		return
	}

	// Full path combining base in config and subpath for specific file
	fullPath := filepath.Join(basePath, subPath)

	// Delete file with os.Remove with err handling
	err := os.Remove(fullPath) // Remove files at the given path
	if err != nil {
		http.Error(w, "Error deleting file", http.StatusInternalServerError) // Could not delete file for whatever reason
		return
	}

	// Log deletion
	_ = db.LogActivity("delete", fullPath, r.Header.Get("X-User-ID"), r.Header.Get("X-Device-ID"), r.RemoteAddr)

	w.WriteHeader(http.StatusOK) // 200 Ok if success
	w.Write([]byte("File deleted")) // Converting string to slice of bytes as response body for http network level compatibility
}

// function to upload files
func UploadFileHandler(w http.ResponseWriter, r *http.Request, basePath string) {
	// Only allow POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request", http.StatusMethodNotAllowed) // No other methods
		return
	}

	// ParseMultipartForm for multipart uploads, the standard encoding for file uploads in browsers.

	// 10 << 30. 10 is the number, 30 is the bits it gets shifted to the left, making 10 * 2^30 for 10 GB. If it were 10 << 20, it would be 10 * 2^20 for 10 MB.
	err := r.ParseMultipartForm(10 << 30) // 10 GB limit
	if err != nil {
		http.Error(w, "Over limit", http.StatusBadRequest) // When upload is over limit
		return
	}

	// FormFile gets uploaded file from form, "file" is the form field name in the multipart form data the frontend must send the file under. ex. <input type="file" name="file" />
	// file = file content, handler = file metadata, err = error
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File could not be read", http.StatusBadRequest) // File couldn't be read
		return
	}

	// defer to definitely run when fn exits
	defer file.Close() // Close to free resources and prevent leaking file handles

	// Dest path from config file. ex. http://localhost:3000/upload?path=C:\Users\Megas\Documents
	destPath := r.URL.Query().Get("path")

	// Use base config dest path as default, and add mentioned filepath to it if given
	var fullDestPath string
	if destPath == "" {
		fullDestPath = basePath
	} else {
		fullDestPath = filepath.Join(basePath, destPath)
	}

	// Join filepath to save file under the dir. ex. C:\Users\Megas\Documents + file.txt = C:\Users\Megas\Documents\file.txt
	// No string concatenation since different OSes have different path separators, so filepath.Join handles that. ie / for Linux and \ for Windows
	dst, err := os.Create(filepath.Join(fullDestPath, handler.Filename)) // handler.Filename is a default from handler file metadata's struct, and os.Create creates a file at the given argument
	// dst is the file the uploaded content should be written to, err is error handling
	if err != nil {
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}
	// defer to close again, freeing resources and preventing leaks
	defer dst.Close()

	// Using input output to copy uploaded file to destination file at the server's filesystem
	io.Copy(dst, file)

	// Log uploads
	_ = db.LogActivity("upload", filepath.Join(fullDestPath, handler.Filename), r.Header.Get("X-User-ID"), r.Header.Get("X-Device-ID"), r.RemoteAddr)

	// Client response at the end
	w.WriteHeader(http.StatusOK) // 200 if success
	w.Write([]byte("File upload success")) // String to byte response body

}