package rpc

import (
	"encoding/json"
	"io"
	"seva/lib/utils"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, e *utils.Error) {
	c.JSON(400, gin.H{"Code": e.Code(), "Body": e.Message()})
}

func Ok(c *gin.Context, code utils.Code, body any) {
	if code > 0 {
		panic("Ok code must be 0 or negative.")
	}
	c.JSON(200, gin.H{"Code": code, "Body": body})
}

func JsonRequestBody(c *gin.Context, v *any) *utils.Error {
	body, be := io.ReadAll(c.Request.Body)
	if be != nil {
		return utils.CreateDefaultErrorFromBase(be)
	}

	be = json.Unmarshal(body, v)
	if be != nil {
		return utils.CreateDefaultErrorFromBase(be)
	}
	return nil
}

func TextRequestBody(c *gin.Context) (string, *utils.Error) {
	body, be := io.ReadAll(c.Request.Body)
	if be != nil {
		return "", utils.CreateDefaultErrorFromBase(be)
	}
	return string(body), nil
}
