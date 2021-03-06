package service

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/find"
	"github.com/lewis-od/wavelength/internal/io"
)

type UpdateService interface {
	Run(version string, lambdaNames []string, roleToAssume *builder.Role)
}

func NewUpdateService(
	finder find.Finder,
	updater builder.Updater,
	printer io.Printer,
	projectName *string) UpdateService {
	return &updateService{
		finder:      finder,
		updater:     updater,
		printer:     printer,
		projectName: projectName,
	}
}

type updateService struct {
	finder      find.Finder
	updater     builder.Updater
	printer     io.Printer
	projectName *string
}

func (u *updateService) Run(version string, lambdaNames []string, roleToAssume *builder.Role) {
	lambdasToUpdate, err := u.finder.FindLambdas(lambdaNames)
	if err != nil {
		u.printer.PrintErr(err)
		return
	}

	artifactBucket, err := u.finder.FindArtifactBucketName()
	if err != nil {
		u.printer.PrintErr(err)
		return
	}

	for _, lambda := range lambdasToUpdate {
		artifactLocation := fmt.Sprintf("%s/%s.zip", version, lambda)

		u.printer.Printlnf("⬆️ Updating %s with code at s3://%s/%s", lambda, artifactBucket, artifactLocation)

		err := u.updater.UpdateCode(lambda, artifactBucket, artifactLocation, roleToAssume)
		if err != nil {
			u.printer.PrintErr(err)
			return
		}
		u.printer.Printlnf("✅ Successfully updated %s", lambda)
	}
}
