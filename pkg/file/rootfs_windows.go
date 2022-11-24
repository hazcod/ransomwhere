package file

func getSystemRootPath() string {
	return os.Getenv("SystemDrive")
}
