package log

import (
	"bytes"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

var (
	// colorStatus returns a new function that returns status-colorized (cyan) strings for the
	// given arguments with fmt.Sprint().
	colorStatus = color.New(color.FgCyan).SprintFunc()

	// colorWarn returns a new function that returns status-colorized (yellow) strings for the
	// given arguments with fmt.Sprint().
	colorWarn = color.New(color.FgYellow).SprintFunc()

	// colorInfo returns a new function that returns info-colorized (green) strings for the
	// given arguments with fmt.Sprint().
	colorInfo = color.New(color.FgGreen).SprintFunc()

	// colorError returns a new function that returns error-colorized (red) strings for the
	// given arguments with fmt.Sprint().
	colorError = color.New(color.FgRed).SprintFunc()

	logger *logrus.Entry

	labelsPath = "/etc/labels"
)

var ( // For Test Mocks
	initLogger = initializeLogger
)

var defaultLogger *logrus.Logger

// FormatLayoutType the layout kind
type FormatLayoutType string

// VgsTextFormat lets use a custom text format
type VgsTextFormat struct {
	ShowInfoLevel   bool
	ShowTimestamp   bool
	TimestampFormat string
}

func NewVgsTextFormat() *VgsTextFormat {
	return &VgsTextFormat{
		ShowInfoLevel:   false,
		ShowTimestamp:   false,
		TimestampFormat: "2006-01-02 15:04:05",
	}
}

// Format formats the log statement
func (f *VgsTextFormat) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	level := strings.ToUpper(entry.Level.String())
	switch level {
	case "INFO":
		b.WriteString(colorInfo(level))
		b.WriteString(": ")
	case "WARNING":
		b.WriteString(colorWarn(level))
		b.WriteString(": ")
	case "DEBUG":
		b.WriteString(colorStatus(level))
		b.WriteString(": ")
	default:
		b.WriteString(colorError(level))
		b.WriteString(": ")
	}
	if f.ShowTimestamp {
		b.WriteString(entry.Time.Format(f.TimestampFormat))
		b.WriteString(" - ")
	}

	b.WriteString(entry.Message)

	if !strings.HasSuffix(entry.Message, "\n") {
		b.WriteByte('\n')
	}
	return b.Bytes(), nil
}

func initializeLogger() error {
	if logger == nil {
		var fields logrus.Fields
		logger = logrus.WithFields(fields)

		format := os.Getenv("VGS_LOG_FORMAT")
		if format == "json" {
			setFormatter("json")
		} else {
			setFormatter("text")
		}
	}
	return nil
}

// setFormatter sets the logrus format to use either text or JSON formatting
func setFormatter(layout FormatLayoutType) {
	switch layout {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.SetFormatter(NewVgsTextFormat())
	}
}

// Logger obtains the logger for use in the codebase
// This is the only way you should obtain a logger
func Logger() *logrus.Entry {
	err := initLogger()
	if err != nil {
		logrus.Warnf("error initializing logrus %v", err)
	}
	return logger
}

// SetLevel sets the logging level
func SetLevel(s string) error {
	level, err := logrus.ParseLevel(s)
	if err != nil {
		return errors.Errorf("Invalid log level '%s'", s)
	}
	Logger().Debugf("logging set to level: %s", level)
	logrus.SetLevel(level)
	return nil
}

// CaptureOutput calls the specified function capturing and returning all logged messages.
func CaptureOutput(f func()) string {
	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	f()
	logrus.SetOutput(os.Stderr)
	return buf.String()
}

// SetOutput sets the outputs for the default logger.
func SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}

// GetLevels returns the list of valid log levels
func GetLevels() []string {
	var levels []string
	for _, level := range logrus.AllLevels {
		levels = append(levels, level.String())
	}
	return levels
}
