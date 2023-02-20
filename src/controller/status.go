package controller

import (
	"github.com/gin-gonic/gin"
)

func (controller *Controller) Status(c *gin.Context) {
	controller.FeedbackOK(c, nil)
}
