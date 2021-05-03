package movie

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/badoll/movie_logger-backend/api"
	"github.com/badoll/movie_logger-backend/db"
	"github.com/badoll/movie_logger-backend/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetMovieListByChart(c *gin.Context) {
	chart := c.Param("chart")
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	var movieList []db.Movie
	var err error
	switch chart {
	case "nowplaying":
		fallthrough
	case "upcoming":
		//目前只有nowplaying和upcoming需要分页
		movieList, err = db.GetCli().GetMovieListByChartWithPage(chart, limit, offset)
	case "newrelease":
		fallthrough
	case "weeklytop":
		movieList, err = db.GetCli().GetMovieListByChart(chart)
	default:
		c.PureJSON(http.StatusOK, api.NewResp(api.ParamErr, "err", api.NilStruct))
		return
	}
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": c.Request.URL, "error": err}).Error()
		c.PureJSON(http.StatusOK, api.NewResp(api.DBErr, "err", api.NilStruct))
		return
	}
	respList := make([]Movie, len(movieList))
	for i, v := range movieList {
		respList[i] = TransMovie(v)
	}
	resp := movieListResp{
		Total:     len(movieList),
		MovieList: respList,
	}
	c.PureJSON(http.StatusOK, api.NewResp(api.Succ, "succ", resp))
}
