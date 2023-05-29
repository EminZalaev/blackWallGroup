package response

import (
	"github.com/gin-gonic/gin"
)

func ReturnError(c *gin.Context, httpStatus int, err error) {
	c.JSON(httpStatus, gin.H{
		"status": httpStatus,
		"detail": err.Error(),
	})
}
