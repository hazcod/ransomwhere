package crypto

import (
	"filippo.io/age"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"testing"
)

func TestDecryptFileWithoutDelete(t *testing.T) {
	bytes := make([]byte, 150)

	_, err := rand.Read(bytes)
	if err != nil {
		t.Fatal(err)
	}

	tmpFile, err := os.CreateTemp(os.TempDir(), "ransomwhere")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	if _, err := tmpFile.Write(bytes); err != nil {
		t.Fatal(err)
	}

	tmpFile.Close()

	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}

	l := logrus.New()
	l.SetLevel(logrus.FatalLevel)

	if err := EncryptFile(l.WithField("a", ""), tmpFile.Name(), []age.Recipient{id.Recipient()}, true); err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	if !fileExists(tmpFile.Name() + CryptedExtension) {
		t.Fatal("encrypted file does not exist")
	}

	if err := DecryptFile(l.WithField("", ""), tmpFile.Name()+CryptedExtension, id, true); err != nil {
		t.Fatalf("could not decrypt file: %v", err)
	}

	if fileExists(tmpFile.Name() + CryptedExtension) {
		t.Fatal("crypted file was not removed after decrypt")
	}

	if !fileExists(tmpFile.Name()) {
		t.Fatal("file was removed after decrypt while delete = false")
	}
}
