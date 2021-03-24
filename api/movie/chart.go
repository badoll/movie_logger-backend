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

func GetMovieListByChart(c *gin.Context) {
	chart := c.Param("chart")
	limit := c.Query("limit")
	offset := c.Query("offset")
	var movieList []db.Movie
	var err error
	switch chart {
	case "nowplaying":
		fallthrough
	case "upcoming":
		/*
			目前只有nowplaying和upcoming有分页的需求
		*/
		if checkIntParam(limit) && checkIntParam(offset) {
			l, _ := strconv.Atoi(limit)
			o, _ := strconv.Atoi(offset)
			movieList, err = db.GetCli().GetMovieListByChartWithPage(chart, l, o)
		} else {
			movieList, err = db.GetCli().GetMovieListByChart(chart)
		}
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
		respList[i] = transMovie(v)
	}
	resp := movieListResp{
		Total:     len(movieList),
		MovieList: respList,
	}
	c.PureJSON(http.StatusOK, api.NewResp(api.Succ, "succ", resp))
}

func checkIntParam(param string) bool {
	if len(param) == 0 {
		return false
	}
	if _, err := strconv.Atoi(param); err != nil {
		return false
	}
	return true
}
