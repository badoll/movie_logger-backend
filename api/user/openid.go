package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/badoll/movie_logger-backend/api"
	"github.com/badoll/movie_logger-backend/config"
	"github.com/badoll/movie_logger-backend/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

/*
	req:
	/user/code/xxx
*/

type wxResp struct {
	OpenID string `json:"openid"`
	// UnionID string `json:"unionid"`
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// 接口：https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html
func GetWXOpenID(c *gin.Context) {
	code := c.Param("code")
	conf := config.GetConfig().WXConf
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		conf.AppID, conf.Secret, code)
	httpResp, err := http.Get(url)
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": url, "error": err}).Error("http request error")
		c.PureJSON(http.StatusOK, api.NewResp(api.HTTPErr, "err", api.NilStruct))
		return
	}
	defer httpResp.Body.Close()
	body, _ := ioutil.ReadAll(httpResp.Body)
	resp := wxResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": url, "error": err}).Error("http request error")
		c.PureJSON(http.StatusOK, api.NewResp(api.HTTPErr, "err", api.NilStruct))
		return
	}
	if resp.ErrCode != 0 {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": url, "error": err, "resp": resp}).Error("code to open_id request error")
		c.PureJSON(http.StatusOK, api.NewResp(api.HTTPErr, "err", resp))
		return
	}
	c.PureJSON(http.StatusOK, api.NewResp(api.Succ, "succ", resp))
}
