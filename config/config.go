package config

import (
	"github.com/BurntSushi/toml"
	"github.com/badoll/movie_logger-backend/logger"
)

var c *Config

// init 初始化配置，path: project_path/config/config.toml
func init() {
	c = new(Config)
	confData, err := toml.DecodeFile("./config/config.toml", c)
	if err != nil {
		logger.GetDefaultLogger().Errorf("init config error: %v", err)
		panic(err)
	}
	logger.GetDefaultLogger().WithField("confData", confData).Debug("init config succ")
}

// Config ...
type Config struct {
	DBConf struct {
		User     string `toml:"User"`
		Password string `toml:"Password"`
		Host     string `toml:"Host"`
		Port     string `toml:"Port"`
		Database string `toml:"Database"`
		Charset  string `toml:"Charset"`
	} `toml:"DBConf"`
}

func GetConfig() Config {
	return *c
}
