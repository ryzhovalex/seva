package main

import (
	"seva/internal/domains"
	"seva/internal/sevent"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func createServer() *gin.Engine {
	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(cors.Default())

	server.POST("/Rpc/Domains/CreateDomain", domains.RpcCreateDomain)
	server.POST("/Rpc/Domains/GetDomains", domains.RpcGetDomains)

	server.POST("/Rpc/Sevent/CreateEvent", sevent.RpcCreateEvent)
	server.POST("/Rpc/Sevent/GetSpecs", sevent.RpcGetSpecs)
	server.POST("/Rpc/Sevent/GetEvents", sevent.RpcGetEvents)

	return server
}

func main() {
	server := createServer()
	server.Run("0.0.0.0:3000")
}
