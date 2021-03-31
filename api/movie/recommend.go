package movie

import (
	"net/http"
	"strconv"

	"github.com/badoll/movie_logger-backend/api"
	"github.com/badoll/movie_logger-backend/api/es"
	"github.com/badoll/movie_logger-backend/db"
	"github.com/badoll/movie_logger-backend/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetRecommendByUser 根据用户open_id拿到推荐的电影list
func GetRecommendByUser(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	userID, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)
	userInter, err := db.GetCli().GetUesrInter(userID)
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": c.Request.URL, "error": err}).Error()
		c.PureJSON(http.StatusOK, api.NewResp(api.DBErr, "err", api.NilStruct))
		return
	}

	rf := RecommendFactor{
		Cate:      db.SplitString(userInter.InterField),
		Director:  db.SplitString(userInter.InterDirector),
		Writer:    db.SplitString(userInter.InterWriter),
		Performer: db.SplitString(userInter.InterPerformer),
	}
	movieList, err := getRecommend(rf, limit, offset)
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": c.Request.URL, "error": err}).Error()
		c.PureJSON(http.StatusOK, api.NewResp(api.SvrErr, "err", api.NilStruct))
		return
	}
	resp := movieListResp{
		Total:     len(movieList),
		MovieList: movieList,
	}
	c.PureJSON(http.StatusOK, api.NewResp(api.Succ, "succ", resp))
}

// GetRecommendByMovie 电影相关推荐
func GetRecommendByMovie(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	movieID, _ := strconv.ParseInt(c.Param("movie_id"), 10, 64)
	movie, err := db.GetCli().GetMovieDetailByMovieID(movieID)
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": c.Request.URL, "error": err}).Error()
		c.PureJSON(http.StatusOK, api.NewResp(api.DBErr, "err", api.NilStruct))
		return
	}
	rf := RecommendFactor{
		Cate:      db.SplitString(movie.Cate),
		Director:  db.SplitString(movie.Director),
		Writer:    db.SplitString(movie.Writer),
		Performer: db.SplitString(movie.Performer),
	}

	// 防止推荐重复的脏数据
	if len(rf.Director) > 1 {
		rf.Director = rf.Director[:1]
	}
	if len(rf.Writer) > 1 {
		rf.Writer = rf.Writer[:1]
	}
	if len(rf.Performer) > 3 {
		rf.Performer = rf.Performer[:3]
	}
	movieList, err := getRecommend(rf, limit, offset)
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": c.Request.URL, "error": err}).Error()
		c.PureJSON(http.StatusOK, api.NewResp(api.SvrErr, "err", api.NilStruct))
		return
	}
	resp := movieListResp{
		Total:     len(movieList),
		MovieList: movieList,
	}
	c.PureJSON(http.StatusOK, api.NewResp(api.Succ, "succ", resp))
}

// RecommendFactor 电影推荐因素
type RecommendFactor struct {
	Cate      []string
	Director  []string
	Writer    []string
	Performer []string
}

// getRecommend 通过es索引拉电影id，通过电影主键id拉数据
func getRecommend(rf RecommendFactor, limit, offset int) (mlist []Movie, err error) {
	shouldConds := map[string][]interface{}{
		"cate":      ToGenericArray(rf.Cate),
		"director":  ToGenericArray(rf.Director),
		"writer":    ToGenericArray(rf.Writer),
		"performer": ToGenericArray(rf.Performer),
	}
	// 通过es索引拉取相关电影id
	movieIDList, err := es.RelativeFilter(shouldConds, limit, offset)
	if err != nil {
		return
	}
	movieList, err := db.GetCli().SelectMovieDetailByMovieIDList(movieIDList)
	if err != nil {
		return
	}
	mlist = make([]Movie, len(movieList))
	for i, v := range movieList {
		mlist[i] = TransMovie(v)
	}
	return
}

// ToGenericArray []string -> []interface{}
func ToGenericArray(arr []string) []interface{} {
	arri := make([]interface{}, len(arr))
	for i, v := range arr {
		arri[i] = v
	}
	return arri
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
		respList[i] = TransMovie(v)
	}
	resp := movieListResp{
		Total:     len(movieList),
		MovieList: respList,
	}
	c.PureJSON(http.StatusOK, api.NewResp(api.Succ, "succ", resp))
}
