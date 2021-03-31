package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/badoll/movie_logger-backend/api"
	"github.com/badoll/movie_logger-backend/config"
	"github.com/badoll/movie_logger-backend/db"
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

type loginResp struct {
	UserID    int64   `json:"user_id"`
	LikeList  []int64 `json:"like_list"`
	NickName  string  `json:"nick_name"`
	AvatarUrl string  `json:"avatar_url"`
}

// 接口：https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html
// Login 使用前端返回code返回用户id以及其他信息（默认注册逻辑）
func Login(c *gin.Context) {
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
	wxResp := wxResp{}
	err = json.Unmarshal(body, &wxResp)
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": url, "error": err}).Error("http request error")
		c.PureJSON(http.StatusOK, api.NewResp(api.HTTPErr, "err", api.NilStruct))
		return
	}
	if wxResp.ErrCode != 0 {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": url, "error": err, "resp": wxResp}).Error("code to open_id request error")
		c.PureJSON(http.StatusOK, api.NewResp(api.HTTPErr, "err", wxResp))
		return
	}
	userInfo, err := db.GetCli().GetUserInfo(wxResp.OpenID)
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": url, "error": err, "resp": wxResp}).Error("GetUserID error")
		c.PureJSON(http.StatusOK, api.NewResp(api.DBErr, "err", api.NilStruct))
		return
	}
	likeList, err := db.GetCli().GetUserLikeList(userInfo.ID)
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": url, "error": err, "resp": wxResp}).Error("GetUserID error")
		c.PureJSON(http.StatusOK, api.NewResp(api.DBErr, "err", api.NilStruct))
		return
	}
	if likeList == nil {
		// 没有则返回空数组
		likeList = make([]int64, 0)
	}
	resp := loginResp{
		UserID:    userInfo.ID,
		LikeList:  likeList,
		NickName:  userInfo.NickName,
		AvatarUrl: userInfo.AvatarUrl,
	}
	c.PureJSON(http.StatusOK, api.NewResp(api.Succ, "succ", resp))
}

type updateReq struct {
	NickName  string `json:"nick_name"`
	AvatarUrl string `json:"avatar_url"`
}

// UpdateUserInfo 更新用户数据
func UpdateUserInfo(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)
	req := updateReq{}
	if err := c.BindJSON(&req); err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": c.Request.URL, "error": err}).Error()
		c.PureJSON(http.StatusOK, api.NewResp(api.ParamErr, "err", api.NilStruct))
		return
	}
	userInfo := db.User{
		NickName:  req.NickName,
		AvatarUrl: req.AvatarUrl,
	}
	if err := db.GetCli().UpdateUserInfo(userID, userInfo); err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": c.Request.URL, "error": err}).Error()
		c.PureJSON(http.StatusOK, api.NewResp(api.DBErr, "err", api.NilStruct))
		return
	}
	c.PureJSON(http.StatusOK, api.NewResp(api.Succ, "succ", api.NilStruct))
}

// IsNewUser 判断是否是新用户（没有选择兴趣类型）
func IsNewUser(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)
	isNew, err := db.GetCli().IsNewUser(userID)
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": c.Request.URL, "error": err}).Error()
		c.PureJSON(http.StatusOK, api.NewResp(api.DBErr, "err", api.NilStruct))
		return
	}
	resp := struct {
		IsNew bool `json:"is_new"`
	}{IsNew: isNew}
	c.PureJSON(http.StatusOK, api.NewResp(api.Succ, "succ", resp))
}
