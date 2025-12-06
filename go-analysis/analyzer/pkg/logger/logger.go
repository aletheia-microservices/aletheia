package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

type CustomFormatter struct{}

func init() {
	Logger.SetLevel(logrus.DebugLevel)
	Logger.SetFormatter(&CustomFormatter{})
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

func getColor(level logrus.Level) int {
	switch level {
	case logrus.TraceLevel:
		return 37 // gray
	case logrus.DebugLevel:
		return 32 // green
	case logrus.InfoLevel:
		return 36 // cyan
	case logrus.WarnLevel:
		return 33 // yellow
	case logrus.ErrorLevel:
		return 31 // red
	case logrus.FatalLevel, logrus.PanicLevel:
		return 35 // magenta
	default:
		return 0 // default color
	}
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var file string
	var line int
	var ok bool
	i := 1

	// iterate over the stack frames to find the correct file and line
	for i < 20 {
		_, file, line, ok = runtime.Caller(i)
		if !strings.Contains(file, "entry.go") && !strings.Contains(file, "logger.go") {
			break
		}
		i++
	}

	if !ok {
		file = "???"
		line = 0
	} else {
		file = strings.TrimPrefix(file, "pkg/")

		dir := filepath.Dir(file)
		splits := strings.Split(dir, "/")
		len := len(splits)
		if len >= 2 {
			file = splits[len-2] + "/" + splits[len-1] + "/" + filepath.Base(file)
		} else if len >= 1 {
			file = splits[len-1] + "/" + filepath.Base(file)
		} else {
			file = filepath.Base(file)
		}
		// remove the "pkg" from the path if it exists
		file = strings.TrimPrefix(file, "pkg/")
	}

	levelColor := getColor(entry.Level)

	// no color
	//levelText := fmt.Sprintf("%-5s", strings.ToUpper(entry.Level.String()[:4]))
	// normal
	levelText := fmt.Sprintf("\x1b[%dm%-5s\x1b[0m", levelColor, strings.ToUpper(entry.Level.String()[:4]))
	// bold
	// levelText := fmt.Sprintf("\x1b[1m\x1b[%dm%-5s\x1b[0m", levelColor, strings.ToUpper(entry.Level.String()[:4]))

	// calculate padding to align the level text
	fileLine := fmt.Sprintf("%s:%d", file, line)
	padding := 40 - len(fileLine)
	if padding < 1 {
		padding = 1
	}

	// check if the message starts with [...] and apply color and bold if true
	/* msgText := entry.Message
    if strings.HasPrefix(msgText, "[") {
        endIdx := strings.Index(msgText, "]")
        if endIdx != -1 {
            // apply color to the text within brackets
            msgText = fmt.Sprintf("\x1b[%dm%s\x1b[0m%s", levelColor, msgText[:endIdx+1], msgText[endIdx+1:])
        }
    } */

	// check if the message starts with [...] (including cases where it has tabs before, e.g. \t\t\t[...]) and apply color and bold if true
	msgText := entry.Message
	trimmedMsgText := strings.TrimLeft(msgText, " \t")
	if strings.HasPrefix(trimmedMsgText, "[") {
		endIdx := strings.Index(trimmedMsgText, "]")
		if endIdx != -1 {
			// find the original position of the end index in msgText
			originalEndIdx := len(msgText) - len(trimmedMsgText) + endIdx + 1
			// apply color to the text within brackets
			msgText = fmt.Sprintf("%s\x1b[%dm%s\x1b[0m%s", msgText[:len(msgText)-len(trimmedMsgText)], levelColor, msgText[len(msgText)-len(trimmedMsgText):originalEndIdx], msgText[originalEndIdx:])
		}
	}

	// format the final message
	msg := fmt.Sprintf("%s%s%s %s\n", fileLine, strings.Repeat(" ", padding), levelText, msgText)
	return []byte(msg), nil
}
