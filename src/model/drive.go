package model

import (
	"go.uber.org/zap"
)

type S3Instance interface {
	GetPreSignedPutURL(fileName string) (string, error)
}

type Drive struct {
	logger                *zap.SugaredLogger
	SystemDriveS3Instance S3Instance
	TeamDriveS3Instance   S3Instance
}

func NewDrive(systemDriveS3Instance S3Instance, teamDriveS3Instance S3Instance, logger *zap.SugaredLogger) *Drive {
	return &Drive{
		logger:                logger,
		SystemDriveS3Instance: systemDriveS3Instance,
		TeamDriveS3Instance:   teamDriveS3Instance,
	}
}
