package main

import (
	_ "github.com/badoll/movie_logger-backend/config"
	"github.com/badoll/movie_logger-backend/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	initRouter(router)
	initServer()
	logger.GetDefaultLogger().Debug("server start!")
	router.Run() // default listen and serve on 8080
}
