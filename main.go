package main

import (
	"context"
	"encoding/json"
	"io"
	"seva/components"
	"seva/internal/domains"
	"seva/lib/rpc"
	"seva/lib/utils"
	"strings"

	"github.com/a-h/templ"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func renderView(
	view templ.Component,
	status int,
	title string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		component := components.Index(components.IndexArgs{View: view, Title: title})
		c.Status(status)
		component.Render(context.Background(), c.Writer)
	}
}

func getFavicon(c *gin.Context) {
	c.Status(200)
	c.File("Static/Favicon.ico")
}

func getStatic(c *gin.Context) {
	c.Status(200)
	c.File("static/" + c.Param("Name"))
}

func renderNotFound(c *gin.Context) {
	c.Status(404)
	components.NotFound().Render(context.Background(), c.Writer)
}

func jsonRequestBody(c *gin.Context, v *any) *utils.Error {
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

func textRequestBody(c *gin.Context) (string, *utils.Error) {
	body, be := io.ReadAll(c.Request.Body)
	if be != nil {
		return "", utils.CreateDefaultErrorFromBase(be)
	}
	return string(body), nil
}

func RpcCreateDomain(c *gin.Context) {
	domain, e := textRequestBody(c)
	if e != nil {
		rpc.Error(c, e)
		return
	}
	domain = strings.Replace(domain, "Domain=", "", 1)

	e = domains.CreateDomain(domain)
	if e != nil {
		rpc.Error(c, e)
		return
	}
	c.Header("HX-Redirect", "/")
	rpc.Ok(c, 0, nil)
}

func createServer() *gin.Engine {
	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(cors.Default())

	server.GET("/", renderView(components.Home(), 200, "Home"))
	server.GET("/CreateDomain", renderView(components.CreateDomain(), 200, "Create Domain"))
	server.GET("/favicon.ico", getFavicon)
	server.GET("/Static/:Name", getStatic)

	server.POST("/Rpc/Domains/Create", RpcCreateDomain)

	return server
}

func main() {
	server := createServer()
	server.Run("0.0.0.0:3000")
}
