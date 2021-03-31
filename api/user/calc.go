package user

import (
	"strings"
	"sync"
	"time"

	"github.com/badoll/movie_logger-backend/api/movie"
	"github.com/badoll/movie_logger-backend/db"
	"github.com/badoll/movie_logger-backend/logger"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

const (
	// 用户喜欢的电影列表中某一个tag出现的次数占比超过该比率则加入用户感兴趣的tag
	RATE2RECOMMEND float32 = 0.35
)

var c = &calBot{userList: map[int64]struct{}{}}

// CalInit 后台每隔一个小时计算更新一次用户推荐因素
func CalInit() {
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for {
			<-ticker.C
			c.Cal()
		}
	}()
}

type calBot struct {
	userList map[int64]struct{}
	mu       sync.RWMutex
}

func AddCalUser(userID int64) {
	c.mu.Lock()
	c.userList[userID] = struct{}{}
	c.mu.Unlock()
}

func (c *calBot) Cal() {
	c.mu.RLock()
	userList := make([]int64, len(c.userList))
	i := 0
	for k := range c.userList {
		userList[i] = k
		i++
	}
	c.mu.RUnlock()
	logger.GetLogger("cal").WithFields(logrus.Fields{"api": "Cal"}).Debug("start to cal")
	CalUserInter(userList)
	logger.GetLogger("cal").WithFields(logrus.Fields{"api": "Cal"}).Debug("done")
}

// CalUserInter 通过用户喜欢的电影列表提取电影推荐因素中占比最大的几个值
func CalUserInter(userList []int64) {
	if len(userList) == 0 {
		return
	}
	for _, userID := range userList {
		likeList, err := db.GetCli().GetUserLikeList(userID)
		if err != nil {
			logger.GetLogger("cal").WithFields(logrus.Fields{"api": "CalUserInter", "error": err}).Error()
			return
		}
		if len(likeList) == 0 {
			return
		}
		daoList, err := db.GetCli().SelectMovieDetailByMovieIDList(likeList)
		if err != nil {
			logger.GetLogger("cal").WithFields(logrus.Fields{"api": "CalUserInter", "error": err}).Error()
			return
		}
		mlist := make([]movie.Movie, len(daoList))
		for i, v := range daoList {
			mlist[i] = movie.TransMovie(v)
		}
		// rfMap 用户喜欢的所有电影的推荐因素集合
		rfMap := map[string]map[string]int{
			"inter_field":     {},
			"inter_director":  {},
			"inter_writer":    {},
			"inter_performer": {},
		}
		for _, m := range mlist {
			for _, v := range m.Cate {
				calFactor("inter_field", v, rfMap)
			}
			d := m.Director
			if len(d) > 1 {
				d = d[:1]
			}
			for _, v := range d {
				calFactor("inter_director", v, rfMap)
			}
			w := m.Writer
			if len(w) > 1 {
				w = w[:1]
			}
			for _, v := range w {
				calFactor("inter_writer", v, rfMap)
			}
			p := m.Performer
			if len(p) > 3 {
				p = p[:3] //防止推荐重复的脏数据
			}
			for _, v := range p {
				calFactor("inter_performer", v, rfMap)
			}
		}
		userInter, err := db.GetCli().GetUesrInter(userID)
		if err != nil {
			logger.GetLogger("cal").WithFields(logrus.Fields{"api": "CalUserInter", "error": err}).Error()
			return
		}
		// 更新用户兴趣
		rf := movie.RecommendFactor{
			Cate:      mergeInter(db.SplitString(userInter.InterField), rfMap["inter_field"]),
			Director:  mergeInter(db.SplitString(userInter.InterDirector), rfMap["inter_director"]),
			Writer:    mergeInter(db.SplitString(userInter.InterWriter), rfMap["inter_writer"]),
			Performer: mergeInter(db.SplitString(userInter.InterPerformer), rfMap["inter_performer"]),
		}
		daoRF := db.UserInter{
			InterField:     strings.Join(rf.Cate, ","),
			InterDirector:  strings.Join(rf.Director, ","),
			InterWriter:    strings.Join(rf.Writer, ","),
			InterPerformer: strings.Join(rf.Performer, ","),
		}
		if err = db.GetCli().UpdateUserInter(userID, daoRF); err != nil {
			logger.GetLogger("cal").WithFields(logrus.Fields{"api": "CalUserInter", "error": err}).Error()
			return
		}
	}
}

// 判断某一个tag出现的次数比重是否超过 RATE2RECOMMEND, 超过则加入userInter
func mergeInter(userInter []string, valMap map[string]int) []string {
	userInterMap := map[string]struct{}{}
	for _, v := range userInter {
		userInterMap[v] = struct{}{}
	}
	total := 0
	for _, v := range valMap {
		total += v
	}
	newInter := []string{}
	for k, v := range valMap {
		rate, _ := decimal.NewFromInt(int64(v)).Div(decimal.NewFromInt(int64(total))).Round(1).Float64()
		if rate > float64(RATE2RECOMMEND) {
			newInter = append(newInter, k)
		}
	}
	for _, v := range newInter {
		if _, ok := userInterMap[v]; !ok {
			userInter = append(userInter, v)
		}
	}
	return userInter
}

func calFactor(factor, val string, rfMap map[string]map[string]int) {
	// When accessing a map key that doesn't exist yet, one receives a zero value.
	// Often, the zero value is a suitable value, for example when using append or doing integer math.
	rfMap[factor][val] += 1
}
