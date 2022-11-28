package crypto

import (
	"bytes"
	"errors"
	"filippo.io/age"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

// shouldDelete deletes a file on path if required by shouldDelete
func shouldDelete(l *logrus.Entry, shouldDelete bool, path string) {
	if !shouldDelete {
		return
	}

	l.WithField("path", path).Debug("deleting")

	if err := os.Remove(path); err != nil {
		l.WithError(err).WithField("path", path).Error("could not delete")
		return
	}

	l.Debug("deleted original")
}

// EncryptFile encrypts a file on the given path using the Recipient public key while optionally removing the original
func EncryptFile(l *logrus.Entry, path string, keys []age.Recipient, delete bool) error {
	cryptPath := path + CryptedExtension

	// if the encrypted version already exists, skip decryption and potentially delete the original
	if _, err := os.Stat(cryptPath); !errors.Is(err, os.ErrNotExist) {
		shouldDelete(l, delete, path)
		return nil
	}

	out := &bytes.Buffer{}

	// initialize the encryption writer
	writer, err := age.Encrypt(out, keys...)
	if err != nil {
		return fmt.Errorf("could setup encrypter: %v", err)
	}
	defer writer.Close()

	// read out all the file bytes, meaning encrypt
	origBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read file %s: %v", path, err)
	}

	// copy over the original file permissions
	filePermissions := os.FileMode(0600)
	fileInfo, err := os.Stat(path)
	if err != nil {
		l.WithError(err).WithField("path", path).Error("could not retrieve file permissions, using 0600")
	} else {
		filePermissions = fileInfo.Mode()
	}

	// use our encryption writer to encrypt it all
	if _, err := writer.Write(origBytes); err != nil {
		return fmt.Errorf("could not write crypted file %s: %v", path, err)
	}

	// write the last encrypted file chunk, else we'll get an unexpected EOF during decryption
	_ = writer.Close()

	// write our encrypted form to the crypted path using the original permissions
	if err := os.WriteFile(cryptPath, out.Bytes(), filePermissions); err != nil {
		return fmt.Errorf("could not write crypted file %s: %v", path, err)
	}

	l.WithField("crypt_path", cryptPath).Debug("encrypted file")

	// if required delete the original file
	shouldDelete(l, delete, path)

	return nil
}
