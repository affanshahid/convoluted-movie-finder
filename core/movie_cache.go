//go:generate mockery --name=MovieCache

package core

import tmdb "github.com/cyruzin/golang-tmdb"

type MovieCache interface {
	GetMovieDetails(id int64) (*tmdb.MovieDetails, error)
	SaveMovieDetails(movie *tmdb.MovieDetails) error
}
