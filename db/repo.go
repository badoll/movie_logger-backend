package db

func (db *DB) GetMovieDetailByDoubanID(doubanID string) (Movie, error) {
	movie := Movie{}
	sql := "select " +
		"douban_id,title,poster,cate,director," +
		"writer,performer,region,language,release_year," +
		"release_date,runtime,rating_score,rating_num," +
		"rating_stars,intro,main_cast,photos " +
		"from movie where douban_id = ? limit 1"
	err := db.Get(&movie, sql, doubanID)
	return movie, err
}

func (db *DB) GetMovieListByTitle(title string) ([]Movie, error) {
	title = title + "%"
	var movieList []Movie
	sql := "select " +
		"douban_id,title,poster,cate,director," +
		"writer,performer,region,language,release_year," +
		"release_date,runtime,rating_score,rating_num," +
		"rating_stars,intro,main_cast,photos " +
		"from movie where title like ?"
	err := db.Select(&movieList, sql, title)
	return movieList, err
}
