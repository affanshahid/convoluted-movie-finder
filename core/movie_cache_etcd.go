package core

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/affanshahid/configo"
	tmdb "github.com/cyruzin/golang-tmdb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const moviePrefix = "movie_"

type MovieCacheEtcd struct {
	client *clientv3.Client
}

func NewMovieCacheEtcd() (*MovieCacheEtcd, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{configo.MustGetString("etcd_url")},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		return nil, err
	}

	return &MovieCacheEtcd{client}, nil
}

func (c *MovieCacheEtcd) GetMovieDetails(id int64) (*tmdb.MovieDetails, error) {
	resp, err := c.client.Get(context.TODO(), getMovieKey(id))
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, nil
	}

	movie := new(tmdb.MovieDetails)
	err = json.Unmarshal(resp.Kvs[0].Value, movie)
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (c *MovieCacheEtcd) SaveMovieDetails(movie *tmdb.MovieDetails) error {
	data, err := json.Marshal(movie)
	if err != nil {
		return err
	}

	_, err = c.client.Put(context.TODO(), getMovieKey(movie.ID), string(data))
	if err != nil {
		return err
	}

	return nil
}

func (c *MovieCacheEtcd) Close() error {
	return c.client.Close()
}

func getMovieKey(id int64) string {
	return moviePrefix + strconv.FormatInt(id, 10)
}

var _ MovieCache = (*MovieCacheEtcd)(nil)
