package es

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/badoll/movie_logger-backend/api"
	"github.com/badoll/movie_logger-backend/db"
	"github.com/badoll/movie_logger-backend/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var bState bulkState

func init() {
	bState = bulkState{
		State: State{Finished: true},
		done:  make(chan struct{}, 1),
	}
}

// bulkState bulk状态
type bulkState struct {
	State
	done chan struct{}
}

func (b *bulkState) TryStart() bool {
	select {
	case b.done <- struct{}{}:
		b.Finished = false
		b.Took = 0
		b.Offset = 0
		b.ErrMsg = ""
		return true
	default:
	}
	return false
}

func (b *bulkState) Done() {
	b.Finished = true
	<-b.done
}

type State struct {
	Finished bool   `json:"finished,omitempty"`
	Offset   int    `json:"offset,omitempty"`
	Took     int    `json:"took,omitempty"`
	ErrMsg   string `json:"err_msg,omitempty"`
}

// GetState 拿到当前bulk状态
func GetState(c *gin.Context) {
	c.PureJSON(http.StatusOK, api.NewResp(api.Succ, "succ", bState))
}

// Reload es重新读db数据入索引
func Reload(c *gin.Context) {
	if !ReloadIndex() {
		c.PureJSON(http.StatusOK, api.NewResp(api.SvrErr, "err", map[string]interface{}{"msg": "reload not done"}))
		return
	}
	c.PureJSON(http.StatusOK, api.NewResp(api.Succ, "succ", map[string]interface{}{"msg": "reload start"}))
}

func ReloadIndex() bool {
	if !bState.TryStart() {
		return false
	}
	go reload()
	return true
}

func reload() {
	var err error
	defer func() {
		if err != nil {
			logger.GetDefaultLogger().WithFields(logrus.Fields{"api": "reload", "error": err}).Error()
			bState.ErrMsg = err.Error()
		}
		bState.Done()
	}()
	count, err := db.GetCli().GetMovieCount()
	if err != nil {
		return
	}
	limit := 100 //服务器内存太小了，只能一次拉取一点
	offset := 0
	took := 0
	var mlist []db.MovieIndex
	indexList := make([]MovieIndex, limit)
	// MovieIndexData 用自增id排序，分批拉顺序不会有问题
	for ; offset < count; offset += limit {
		if offset+limit > count {
			// 最后一次循环
			limit = count - offset
		}
		mlist, err = db.GetCli().GetMovieIndexData(limit, offset)
		if err != nil {
			return
		}
		for i, v := range mlist {
			indexList[i] = transMovie(v)
		}
		took, err = bulk(indexList[:len(mlist)])
		if err != nil {
			return
		}
		bState.Took += took
		bState.Offset = offset
	}
}

type bulkResp struct {
	Took   int  `json:"took"`
	Errors bool `json:"errors"`
}

func bulk(mlist []MovieIndex) (took int, err error) {
	buffer := indexBuffer(mlist)
	url := ESHost + "/_bulk"
	respBody, err := api.Request("POST", url, buffer.String())
	if err != nil {
		return
	}
	resp := bulkResp{}
	if err = json.Unmarshal(respBody, &resp); err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": url, "error": err, "resp": string(respBody)}).Error()
		return
	}
	if resp.Errors {
		err = fmt.Errorf("es bulk error")
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": url, "error": err, "resp": resp}).Error()
		return
	}
	took = resp.Took
	return
}

func indexBuffer(mlist []MovieIndex) bytes.Buffer {
	buf := bytes.Buffer{}
	for _, v := range mlist {
		h := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": ESIndex,
				"_id":    v.MovieID,
			},
		}
		jsonB, _ := json.Marshal(h)
		buf.Write(jsonB)
		buf.WriteByte('\n')
		jsonB, _ = json.Marshal(v)
		buf.Write(jsonB)
		buf.WriteByte('\n')
	}
	return buf
}

// MovieIndex 用于在es作index的数据
type MovieIndex struct {
	MovieID     int64    `json:"-"`
	Title       string   `json:"title"`
	Cate        []string `json:"cate"`
	Director    []string `json:"director"`
	Writer      []string `json:"writer"`
	Performer   []string `json:"performer"`
	Region      []string `json:"region"`
	Language    []string `json:"language"`
	ReleaseYear int      `json:"release_year,omitempty"` // 未知年份传空值
	RatingScore int      `json:"rating_score"`
}

// transMovie 转换dao数据类型
func transMovie(mdao db.MovieIndex) MovieIndex {
	releaseYear, err := strconv.Atoi(mdao.ReleaseYear)
	if err != nil {
		releaseYear = 0
	}
	// 保留1位小数点精度
	// ratingScore, _ := decimal.NewFromInt(int64(mdao.RatingScore)).Div(decimal.NewFromInt(100)).Round(1).Float64()
	return MovieIndex{
		MovieID:     mdao.MovieID,
		Title:       mdao.Title,
		Cate:        db.SplitString(mdao.Cate),
		Director:    db.SplitString(mdao.Director),
		Writer:      db.SplitString(mdao.Writer),
		Performer:   db.SplitString(mdao.Performer),
		Region:      db.SplitString(mdao.Region),
		Language:    db.SplitString(mdao.Language),
		ReleaseYear: releaseYear,
		RatingScore: int(mdao.RatingScore),
	}
}
