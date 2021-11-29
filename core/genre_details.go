package core

import tmdb "github.com/cyruzin/golang-tmdb"

type GenrePeriodDetails struct {
	Id     int64
	Name   string
	Pct    float64
	Movies []*tmdb.MovieDetails
}
