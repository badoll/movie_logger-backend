package main

import (
	"github.com/badoll/movie_logger-backend/api/movie"
	"github.com/badoll/movie_logger-backend/db"
	"github.com/gin-gonic/gin"
)

// initRouter 注册api
func initRouter(router *gin.Engine) {
	v := router.Group("/movie_logger")
	// 电影相关接口
	v.GET("/movie/:douban_id", movie.GetMovieDetail)
	v.POST("/movie/search", movie.SearchMovie)
}

func initServer() {
	db.Init()
}
