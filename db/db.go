package db

import (
	"fmt"

	"github.com/badoll/movie_logger-backend/config"
	"github.com/badoll/movie_logger-backend/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var cli *DB

// DB 封装db客户端
type DB struct {
	*sqlx.DB
}

func Init() {
	conf := config.GetConfig().DBConf
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", conf.User,
		conf.Password, conf.Host, conf.Port, conf.Database, conf.Charset)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		logger.GetDefaultLogger().Error("init db error: %v", err)
		panic(err)
	}
	logger.GetDefaultLogger().WithField("db", db).Debug("init db succ")
	cli = new(DB)
	cli.DB = db
}

func GetCli() *DB {
	return cli
}
