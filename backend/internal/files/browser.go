// files package, like everything elser in the files dir
package files

// imports
import (
	"os" // System files
	"path/filepath" // File paths
	"time" // time
)

// Defining File Structure
type FileItem struct {
	Name string `json:"name"`
	Path string `json:"path"`
	IsDir bool  `json:"isDir"`
	Size int64 `json:"size"`
	ModTime time.Time `json:"modTime"`
	MIMEType string `json:"mimeType"` 
}

// Function to list all directory contents
func ListDirectory(dirPath string, basePath string) ([]FileItem, error) {
	entries, err := os.ReadDir(dirPath) // Returns dir content and errors respectively

	if err != nil {
		return nil, err // Return error if any
	}

	items := []FileItem{} // Slice (array) to hold file items

	// Loop through each entry in the directory
	for _, entry := range entries {
		// Gets filepath import alongside file name for path
		fullPath := filepath.Join(dirPath, entry.Name()) // Joins dir path and entry name to get full path to the file in question. Ex: C:\Users\Username\Documents\ + file.txt = C:\Users\Username\Documents\file.txt
		
		// Symlink directory handling
		info, err := os.Stat(fullPath) // Symlink junction etc handling
		if err != nil {
			continue
		}

		isDir := info.IsDir() // Checks if the entry is a directory or not

		// Remove faux directories (symlinks that point outside the base path)
		if isDir {
			// test if path is a valid directory by trying to read it
			_, err := os.ReadDir(fullPath)
			// if no error, continue with it
			if err != nil {
				continue
			}

		}

		// compute relative path from basePath
        relativePath, err := filepath.Rel(basePath, fullPath)
        if err != nil {
            relativePath = entry.Name()
        }

		item := FileItem {
			Name: entry.Name(),
			Path: filepath.ToSlash(relativePath), // Convert to slash for consistency across OS
			IsDir: isDir,
			Size: info.Size(),
			ModTime: info.ModTime(),
			MIMEType: getMIMEType(entry.Name()),
		}

		items = append(items, item) // Append the item to the array holding all file items
	}

	return items, nil // Return the array of file items and nil for error
}

// Helper function to get MIME type based on file extension
func getMIMEType(filename string) string {
	ext := filepath.Ext(filename) // The file extension, ex: .txt, .jpg, etc.

	// Switch to assign correct MIME type from file extension
	switch ext {
	case ".mp4", ".avi", ".mkv":
		return "video"
	case ".mp3", ".wav":
		return "audio"
	case ".jpg", ".png", ".gif", ".jpeg", ".webp":
		return "image"
	case ".pdf":
		return "pdf"
	case ".txt", ".md":
		return "text"
	default:
		return "other"
	}
}
