package movie

import (
	"net/http"
	"strconv"

	"github.com/badoll/movie_logger-backend/api"
	"github.com/badoll/movie_logger-backend/db"
	"github.com/badoll/movie_logger-backend/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetRecommendByUser 根据用户open_id拿到推荐的电影list
func GetRecommendByUser(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)
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
	movieID, _ := strconv.ParseInt(c.Param("movie_id"), 10, 64)
	movieList, err := db.GetCli().GetRecommendByMovie(movieID)
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

// GetRecommendDefault 默认推荐
func GetRecommendDefault(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	movieList, err := db.GetCli().GetTopMovieWithPage(limit, offset)
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
