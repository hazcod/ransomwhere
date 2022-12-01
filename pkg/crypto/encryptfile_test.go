package crypto

import (
	"errors"
	"filippo.io/age"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func TestEncryptFileWithDelete(t *testing.T) {
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

	if err := EncryptFile(l.WithField("", ""), tmpFile.Name(), []age.Recipient{id.Recipient()}, true); err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	if fileExists(tmpFile.Name()) {
		t.Fatal("file was NOT removed while delete = false")
	}
}

func TestEncryptFileWithoutDelete(t *testing.T) {
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

	if err := EncryptFile(l.WithField("a", ""), tmpFile.Name(), []age.Recipient{id.Recipient()}, false); err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	if !fileExists(tmpFile.Name()) {
		t.Fatal("file was removed while delete = false")
	}
}
