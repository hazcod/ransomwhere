package file

import "os"

// GetHomeDirectoryWithFallback will return the current users home directory with a fallback to the system drive
func GetHomeDirectoryWithFallback() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return getSystemRootPath()
	}

	return homeDir
}
