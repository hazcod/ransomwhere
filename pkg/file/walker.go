package file

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"path/filepath"
)

func WalkFiles(directory string, l *logrus.Logger, f func(name, filePath string) error) error {
	dirFS := os.DirFS(directory)

	if err := fs.WalkDir(dirFS, ".", func(relPath string, d fs.DirEntry, err error) error {
		path := filepath.Join(directory, relPath)

		logger := l.WithField("path", path)

		if err != nil {
			logger.WithError(err).Warnf("error while walking %s", path)
			return nil
		}

		if d.IsDir() {
			logger.Trace("skipping walk of directory")
			return nil
		}

		logger.Trace("walking")

		if err := f(d.Name(), path); err != nil {
			logger.WithError(err).Warn("function returned error")
		} else {
			logger.Trace("function completed successfully")
		}

		return nil

	}); err != nil {
		return fmt.Errorf("could not walk "+directory+": %v", err)
	}

	return nil
}
