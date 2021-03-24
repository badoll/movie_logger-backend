package movie

import (
	"net/http"

	"github.com/badoll/movie_logger-backend/api"
	"github.com/badoll/movie_logger-backend/db"
	"github.com/badoll/movie_logger-backend/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

/*
	req:
	{
		"title":"xxx",
	}
*/

/*
	resp:
	{
		base: {
			status: 0,
			message: "succ",
		},
		result: {
			total: 0,
			movie_list: [],
		}
	}
*/
type searchReq struct {
	Title  string `json:"title"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}
type movieListResp struct {
	Total     int     `json:"total"`
	MovieList []Movie `json:"movie_list"`
}

func SearchMovie(c *gin.Context) {
	req := searchReq{}
	if err := c.BindJSON(&req); err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": c.Request.URL, "error": err}).Error()
		c.PureJSON(http.StatusOK, api.NewResp(api.ParamErr, "err", api.NilStruct))
		return
	}
	movieList, err := db.GetCli().GetMovieListByTitle(req.Title, req.Limit, req.Offset)
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": c.Request.URL, "error": err}).Error()
		c.PureJSON(http.StatusOK, api.NewResp(api.DBErr, "err", api.NilStruct))
		return
	}
	respList := make([]Movie, len(movieList))
	for i, v := range movieList {
		respList[i] = transMovie(v)
	}
	resp := movieListResp{Total: len(movieList), MovieList: respList}
	c.PureJSON(http.StatusOK, api.NewResp(api.Succ, "succ", resp))
}
