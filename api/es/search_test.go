package es

import (
	"encoding/json"
	"testing"
)

func TestMarshal(t *testing.T) {
	from := 0
	size := 10
	q := searchBody{
		Sort: []map[string]interface{}{
			{"_score": "desc"}, {"rating_score": "desc"},
		},
		From: &from,
		Size: &size,
		Query: boolQuery{
			Bool: boolCond{
				Must: []mustCond{
					{Match: map[string]interface{}{"title": "杀手"}},
				},
				Filter: []filterCond{
					{Term: map[string]interface{}{"cate": "剧情"}},
					{Term: map[string]interface{}{"cate": "动画"}},
				},
			},
		},
	}
	b, _ := json.Marshal(q)
	t.Log(string(b))
}
