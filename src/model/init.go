package model

import (
	"github.com/illacloud/illa-supervisior-backend/src/utils/config"
)

var Config *config.Config

func init() {
	var err error
	Config, err = config.GetConfig()
	if err != nil {
		panic(err)
	}
}
