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
	// try opening the encrypted file
	cryptedFile, err := os.Open(cryptedPath)
	if err != nil {
		return fmt.Errorf("could not read file: %v", err)
	}
	defer cryptedFile.Close()

	// initliase the decrypter writer
	reader, err := age.Decrypt(cryptedFile, identity)
	if err != nil {
		return fmt.Errorf("could setup encrypter: %v", err)
	}

	// copy over the same file permissions
	filePermissions := os.FileMode(0600)
	fileInfo, err := os.Stat(cryptedPath)
	if err != nil {
		l.WithError(err).WithField("path", cryptedPath).Error("could not retrieve file permissions, using 0600")
	} else {
		filePermissions = fileInfo.Mode()
	}

	// read everything in, meaning decrypt everything
	origBytes, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("could not read decrypted bytes: %v (size %d)", err, fileInfo.Size())
	}

	// get the original file path to restore to
	origPath := strings.TrimSuffix(cryptedPath, CryptedExtension)

	// write the decrypted bytes back to the original location with the original permissions
	if err := os.WriteFile(origPath, origBytes, filePermissions); err != nil {
		return fmt.Errorf("could not write crypted file %s: %v", cryptedPath, err)
	}

	l.WithField("path", origPath).Debug("decrypted file")

	// delete the encrypted version if necessaru
	shouldDelete(l, delete, cryptedPath)

	return nil
}
