package error_logger

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/clock"
	"github.com/lewis-od/wavelength/internal/io"
	"strings"
	"time"
)

type ErrorLogger interface {
	AddError(result *builder.BuildResult)
	WriteLogFile() error
}

func NewErrorLogger(fs io.Filesystem, c clock.Clock, logFileName string) ErrorLogger {
	return &errorLogger{
		fs:          fs,
		c:           c,
		logFileName: logFileName,
		errors:      make([]*builder.BuildResult, 0, 10),
	}
}

type errorLogger struct {
	fs          io.Filesystem
	c           clock.Clock
	logFileName string
	errors      []*builder.BuildResult
}

func (el *errorLogger) AddError(result *builder.BuildResult) {
	el.errors = append(el.errors, result)
}

func (el *errorLogger) WriteLogFile() error {
	var logBuilder strings.Builder
	now := el.c.Now()

	logBuilder.WriteString(fmt.Sprintf("Errors encountered during build at %s\n", now.Format(time.RFC3339)))
	for _, buildResult := range el.errors {
		logBuilder.WriteString(fmt.Sprintf("Lambda: %s\n", buildResult.LambdaName))
		logBuilder.WriteString(fmt.Sprintf("Go error: %s\n", buildResult.Error))
		logBuilder.WriteString("Build output:\n")
		logBuilder.Write(buildResult.Output)
		logBuilder.WriteString("\n\n")
	}
	for i := 0; i < 80; i++ {
		logBuilder.WriteString("-")
	}
	logBuilder.WriteString("\n")

	return el.fs.AppendToFile("wavelength.log", logBuilder.String())
}
