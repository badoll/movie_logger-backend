package user

import (
	"net/http"
	"strings"

	"github.com/badoll/movie_logger-backend/api"
	"github.com/badoll/movie_logger-backend/db"
	"github.com/badoll/movie_logger-backend/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type likeReq struct {
	UserID  int64 `json:"user_id"`
	MovieID int64 `json:"movie_id"`
	Like    bool  `json:"like"`
}

func LikeMovie(c *gin.Context) {
	req := likeReq{}
	if err := c.BindJSON(&req); err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": c.Request.URL, "error": err}).Error()
		c.PureJSON(http.StatusOK, api.NewResp(api.ParamErr, "err", api.NilStruct))
		return
	}
	if err := db.GetCli().LikeMovie(req.UserID, req.MovieID, req.Like); err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": c.Request.URL, "error": err}).Error()
		c.PureJSON(http.StatusOK, api.NewResp(api.DBErr, "err", api.NilStruct))
		return
	}
	// 用户点赞电影后异步更新用户电影推荐因素(导演/编剧/演员/电影类型)
	AddCalUser(req.UserID)
	c.PureJSON(http.StatusOK, api.NewResp(api.Succ, "succ", api.NilStruct))
}

type interFieldReq struct {
	UserID     int64    `json:"user_id"`
	InterField []string `json:"inter_field"`
}

// SetUserInterField 设置用户感兴趣的电影类
func SetUserInterField(c *gin.Context) {
	req := interFieldReq{}
	if err := c.BindJSON(&req); err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": c.Request.URL, "error": err}).Error()
		c.PureJSON(http.StatusOK, api.NewResp(api.ParamErr, "err", api.NilStruct))
		return
	}
	interField := strings.Join(req.InterField, ",")
	if err := db.GetCli().SetUserInterField(req.UserID, interField); err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": c.Request.URL, "error": err}).Error()
		c.PureJSON(http.StatusOK, api.NewResp(api.DBErr, "err", api.NilStruct))
		return
	}
	c.PureJSON(http.StatusOK, api.NewResp(api.Succ, "succ", api.NilStruct))
}
