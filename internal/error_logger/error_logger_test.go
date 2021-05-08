package error_logger_test

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/error_logger"
	"github.com/lewis-od/wavelength/internal/testutil/mock_filesystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type fixedClock struct {
	t time.Time
}

func (f *fixedClock) Now() time.Time {
	return f.t
}

const expectedReport = `Errors encountered during build at 2021-05-08T15:04:05Z
Lambda: lambda-one
Go error: this is an error
Build output:
building lambda-one... done

Lambda: lambda-two
Go error: this is an error
Build output:
building lambda-two... done

--------------------------------------------------------------------------------
`

func TestErrorLogger_WriteLogFile(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2021-05-08T15:04:05Z")

	fileName := "wavelength.log"
	mockClock := &fixedClock{t: now}
	var mockFilesystem *mock_filesystem.MockFilesystem
	var errorLogger error_logger.ErrorLogger

	setupTest := func() {
		mockFilesystem = new(mock_filesystem.MockFilesystem)
		errorLogger = error_logger.NewErrorLogger(mockFilesystem, mockClock, fileName)
	}

	t.Run("Success", func(t *testing.T) {
		setupTest()
		resultOne := &builder.BuildResult{
			LambdaName: "lambda-one",
			Error:      fmt.Errorf("this is an error"),
			Output:     []byte("building lambda-one... done"),
		}
		resultTwo := &builder.BuildResult{
			LambdaName: "lambda-two",
			Error:      fmt.Errorf("this is an error"),
			Output:     []byte("building lambda-two... done"),
		}

		mockFilesystem.On("AppendToFile", fileName, expectedReport).Return(nil)

		errorLogger.AddError(resultOne)
		errorLogger.AddError(resultTwo)
		err := errorLogger.WriteLogFile()

		assert.Nil(t, err)
		mockFilesystem.AssertExpectations(t)
	})
	t.Run("Error", func(t *testing.T) {
		setupTest()

		expectedErr := fmt.Errorf("there was an error")
		mockFilesystem.On("AppendToFile", fileName, mock.Anything).Return(expectedErr)

		err := errorLogger.WriteLogFile()
		assert.Equal(t, expectedErr, err)
	})

}
