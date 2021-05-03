package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

/**************************************Movie*************************************************/

var movie_col = "m.id,m.title,m.poster,m.cate,m.director," +
	"m.writer,m.performer,m.region,m.language,m.release_year," +
	"m.release_date,m.runtime,m.rating_score,m.rating_num," +
	"m.rating_stars,m.intro,m.main_cast,m.photos"

var movie_index_col = "m.id,m.title,m.cate,m.director," +
	"m.writer,m.performer,m.region,m.language,m.release_year," +
	"m.rating_score"

func (db *DB) GetMovieDetailByMovieID(movieID int64) (Movie, error) {
	movie := Movie{}
	sql := "select " + movie_col + " from movie m where m.id = ? limit 1"
	err := db.Get(&movie, sql, movieID)
	return movie, err
}

func (db *DB) SelectMovieDetailByMovieIDList(movieIDs []int64) (mlist []Movie, err error) {
	if len(movieIDs) == 0 {
		return
	}
	sql := "select " + movie_col + " from movie m where m.id in (?)"
	query, args, err := sqlx.In(sql, movieIDs)
	if err != nil {
		return
	}
	rawList := []Movie{}
	err = db.Select(&rawList, query, args...)
	if err != nil {
		return
	}
	// 按movieIDs中顺序排序
	m := map[int64]Movie{}
	for _, v := range rawList {
		m[v.MovieID] = v
	}
	for _, id := range movieIDs {
		if movie, ok := m[id]; ok {
			mlist = append(mlist, movie)
		}
	}
	return
}

// GetMovieListByTitle
//使用索引模糊搜索前缀，后面再优化模糊搜索中间内容
func (db *DB) GetMovieListByTitle(title string, limit, offset int) (mlist []Movie, err error) {
	title = title + "%"
	sql := "select " + movie_col + " from movie m where m.title like ? order by m.rating_score desc limit ? offset ?"
	err = db.Select(&mlist, sql, title, limit, offset)
	return

}

// GetMovieListByChart 排行榜电影id联表查电影详情
// 评分人数 > 评分 > 上映日期 > 上映年份 > id
func (db *DB) GetMovieListByChart(chart string) (mlist []Movie, err error) {
	sql := fmt.Sprintf("select "+movie_col+" from %s as c join movie as m on c.douban_id = m.douban_id "+
		"order by m.rating_num desc, m.rating_score desc, m.release_date desc, m.release_year desc, m.douban_id desc", chart)
	err = db.Select(&mlist, sql)
	return
}

// GetMovieListByChartWithPage 排行榜电影id联表查电影详情（分页）
// 评分人数 > 评分 > 上映日期 > 上映年份 > id
func (db *DB) GetMovieListByChartWithPage(chart string, limit, offset int) (mlist []Movie, err error) {
	sql := fmt.Sprintf("select "+movie_col+" from %s as c join movie as m on c.douban_id = m.douban_id "+
		"order by m.rating_num desc, m.rating_score desc, m.release_date desc, m.release_year desc, m.douban_id desc "+
		"limit ? offset ?", chart)
	err = db.Select(&mlist, sql, limit, offset)
	return
}

// GetRecommendNumByUser 推荐池里可推荐的电影总数
func (db *DB) GetRecommendNumByUser(userID int64) (total int, err error) {
	sql := "select count(*) from recommend rc join movie m on rc.movie_id = m.id " +
		"and rc.user_id = ? and rc.has_recom = 0"
	err = db.Get(&total, sql, userID)
	return
}

// GetTopMovieByPage 高分电影 分页
func (db *DB) GetTopMovieWithPage(limit, offset int) (mlist []Movie, err error) {
	sql := "select " + movie_col + " from movie m order by m.rating_score desc limit ? offset ?"
	err = db.Select(&mlist, sql, limit, offset)
	return
}

func (db *DB) GetMovieIndexData(limit, offset int) (mlist []MovieIndex, err error) {
	sql := "select " + movie_index_col + " from movie m order by m.id limit ? offset ?"
	err = db.Select(&mlist, sql, limit, offset)
	return
}

func (db *DB) GetMovieCount() (cnt int, err error) {
	sql := "select count(*) from movie"
	err = db.Get(&cnt, sql)
	return
}

/**************************************User**************************************************/

// GetUserInfo 根据open_id获取user主键id，若无此用户则创建
func (db *DB) GetUserInfo(openID string) (user User, err error) {
	tx, err := db.Beginx()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	u := make([]User, 0)
	sql := "select id, nick_name, avatar_url from user where open_id = ?"
	err = tx.Select(&u, sql, openID)
	if err != nil {
		return
	}
	if len(u) > 0 {
		return u[0], nil
	}
	// 没有则创建
	sql = "insert into user (open_id) values (?)"
	result, err := tx.Exec(sql, openID)
	if err != nil {
		return
	}
	userID, err := result.LastInsertId()
	user.ID = userID
	return
}

// IsNewUser 判断是否是新用户（没有选择兴趣类型）
func (db *DB) IsNewUser(userID int64) (isNew bool, err error) {
	isNew = true
	sql := "select count(*) from user where id = ? and length(inter_field) != 0"
	cnt := 0
	err = db.Get(&cnt, sql, userID)
	if err != nil {
		return
	}
	if cnt > 0 {
		isNew = false
	}
	return
}

// SetUserInterField 设置用户感兴趣的电影类型
func (db *DB) SetUserInterField(userID int64, interField string) error {
	sql := "insert into user (id, inter_field) values (?,?) on duplicate key update " +
		"inter_field = values(inter_field)"
	_, err := db.Exec(sql, userID, interField)
	return err
}

// LikeMovie 设置用户喜欢或不喜欢
func (db *DB) LikeMovie(userID, movieID int64, likeStatus bool) error {
	like := 0
	if likeStatus {
		like = 1
	}
	sql := "insert into user_collect (user_id, movie_id, likes) values (?, ?, ?) " +
		"on duplicate key update likes = values(likes)"
	_, err := db.Exec(sql, userID, movieID, like)
	return err
}

// GetUserLikeList 用户喜欢的电影，按点赞时间排序
func (db *DB) GetUserLikeList(userID int64) (ids []int64, err error) {
	sql := "select movie_id from user_collect where user_id = ? and likes = 1 order by update_time desc"
	err = db.Select(&ids, sql, userID)
	return
}

// GetUserLikeListWithPage 用户喜欢的电影，按点赞时间排序(分页)
func (db *DB) GetUserLikeListWithPage(userID int64, limit, offset int) (ids []int64, err error) {
	sql := "select movie_id from user_collect where user_id = ? and likes = 1 order by update_time desc limit ? offset ?"
	err = db.Select(&ids, sql, userID, limit, offset)
	return
}

// GetUesrInter 获取用户兴趣
func (db *DB) GetUesrInter(userID int64) (userInter UserInter, err error) {
	sql := "select inter_field, inter_director, inter_writer, inter_performer from user where id = ?"
	err = db.Get(&userInter, sql, userID)
	return
}

func (db *DB) UpdateUserInter(userID int64, userInter UserInter) error {
	sql := "update user set inter_field = ?, inter_director = ?, inter_writer = ?, inter_performer = ? " +
		"where id = ?"
	_, err := db.Exec(sql, userInter.InterField, userInter.InterDirector, userInter.InterWriter,
		userInter.InterPerformer, userID)
	return err
}

func (db *DB) UpdateUserInfo(userID int64, userInfo User) error {
	sql := "update user set nick_name = ?, avatar_url = ? where id = ?"
	_, err := db.Exec(sql, userInfo.NickName, userInfo.AvatarUrl, userID)
	return err
}
