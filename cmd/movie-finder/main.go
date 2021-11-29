package main

import (
	"os"

	"github.com/affanshahid/configo"
	"github.com/affanshahid/convoluted-movie-finder/core"
	"github.com/affanshahid/convoluted-movie-finder/rpc"
	tmdb "github.com/cyruzin/golang-tmdb"
)

func init() {
	if err := configo.Initialize(os.DirFS("./config")); err != nil {
		panic(err)
	}
}

func main() {
	c, err := tmdb.Init(configo.MustGetString("tmdb_api_key"))
	if err != nil {
		panic(err)
	}

	cache, err := core.NewMovieCacheEtcd()
	if err != nil {
		panic(err)
	}
	s := core.NewMovieService(c, cache)

	err = rpc.Serve(s)
	if err != nil {
		panic(err)
	}
}
