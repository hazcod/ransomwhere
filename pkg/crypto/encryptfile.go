package crypto

import (
	"bytes"
	"errors"
	"filippo.io/age"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func EncryptFile(l *logrus.Entry, path, cryptExtension string, keys []age.Recipient, delete bool) error {
	cryptPath := path + cryptExtension

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

	if _, err := writer.Write(origBytes); err != nil {
		return fmt.Errorf("could not write crypted file %s: %v", path, err)
	}

	if err := os.WriteFile(cryptPath, out.Bytes(), 0600); err != nil {
		return fmt.Errorf("could not write crypted file %s: %v", path, err)
	}

	l.WithField("crypt_path", cryptPath).Debug("encrypted file")

	if delete {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("could not delete %s: %v", path, err)
		} else {
			l.Debug("deleted original")
		}
	}

	return nil
}
