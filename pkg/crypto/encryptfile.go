package crypto

import (
	"bytes"
	"errors"
	"filippo.io/age"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func EncryptFile(l *logrus.Entry, path string, keys []age.Recipient, delete bool) error {
	cryptPath := path + CryptedExtension

	if _, err := os.Stat(cryptPath); !errors.Is(err, os.ErrNotExist) {
		return nil
	}

	out := &bytes.Buffer{}

	writer, err := age.Encrypt(out, keys...)
	if err != nil {
		return fmt.Errorf("could setup encrypter: %v", err)
	}

	origBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read file %s: %v", path, err)
	}

	filePermissions := os.FileMode(0600)
	fileInfo, err := os.Stat(path)
	if err != nil {
		l.WithError(err).WithField("path", path).Error("could not retrieve file permissions, using 0600")
	} else {
		filePermissions = fileInfo.Mode()
	}

	if _, err := writer.Write(origBytes); err != nil {
		return fmt.Errorf("could not write crypted file %s: %v", path, err)
	}

	if err := os.WriteFile(cryptPath, out.Bytes(), filePermissions); err != nil {
		return fmt.Errorf("could not write crypted file %s: %v", path, err)
	}

	l.WithField("crypt_path", cryptPath).Debug("encrypted file")

	if delete {
		l.WithField("path", path).Debug("deleting")

		if err := os.Remove(path); err != nil {
			return fmt.Errorf("could not delete %s: %v", path, err)
		} else {
			l.Debug("deleted original")
		}
	}

	return nil
}
