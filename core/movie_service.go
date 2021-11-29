package core

import (
	"errors"
	"strconv"
	"time"

	tmdb "github.com/cyruzin/golang-tmdb"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	"golang.org/x/sync/errgroup"
)

const timeFormat = "2006-01-02"

var ErrGenreNotFound = errors.New("genre not found")

type Operator uint8

const (
	OpLt Operator = iota
	OpEq
	OpGt
)

type multiMovieDetailsMsg struct {
	movies []*tmdb.MovieDetails
	ack    chan bool
}

type movieDetailsMsg struct {
	movie *tmdb.MovieDetails
	ack   chan bool
}

type totalMsg struct {
	total int64
	err   error
}

type MovieService struct {
	client TmdbClient
	cache  MovieCache
}

func NewMovieService(client TmdbClient, cache MovieCache) *MovieService {
	return &MovieService{client, cache}
}

func (s *MovieService) FetchGenrePeriodDetailsWithRevenueFilter(
	genreId int64,
	startDate time.Time,
	endDate time.Time,
	revenue int64,
	revenueCheckOperator Operator,
) (GenrePeriodDetails, error) {
	var genreDetails GenrePeriodDetails
	genreDetails.Id = genreId

	genreResult, err := s.client.GetGenreMovieList(nil)
	if err != nil {
		return genreDetails, err
	}

	found := false
	for _, g := range genreResult.Genres {
		if g.ID == genreId {
			genreDetails.Name = g.Name
			found = true
			break
		}
	}

	if !found {
		return genreDetails, ErrGenreNotFound
	}

	totalMsgChan := make(chan totalMsg)

	go func() {
		total, err := s.getTotalMoviesInPeriod(startDate, endDate)
		totalMsgChan <- totalMsg{total, err}
	}()

	result, err := s.client.GetDiscoverMovie(map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
		"with_genres":      strconv.FormatInt(genreId, 10),
	})

	if err != nil {
		return genreDetails, err
	}

	totalPages := result.TotalPages

	var eg errgroup.Group
	moviesChan := make(chan multiMovieDetailsMsg, 10)

	for i := int64(1); i <= totalPages; i++ {
		page := i
		ackChannel := make(chan bool)
		eg.Go(func() error {
			movies, err := s.getMovieDetailsFromPage(
				genreId,
				startDate,
				endDate,
				page,
				revenue,
				revenueCheckOperator,
			)

			if err != nil {
				return err
			}

			moviesChan <- multiMovieDetailsMsg{movies, ackChannel}
			<-ackChannel

			return nil
		})
	}

	go func() {
		for {
			msg := <-moviesChan
			genreDetails.Movies = append(genreDetails.Movies, msg.movies...)
			msg.ack <- true
		}
	}()

	totalResult := <-totalMsgChan
	if totalResult.err != nil {
		return genreDetails, totalResult.err
	}

	err = eg.Wait()
	if err != nil {
		return genreDetails, err
	}

	genreDetails.Pct = (float64(len(genreDetails.Movies)) / float64(totalResult.total)) * 100

	return genreDetails, nil
}

func (s *MovieService) getMovieDetailsFromPage(
	genreId int64,
	start, end time.Time,
	page int64,
	revenue int64,
	revenueCheckOperator Operator,
) (ret []*tmdb.MovieDetails, err error) {
	result, err := s.client.GetDiscoverMovie(map[string]string{
		"release_date.gte": start.Format(timeFormat),
		"release_date.lte": end.Format(timeFormat),
		"with_genres":      strconv.FormatInt(genreId, 10),
		"page":             strconv.FormatInt(page, 10),
	})

	if err != nil {
		return nil, err
	}

	movieChan := make(chan movieDetailsMsg, 10)
	var eg errgroup.Group

	for _, movie := range result.Results {
		lMovie := movie
		ack := make(chan bool)
		eg.Go(func() error {
			var result *tmdb.MovieDetails
			cachedMovie, err := s.cache.GetMovieDetails(lMovie.ID)
			if err != nil && !isConnectivityError(err) {
				return err
			}

			if cachedMovie == nil && err == nil {
				result, err = s.client.GetMovieDetails(int(lMovie.ID), nil)
				if err != nil {
					return err
				}

				err = s.cache.SaveMovieDetails(result)
				if err != nil && !isConnectivityError(err) {
					return err
				}
			} else {
				result = cachedMovie
			}

			switch {
			case revenueCheckOperator == OpGt && result.Revenue > revenue:
				fallthrough
			case revenueCheckOperator == OpLt && result.Revenue < revenue:
				fallthrough
			case revenueCheckOperator == OpEq && result.Revenue == revenue:
				movieChan <- movieDetailsMsg{result, ack}
				<-ack
			}

			return nil
		})
	}

	go func() {
		for {
			msg := <-movieChan
			ret = append(ret, msg.movie)
			msg.ack <- true
		}
	}()

	err = eg.Wait()

	return ret, err
}

func (s *MovieService) getTotalMoviesInPeriod(start, end time.Time) (int64, error) {
	result, err := s.client.GetDiscoverMovie(map[string]string{
		"release_date.gte": start.Format(timeFormat),
		"release_date.lte": end.Format(timeFormat),
	})

	if err != nil {
		return 0, err
	}

	return result.TotalResults, nil
}

func isConnectivityError(err error) bool {
	switch err {
	case rpctypes.ErrTimeout:
		fallthrough
	case rpctypes.ErrTimeoutDueToLeaderFail:
		fallthrough
	case rpctypes.ErrGRPCTimeoutDueToConnectionLost:
		return true
	default:
		return false
	}
}
