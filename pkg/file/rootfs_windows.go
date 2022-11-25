package file

import (
	"os"
)

func getSystemRootPath() string {
	return os.Getenv("SystemDrive")
}
