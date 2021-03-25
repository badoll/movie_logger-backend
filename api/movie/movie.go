package movie

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/badoll/movie_logger-backend/api"
	"github.com/badoll/movie_logger-backend/db"
	"github.com/badoll/movie_logger-backend/logger"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

/*
	req: /movie?id=xxx
*/

/*
	resp:
	{
		base: {
			status: 0,
			message: "succ",
		},
		result: {

		}
	}
*/

type cast struct {
	Name   string `json:"name"`
	Role   string `json:"role"`
	Avatar string `json:"avatar"`
}

type Movie struct {
	MovieID     int64    `json:"id"`
	Title       string   `json:"title"`
	Poster      string   `json:"poster"`
	Cate        []string `json:"cate"`
	Director    []string `json:"director"`
	Writer      []string `json:"writer"`
	Performer   []string `json:"performer"`
	Region      []string `json:"region"`
	Language    []string `json:"language"`
	ReleaseYear string   `json:"release_year"`
	ReleaseDate []string `json:"release_date"`
	Runtime     []string `json:"runtime"`
	RatingScore float64  `json:"rating_score"`
	RatingNum   uint     `json:"rating_num"`
	RatingStars []string `json:"rating_stars"`
	Intro       string   `json:"intro"`
	MainCast    []cast   `json:"main_cast"`
	Photos      []string `json:"photos"`
}

func GetMovieDetail(c *gin.Context) {
	movieID, _ := strconv.ParseInt(c.Param("movie_id"), 10, 64)
	movie, err := db.GetCli().GetMovieDetailByMovieID(movieID)
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": c.Request.URL, "error": err}).Error()
		c.PureJSON(http.StatusOK, api.NewResp(api.DBErr, "err", api.NilStruct))
		return
	}
	c.PureJSON(http.StatusOK, api.NewResp(api.Succ, "succ", transMovie(movie)))
}

// transMovie 转换dao数据类型
func transMovie(mdao db.Movie) Movie {
	mainCast := make([]cast, 0)
	if err := json.Unmarshal([]byte(mdao.MainCast), &mainCast); err != nil {
		logger.GetDefaultLogger().WithField("error", err).Error("ummarshal main_cast error")
	}
	// 保留1位小数点精度
	ratingScore, _ := decimal.NewFromInt(int64(mdao.RatingScore)).Div(decimal.NewFromInt(100)).Round(1).Float64()
	return Movie{
		MovieID:     mdao.MovieID,
		Title:       mdao.Title,
		Poster:      mdao.Poster,
		Cate:        db.SplitString(mdao.Cate),
		Director:    db.SplitString(mdao.Director),
		Writer:      db.SplitString(mdao.Writer),
		Performer:   db.SplitString(mdao.Performer),
		Region:      db.SplitString(mdao.Region),
		Language:    db.SplitString(mdao.Language),
		ReleaseYear: mdao.ReleaseYear,
		ReleaseDate: db.SplitString(mdao.ReleaseDate),
		Runtime:     db.SplitString(mdao.Runtime),
		RatingScore: ratingScore,
		RatingNum:   mdao.RatingNum,
		RatingStars: db.SplitString(mdao.RatingStars),
		Intro:       mdao.Intro,
		MainCast:    mainCast,
		Photos:      db.SplitString(mdao.Photos),
	}
}
