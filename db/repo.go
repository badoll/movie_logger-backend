package db

import "fmt"

var movie_col = "m.douban_id,m.title,m.poster,m.cate,m.director," +
	"m.writer,m.performer,m.region,m.language,m.release_year," +
	"m.release_date,m.runtime,m.rating_score,m.rating_num," +
	"m.rating_stars,m.intro,m.main_cast,m.photos"

func (db *DB) GetMovieDetailByDoubanID(doubanID string) (Movie, error) {
	movie := Movie{}
	sql := "select " + movie_col + " from movie m where m.douban_id = ? limit 1"
	err := db.Get(&movie, sql, doubanID)
	return movie, err
}

// GetMovieListByTitle
//使用索引模糊搜索前缀，后面再优化模糊搜索中间内容
func (db *DB) GetMovieListByTitle(title string, limit, offset int) (mlist []Movie, err error) {
	title = title + "%"
	sql := "select " + movie_col + " from movie m where m.title like ? order by m.rating_score desc limit ? offset ?"
	err = db.Select(&mlist, sql, title, limit, offset)
	return

}

// GetMovieListBychart 排行榜电影id联表查电影详情
func (db *DB) GetMovieListByChart(chart string) (mlist []Movie, err error) {
	sql := fmt.Sprintf("select "+movie_col+" from %s as c join movie as m on c.douban_id = m.douban_id", chart)
	err = db.Select(&mlist, sql)
	return
}

// GetMovieListByChartWithPage 排行榜电影id联表查电影详情（分页）
func (db *DB) GetMovieListByChartWithPage(chart string, limit, offset int) (mlist []Movie, err error) {
	sql := fmt.Sprintf("select "+movie_col+" from %s as c join movie as m on c.douban_id = m.douban_id limit ? offset ?",
		chart)
	err = db.Select(&mlist, sql, limit, offset)
	return
}

// GetRecommendNumByUser 推荐池里可推荐的电影总数
func (db *DB) GetRecommendNumByUser(userID string) (total int, err error) {
	sql := "select count(*) from recommend rc join movie m on rc.movie_id = m.id " +
		"and rc.user_id = ? and rc.has_recom = 0"
	err = db.Get(&total, sql, userID)
	return
}

// GetRecommendByUser 根据用户open_id拿到推荐的电影list
func (db *DB) GetRecommendByUser(userID string) (mlist []Movie, err error) {
	// TODO 先返回10个高分电影
	sql := "select " + movie_col + " from movie m order by rating_score desc limit 10"
	err = db.Select(&mlist, sql)
	return
	// sql := "select " + movie_col + " from recommend rc join movie m on rc.movie_id = m.id " +
	// 	"and rc.user_id = ? and rc.has_recom = 0"
	// err = db.Select(&mlist, sql, userID)
	// return
}

// GetRecommendByMovie 电影相关推荐
func (db *DB) GetRecommendByMovie(doubanID string) (mlist []Movie, err error) {
	// TODO 先返回10个高分电影
	sql := "select " + movie_col + " from movie m order by rating_score desc limit 10"
	err = db.Select(&mlist, sql)
	return
}
