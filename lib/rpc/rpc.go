package rpc

import (
	"seva/lib/utils"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, e *utils.Error) {
	c.JSON(400, gin.H{"Code": e.Code(), "Body": e.Message()})
}

func Ok(c *gin.Context, body any) {
	c.JSON(200, gin.H{"Code": 0, "Body": body})
}
