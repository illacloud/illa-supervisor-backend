package controller

import (
	"github.com/illacloud/illa-supervisior-backend/src/model"
	"github.com/illacloud/illa-supervisior-backend/src/utils/tokenvalidator"
)

type Controller struct {
	Storage               *model.Storage
	Drive                 *model.Drive
	RequestTokenValidator *tokenvalidator.RequestTokenValidator
}

func NewController(storage *model.Storage, drive *model.Drive, validator *tokenvalidator.RequestTokenValidator) *Controller {
	return &Controller{
		Storage:               storage,
		Drive:                 drive,
		RequestTokenValidator: validator,
	}
}
