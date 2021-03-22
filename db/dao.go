package db

import "strings"

// Movie DB raw data
type Movie struct {
	DoubanID    string `db:"douban_id"`
	Title       string `db:"title"`
	Poster      string `db:"poster"`
	Cate        string `db:"cate"`
	Director    string `db:"director"`
	Writer      string `db:"writer"`
	Performer   string `db:"performer"`
	Region      string `db:"region"`
	Language    string `db:"language"`
	ReleaseYear string `db:"release_year"`
	ReleaseDate string `db:"release_date"`
	Runtime     string `db:"runtime"`
	RatingScore uint   `db:"rating_score"`
	RatingNum   uint   `db:"rating_num"`
	RatingStars string `db:"rating_stars"`
	Intro       string `db:"intro"`
	MainCast    string `db:"main_cast"`
	Photos      string `db:"photos"`
}

// SplitString 分隔以','为分隔符的string数据
func SplitString(data string) []string {
	if len(data) == 0 {
		return []string{}
	}
	return strings.Split(data, ",")
}

type User struct {
}
