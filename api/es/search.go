package es

import (
	"encoding/json"
	"strconv"

	"github.com/badoll/movie_logger-backend/api"
	"github.com/badoll/movie_logger-backend/logger"
	"github.com/sirupsen/logrus"
)

const (
	// 默认传入的should条件要匹配40%, 四舍五入
	MINSHOULDMATCH = "40%"
)

type searchBody struct {
	Sort  []map[string]interface{} `json:"sort,omitempty"`
	From  *int                     `json:"from,omitempty"`
	Size  *int                     `json:"size,omitempty"`
	Query interface{}              `json:"query"`
}

// GET movie/_search?pretty
// {"from":0,"size":3,"query":{"match":{"title":"杀手"}}}
// type fullTextQuery struct {
// 	Match map[string]interface{} `json:"match"`
// }

// GET movie/_search?pretty
// {"sort":["_score",{"rating_score":"desc"}],"from":0,"size":3,"query":{"bool":{"must":[{"match":{"title":"杀手"}}],"filter":[{"term":{"cate":"剧情"}},{"term":{"cate":"动作"}}]}}}
type boolQuery struct {
	Bool boolCond `json:"bool"`
}

type boolCond struct {
	Must           []mustCond   `json:"must,omitempty"`
	Filter         []filterCond `json:"filter,omitempty"`
	Should         []filterCond `json:"should,omitempty"`
	MinShouldMatch string       `json:"minimum_should_match,omitempty"` //should条件中匹配的最小百分比 eg: 40%
}

type mustCond struct {
	Match map[string]interface{} `json:"match"`
}

type filterCond struct {
	Term map[string]interface{} `json:"term"`
}

type searchResp struct {
	Took int  `json:"took"`
	Hits hits `json:"hits"`
}

type hits struct {
	Total total     `json:"total"`
	Hits  []hitItem `json:"hits"`
}

type total struct {
	Value int `json:"value"`
}

type hitItem struct {
	ID string `json:"_id"`
}

// title(required): 搜索的电影标题
// filter(optional): { "cate": ["剧情","动作"], "performer": [“”,“”,“”]}
// limit/offset (required)
// Search 用于搜索页，搜索标题并过滤电影类型/演员/导演/编剧
func Search(title string, filter map[string][]interface{}, limit, offset int) (idList []int64, err error) {
	// 没有传标题则只过滤
	if len(title) == 0 {
		return Filter(filter, limit, offset)
	}
	reqBody := searchBody{
		// 排序优先级：标题相似度得分 > 电影评分
		Sort: []map[string]interface{}{
			{"_score": "desc"}, {"rating_score": "desc"},
		},
		From: &offset,
		Size: &limit,
		Query: boolQuery{
			Bool: boolCond{
				Must: []mustCond{
					{Match: map[string]interface{}{"title": title}},
				},
				Filter: getFilterConds(filter),
			},
		},
	}
	return getSearchResp(reqBody)
}

// filter(required): { "cate": ["剧情","动作"], "performer": [“”,“”,“”]}
// limit/offset (required)
// Filter 用于过滤电影类型/演员/导演/编剧
func Filter(filter map[string][]interface{}, limit, offset int) (idList []int64, err error) {
	reqBody := searchBody{
		// filter context 没有_score, 用电影评分排序
		Sort: []map[string]interface{}{
			{"rating_score": "desc"},
		},
		From: &offset,
		Size: &limit,
		Query: boolQuery{
			Bool: boolCond{
				Filter: getFilterConds(filter),
			},
		},
	}
	return getSearchResp(reqBody)
}

// RelativeFilter 使用should语句检索相关电影
func RelativeFilter(should map[string][]interface{}, limit, offset int) (idList []int64, err error) {
	reqBody := searchBody{
		// filter context 没有_score, 用电影评分排序
		Sort: []map[string]interface{}{
			{"rating_score": "desc"},
		},
		From: &offset,
		Size: &limit,
		Query: boolQuery{
			Bool: boolCond{
				Should:         getFilterConds(should),
				MinShouldMatch: MINSHOULDMATCH, //默认传入的should条件要匹配40%, 四舍五入
			},
		},
	}
	return getSearchResp(reqBody)
}

// filter: { "cate": ["剧情","动作"] }
// => {"filter":[{"term":{"cate":"剧情"}},{"term":{"cate":"动画"}}]}
func getFilterConds(filter map[string][]interface{}) []filterCond {
	// filter(omitempty) 为空的情况下不会传该值
	filterConds := []filterCond{}
	for k, v := range filter {
		// 同一个tag多个过滤值需要拆分开
		for _, vv := range v {
			filterConds = append(filterConds, filterCond{
				Term: map[string]interface{}{
					k: vv,
				},
			})
		}
	}
	return filterConds
}

func getSearchResp(reqBody searchBody) (idList []int64, err error) {
	jsonB, err := json.Marshal(reqBody)
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": "_search", "error": err, "req": reqBody}).Error()
		return
	}
	url := ESHost + "/movie/_search"
	respBody, err := api.Request("GET", url, string(jsonB))
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": "_search", "error": err, "req": string(jsonB)}).Error()
		return
	}
	resp := searchResp{}
	if err = json.Unmarshal(respBody, &resp); err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": "_search", "error": err, "resp": string(respBody)}).Error()
		return
	}
	for _, hitItem := range resp.Hits.Hits {
		id, err := strconv.ParseInt(hitItem.ID, 10, 64)
		if err != nil {
			return idList, err
		}
		idList = append(idList, id)
	}
	logger.GetLogger("eslog").WithFields(logrus.Fields{"api": "_search", "req": string(jsonB), "resp": string(respBody)}).Debug()
	return
}
