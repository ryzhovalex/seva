package rpc

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, e int) {
	c.Header("code", strconv.Itoa(e))
}

func Ok(c *gin.Context, body any) {
	c.Header("code", "0")
	c.JSON(200, body)
}
