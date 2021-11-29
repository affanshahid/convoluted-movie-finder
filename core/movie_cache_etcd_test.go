package core

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/affanshahid/configo"
	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"
)

var etcdServer *embed.Etcd
var client *clientv3.Client

var someMovie *tmdb.MovieDetails = &tmdb.MovieDetails{
	ID:    1,
	Title: "Some Movie",
}

func init() {
	if err := configo.Initialize(
		os.DirFS("../config"),
		configo.WithDeploymentFromEnv("APP_ENV"),
	); err != nil {
		panic(err)
	}
}

func setupEtcd(t *testing.T) {
	lcurl, err := url.Parse("http://" + configo.MustGetString("etcd_url"))
	if err != nil {
		t.Fatal("unable to parse url", err)
	}

	cfg := embed.NewConfig()
	cfg.Dir, err = ioutil.TempDir("", "")
	if err != nil {
		t.Fatal("unable to create temp dir", err)
	}
	cfg.LCUrls = []url.URL{*lcurl}
	cfg.ACUrls = []url.URL{*lcurl}

	etcdServer, err = embed.StartEtcd(cfg)
	if err != nil {
		t.Fatal("unable to start etcd", err)
	}

	select {
	case <-etcdServer.Server.ReadyNotify():
	case <-time.After(10 * time.Second):
		etcdServer.Server.Stop()
		t.Fatal("timeout while starting etcd")
	}

	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{configo.MustGetString("etcd_url")},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatal("unable to start etcd", err)
	}
}

func TestSaveMovieDetails(t *testing.T) {
	setupEtcd(t)
	defer cleanupEtcd()

	cache, err := NewMovieCacheEtcd()
	assert.Nilf(t, err, "expected err to be nil")

	err = cache.SaveMovieDetails(someMovie)
	assert.Nilf(t, err, "expected err to be nil")

	expected, err := json.Marshal(someMovie)
	assert.Nilf(t, err, "expected err to be nil")

	result, err := client.Get(context.Background(), getMovieKey(someMovie.ID))
	assert.Nilf(t, err, "expected err to be nil")

	assert.Len(t, result.Kvs, 1, "expected non empty response")
	assert.Equal(t, expected, result.Kvs[0].Value)
}

func TestGetMovieDetails(t *testing.T) {
	setupEtcd(t)
	defer cleanupEtcd()

	cache, err := NewMovieCacheEtcd()
	assert.Nilf(t, err, "expected err to be nil")

	data, err := json.Marshal(someMovie)
	assert.Nilf(t, err, "expected err to be nil")

	_, err = client.Put(context.Background(), getMovieKey(someMovie.ID), string(data))
	assert.Nilf(t, err, "expected err to be nil")

	assert.Nilf(t, err, "expected err to be nil")

	result, err := cache.GetMovieDetails(someMovie.ID)
	assert.Nilf(t, err, "expected err to be nil")

	assert.Equal(t, someMovie, result)
}

func TestGetMovieDetailsReturnsNilIfNotSaved(t *testing.T) {
	setupEtcd(t)
	defer cleanupEtcd()

	cache, err := NewMovieCacheEtcd()
	assert.Nilf(t, err, "expected err to be nil")

	result, err := cache.GetMovieDetails(someMovie.ID)
	assert.Nilf(t, err, "expected err to be nil")

	assert.Equal(t, (*tmdb.MovieDetails)(nil), result)
}

func cleanupEtcd() {
	client.Close()
	etcdServer.Close()
	os.RemoveAll(etcdServer.Config().Dir)
}
