package file

import "os"

func GetHomeDirectoryWithFallback() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return getSystemRootPath()
	}

	return homeDir
}
