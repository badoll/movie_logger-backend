package main

import (
	"github.com/badoll/movie_logger-backend/api/movie"
	"github.com/badoll/movie_logger-backend/api/user"
	"github.com/badoll/movie_logger-backend/db"
	"github.com/gin-gonic/gin"
)

// initRouter 注册api
func initRouter(router *gin.Engine) {
	v := router.Group("/movie_logger")
	// 电影相关接口
	v.GET("/movie/id/:douban_id", movie.GetMovieDetail)
	v.GET("/movie/chart/:chart", movie.GetMovieListByChart)
	v.GET("/movie/recommend/user/:user_id", movie.GetRecommendByUser)
	v.GET("/movie/recommend/movie/:movie_id", movie.GetRecommendByMovie)
	v.POST("/movie/search", movie.SearchMovie)

	// 用户相关接口
	v.GET("/user/wx_login/:code", user.GetWXOpenID)
}

func initServer() {
	db.Init()
}
