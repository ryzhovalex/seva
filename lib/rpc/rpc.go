package rpc

import (
	"seva/lib/utils"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, e *utils.Error) {
	c.JSON(400, gin.H{"Code": e.Code(), "Body": e.Message()})
}

func Ok(c *gin.Context, code utils.Code, body any) {
	if code >= 0 {
		panic("Ok code must be negative.")
	}
	c.JSON(200, gin.H{"Code": code, "Body": body})
}
