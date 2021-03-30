package api

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/badoll/movie_logger-backend/logger"
	"github.com/sirupsen/logrus"
)

// Request 封装http请求, 默认set Content-Type="application/json"
func Request(method, url string, body string) (respBody []byte, err error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": url, "error": err}).Error()
		return
	}
	req.Header.Set("Content-Type", "application/json")
	cli := http.DefaultClient
	resp, err := cli.Do(req)
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": url, "error": err}).Error()
		return
	}
	defer resp.Body.Close()
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.GetDefaultLogger().WithFields(logrus.Fields{"api": url, "error": err, "resp": resp}).Error()
	}
	return
}
