//go:generate mockery --name=TmdbClient

package core

import tmdb "github.com/cyruzin/golang-tmdb"

type TmdbClient interface {
	GetGenreMovieList(urlOptions map[string]string) (*tmdb.GenreMovieList, error)
	GetDiscoverMovie(urlOptions map[string]string) (*tmdb.DiscoverMovie, error)
	GetMovieDetails(id int, urlOptions map[string]string) (*tmdb.MovieDetails, error)
}
