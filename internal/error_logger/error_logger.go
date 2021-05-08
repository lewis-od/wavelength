package error_logger

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/clock"
	"github.com/lewis-od/wavelength/internal/io"
	"strings"
	"time"
)

type WavelengthError struct {
	Lambda string
	Output []byte
}

type ErrorLogger interface {
	AddError(result *WavelengthError)
	WriteLogFile() error
}

func NewErrorLogger(fs io.Filesystem, c clock.Clock, logFileName string) ErrorLogger {
	return &errorLogger{
		fs:          fs,
		c:           c,
		logFileName: logFileName,
		errors:      make([]*WavelengthError, 0, 10),
	}
}

type errorLogger struct {
	fs          io.Filesystem
	c           clock.Clock
	logFileName string
	errors      []*WavelengthError
}

func (el *errorLogger) AddError(result *WavelengthError) {
	el.errors = append(el.errors, result)
}

func (el *errorLogger) WriteLogFile() error {
	var logBuilder strings.Builder
	now := el.c.Now()

	logBuilder.WriteString(fmt.Sprintf("Errors encountered during build at %s\n", now.Format(time.RFC3339)))
	for _, wavelengthError := range el.errors {
		logBuilder.WriteString(fmt.Sprintf("Lambda: %s\n", wavelengthError.Lambda))
		logBuilder.WriteString("Build output:\n")
		logBuilder.Write(wavelengthError.Output)
		logBuilder.WriteString("\n\n")
	}
	for i := 0; i < 80; i++ {
		logBuilder.WriteString("-")
	}
	logBuilder.WriteString("\n")

	return el.fs.AppendToFile("wavelength.log", logBuilder.String())
}
