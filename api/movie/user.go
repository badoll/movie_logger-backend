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

// Like 用户喜欢的电影
func Like(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	movieIDList, err := db.GetCli().GetUserLikeListWithPage(userID, limit, offset)
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": c.Request.URL, "error": err}).Error()
		c.PureJSON(http.StatusOK, api.NewResp(api.DBErr, "err", api.NilStruct))
		return
	}
	daoList, err := db.GetCli().SelectMovieDetailByMovieIDList(movieIDList)
	if err != nil {
		return
	}
	mlist := make([]Movie, len(daoList))
	for i, v := range daoList {
		mlist[i] = TransMovie(v)
	}
	resp := movieListResp{
		Total:     len(mlist),
		MovieList: mlist,
	}
	c.PureJSON(http.StatusOK, api.NewResp(api.Succ, "succ", resp))
}
