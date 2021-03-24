package movie

import (
	"net/http"

	"github.com/badoll/movie_logger-backend/api"
	"github.com/badoll/movie_logger-backend/db"
	"github.com/badoll/movie_logger-backend/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetRecommendByUser 根据用户open_id拿到推荐的电影list
func GetRecommendByUser(c *gin.Context) {
	userID := c.Param("user_id")
	movieList, err := db.GetCli().GetRecommendByUser(userID)
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": c.Request.URL, "error": err}).Error()
		c.PureJSON(http.StatusOK, api.NewResp(api.DBErr, "err", api.NilStruct))
		return
	}
	respList := make([]Movie, len(movieList))
	for i, v := range movieList {
		respList[i] = transMovie(v)
	}
	resp := movieListResp{
		Total:     len(movieList),
		MovieList: respList,
	}
	c.PureJSON(http.StatusOK, api.NewResp(api.Succ, "succ", resp))
}

// GetRecommendByMovie 电影相关推荐
func GetRecommendByMovie(c *gin.Context) {
	doubanID := c.Param("movie_id")
	movieList, err := db.GetCli().GetRecommendByMovie(doubanID)
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": c.Request.URL, "error": err}).Error()
		c.PureJSON(http.StatusOK, api.NewResp(api.DBErr, "err", api.NilStruct))
		return
	}
	respList := make([]Movie, len(movieList))
	for i, v := range movieList {
		respList[i] = transMovie(v)
	}
	resp := movieListResp{
		Total:     len(movieList),
		MovieList: respList,
	}
	c.PureJSON(http.StatusOK, api.NewResp(api.Succ, "succ", resp))
}
