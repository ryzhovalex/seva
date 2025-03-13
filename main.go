package main

import (
	"flag"
	"seva/lib/bone"
	"seva/lib/shell"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Event struct {
	Domain string
	// Particular number of event in domain.
	Order int
	// Time of event injection.
	Created_Sec int
	// Integer type of an event. Each project has own unsigned set of types,
	// starting from 1.
	Type int
}

func create_server() *gin.Engine {
	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(cors.Default())

	return server
}

func main() {
	shell_enabled := flag.Bool("shell", false, "Enables shell mode.")
	bone.Init("seva")
	shell.Init("seva")

	if *shell_enabled {
		shell.Run()
		return
	}

	server := create_server()
	server.Run("0.0.0.0:3000")
}
