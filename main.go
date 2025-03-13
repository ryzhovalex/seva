package main

import (
	"seva/lib/bone"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func create_server() *gin.Engine {
	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(cors.Default())

	return server
}

func convert_sec_to_str(sec int) string {
	return bone.Date_Sec(sec, "2006-01-02 15:04")
}

func main() {
	bone.Init("seva")
	server := create_server()
	server.Run("0.0.0.0:3000")
}
