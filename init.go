package main

import (
	"github.com/badoll/movie_logger-backend/api/es"
	"github.com/badoll/movie_logger-backend/api/movie"
	"github.com/badoll/movie_logger-backend/api/user"
	"github.com/badoll/movie_logger-backend/db"
	"github.com/gin-gonic/gin"
)

// initRouter 注册api
func initRouter(router *gin.Engine) {
	v := router.Group("/movie_logger")
	// 电影相关接口
	v.GET("/movie/id/:movie_id", movie.GetMovieDetail)
	v.GET("/movie/chart/:chart", movie.GetMovieListByChart)
	v.GET("/movie/user/like/:user_id", movie.Like)
	v.GET("/movie/recommend/user/:user_id", movie.GetRecommendByUser)
	v.GET("/movie/recommend/movie/:movie_id", movie.GetRecommendByMovie)
	v.GET("/movie/recommend/default", movie.GetRecommendDefault)
	v.POST("/movie/search", movie.SearchMovie)

	// 用户相关接口
	v.GET("/user/account/login/:code", user.Login)
	v.GET("/user/account/is_new/:user_id", user.IsNewUser)

	v.POST("/user/action/like", user.LikeMovie)
	v.POST("/user/action/set_inter_field", user.SetUserInterField)

	// es接口
	v.GET("/es/state", es.GetState)
	v.GET("/es/reload", es.Reload)

	v.GET("/test/:test_id", user.GetTest)
}

func initServer() {
	db.Init()
	user.CalInit()
}
