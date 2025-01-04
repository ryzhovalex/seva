package main

import (
	"context"
	"seva/components"

	"github.com/a-h/templ"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func viewRenderer(
	view templ.Component,
	status int,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		component := components.Index(components.IndexArgs{View: view})
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
	c.File("static/" + c.Param("name"))
}

func renderNotFound(c *gin.Context) {
	c.Status(404)
	components.NotFound().Render(context.Background(), c.Writer)
}

func createServer() *gin.Engine {
	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(cors.Default())

	server.GET("/", viewRenderer(components.NotFound(), 200))
	server.GET("/favicon.ico", getFavicon)
	server.GET("/Static", getStatic)
	return server
}

func main() {
	server := createServer()
	server.Run("0.0.0.0:3000")
}
