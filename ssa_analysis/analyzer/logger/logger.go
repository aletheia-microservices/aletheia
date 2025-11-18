package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

type CustomFormatter struct{}

func init() {
	Logger.SetLevel(logrus.DebugLevel)
}

func SetOutput(output string) error {
	var writer io.Writer

	switch output {
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	default:
		file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file %s: %v", output, err)
		}
		writer = file
	}

	Logger.SetOutput(writer)
	return nil
}
