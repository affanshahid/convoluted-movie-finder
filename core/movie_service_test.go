package core

import (
	"errors"
	"testing"
	"time"

	"github.com/affanshahid/convoluted-movie-finder/core/mocks"
	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var genreList = &tmdb.GenreMovieList{
	Genres: []struct {
		ID   int64  "json:\"id\""
		Name string "json:\"name\""
	}{
		{28, "Action"},
		{29, "Sci-fi"},
	},
}

var allMoviesDiscover = &tmdb.DiscoverMovie{
	TotalPages:   1,
	Page:         1,
	TotalResults: 4,
	DiscoverMovieResults: &tmdb.DiscoverMovieResults{
		Results: []struct {
			VoteCount        int64   "json:\"vote_count\""
			ID               int64   "json:\"id\""
			Video            bool    "json:\"video\""
			VoteAverage      float32 "json:\"vote_average\""
			Title            string  "json:\"title\""
			Popularity       float32 "json:\"popularity\""
			PosterPath       string  "json:\"poster_path\""
			OriginalLanguage string  "json:\"original_language\""
			OriginalTitle    string  "json:\"original_title\""
			GenreIDs         []int64 "json:\"genre_ids\""
			BackdropPath     string  "json:\"backdrop_path\""
			Adult            bool    "json:\"adult\""
			Overview         string  "json:\"overview\""
			ReleaseDate      string  "json:\"release_date\""
		}{
			{
				ID:    1,
				Title: "Some Movie",
			},
			{
				ID:    2,
				Title: "Some Movie 1",
			},
			{
				ID:    3,
				Title: "Some Movie 2",
			},
			{
				ID:    4,
				Title: "Some Movie 3",
			},
		},
	},
}

var actionMoviesDiscover = &tmdb.DiscoverMovie{
	TotalPages:   1,
	Page:         1,
	TotalResults: 1,
	DiscoverMovieResults: &tmdb.DiscoverMovieResults{
		Results: []struct {
			VoteCount        int64   "json:\"vote_count\""
			ID               int64   "json:\"id\""
			Video            bool    "json:\"video\""
			VoteAverage      float32 "json:\"vote_average\""
			Title            string  "json:\"title\""
			Popularity       float32 "json:\"popularity\""
			PosterPath       string  "json:\"poster_path\""
			OriginalLanguage string  "json:\"original_language\""
			OriginalTitle    string  "json:\"original_title\""
			GenreIDs         []int64 "json:\"genre_ids\""
			BackdropPath     string  "json:\"backdrop_path\""
			Adult            bool    "json:\"adult\""
			Overview         string  "json:\"overview\""
			ReleaseDate      string  "json:\"release_date\""
		}{
			{
				ID:    1,
				Title: "Some Movie",
			},
		},
	},
}

var scifiMoviesDiscover = &tmdb.DiscoverMovie{
	TotalPages:   1,
	Page:         1,
	TotalResults: 1,
	DiscoverMovieResults: &tmdb.DiscoverMovieResults{
		Results: []struct {
			VoteCount        int64   "json:\"vote_count\""
			ID               int64   "json:\"id\""
			Video            bool    "json:\"video\""
			VoteAverage      float32 "json:\"vote_average\""
			Title            string  "json:\"title\""
			Popularity       float32 "json:\"popularity\""
			PosterPath       string  "json:\"poster_path\""
			OriginalLanguage string  "json:\"original_language\""
			OriginalTitle    string  "json:\"original_title\""
			GenreIDs         []int64 "json:\"genre_ids\""
			BackdropPath     string  "json:\"backdrop_path\""
			Adult            bool    "json:\"adult\""
			Overview         string  "json:\"overview\""
			ReleaseDate      string  "json:\"release_date\""
		}{
			{
				ID:    2,
				Title: "Some Movie 1",
			},
			{
				ID:    3,
				Title: "Some Movie 2",
			},
			{
				ID:    4,
				Title: "Some Movie 3",
			},
		},
	},
}

var someMovieDetails = &tmdb.MovieDetails{
	ID:      1,
	Title:   "Some Movie",
	Revenue: 1000,
}

var cachedMovieDetails = &tmdb.MovieDetails{
	ID:      1,
	Title:   "Some Movie",
	Revenue: 10000,
}

var someMovie1Details = &tmdb.MovieDetails{
	ID:      2,
	Title:   "Some Movie 1",
	Revenue: 1000,
}

var someMovie2Details = &tmdb.MovieDetails{
	ID:      3,
	Title:   "Some Movie 2",
	Revenue: 1000,
}

var someMovie3Details = &tmdb.MovieDetails{
	ID:      4,
	Title:   "Some Movie 3",
	Revenue: 1000,
}

var startDate = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
var endDate = time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)

func TestFetchGenrePeriodDetailsWithRevenueFilter(t *testing.T) {
	t.Parallel()
	mockClient := new(mocks.TmdbClient)
	mockCache := new(mocks.MovieCache)

	var nilmap map[string]string

	mockCache.On("GetMovieDetails", mock.Anything).Return(nil, nil)
	mockCache.On("SaveMovieDetails", mock.Anything).Return(nil)

	mockClient.On("GetGenreMovieList", nilmap).Return(genreList, nil)
	mockClient.On("GetDiscoverMovie", map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
	}).Return(allMoviesDiscover, nil)
	mockClient.On("GetDiscoverMovie", map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
		"with_genres":      "28",
	}).Return(actionMoviesDiscover, nil)
	mockClient.On("GetDiscoverMovie", map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
		"with_genres":      "28",
		"page":             "1",
	}).Return(actionMoviesDiscover, nil)
	mockClient.On("GetMovieDetails", 1, nilmap).Return(someMovieDetails, nil)

	expected := GenrePeriodDetails{
		Id:     28,
		Name:   "Action",
		Pct:    25,
		Movies: []*tmdb.MovieDetails{someMovieDetails},
	}

	svc := NewMovieService(mockClient, mockCache)

	result, err := svc.FetchGenrePeriodDetailsWithRevenueFilter(expected.Id, startDate, endDate, 1, OpGt)
	assert.Nilf(t, err, "expected error to be nil")
	assert.Equal(t, expected, result)
}

func TestFetchGenrePeriodDetailsWithRevenueFilterErrorsWithInvalidGenre(t *testing.T) {
	t.Parallel()
	mockClient := new(mocks.TmdbClient)
	mockCache := new(mocks.MovieCache)
	var nilmap map[string]string

	mockCache.On("GetMovieDetails", mock.Anything).Return(nil, nil)
	mockCache.On("SaveMovieDetails", mock.Anything).Return(nil)

	mockClient.On("GetGenreMovieList", nilmap).Return(genreList, nil)

	svc := NewMovieService(mockClient, mockCache)

	_, err := svc.FetchGenrePeriodDetailsWithRevenueFilter(21, startDate, endDate, 1, OpGt)
	assert.NotNil(t, err, "expected error to not be nil")
	assert.Equal(t, ErrGenreNotFound, err)
}

func TestFetchGenrePeriodDetailsWithRevenueFilterErrorsWhenTotalNotAvailable(t *testing.T) {
	t.Parallel()
	mockClient := new(mocks.TmdbClient)
	mockCache := new(mocks.MovieCache)
	var nilmap map[string]string

	mockCache.On("GetMovieDetails", mock.Anything).Return(nil, nil)
	mockCache.On("SaveMovieDetails", mock.Anything).Return(nil)

	expectedError := errors.New("some error occurred")

	mockClient.On("GetGenreMovieList", nilmap).Return(genreList, nil)
	mockClient.On("GetDiscoverMovie", map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
	}).Return(nil, expectedError)
	mockClient.On("GetDiscoverMovie", map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
		"with_genres":      "28",
	}).Return(actionMoviesDiscover, nil)
	mockClient.On("GetDiscoverMovie", map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
		"with_genres":      "28",
		"page":             "1",
	}).Return(actionMoviesDiscover, nil)
	mockClient.On("GetMovieDetails", 1, nilmap).Return(someMovieDetails, nil)

	svc := NewMovieService(mockClient, mockCache)

	_, err := svc.FetchGenrePeriodDetailsWithRevenueFilter(28, startDate, endDate, 1, OpGt)
	assert.NotNil(t, err, "expected error to not be nil")
	assert.Equal(t, expectedError, err)
}

func TestFetchGenrePeriodDetailsWithRevenueFilterErrorsWhenMovieDetailsNotAvailable(t *testing.T) {
	t.Parallel()
	mockClient := new(mocks.TmdbClient)
	mockCache := new(mocks.MovieCache)
	var nilmap map[string]string

	mockCache.On("GetMovieDetails", mock.Anything).Return(nil, nil)
	mockCache.On("SaveMovieDetails", mock.Anything).Return(nil)

	expectedError := errors.New("some error occurred")

	mockClient.On("GetGenreMovieList", nilmap).Return(genreList, nil)
	mockClient.On("GetDiscoverMovie", map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
	}).Return(allMoviesDiscover, nil)
	mockClient.On("GetDiscoverMovie", map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
		"with_genres":      "29",
	}).Return(scifiMoviesDiscover, nil)
	mockClient.On("GetDiscoverMovie", map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
		"with_genres":      "29",
		"page":             "1",
	}).Return(scifiMoviesDiscover, nil)
	mockClient.On("GetMovieDetails", 2, nilmap).Return(nil, expectedError)
	mockClient.On("GetMovieDetails", 3, nilmap).Return(someMovie2Details, nil)
	mockClient.On("GetMovieDetails", 4, nilmap).Return(someMovie3Details, nil)

	svc := NewMovieService(mockClient, mockCache)

	_, err := svc.FetchGenrePeriodDetailsWithRevenueFilter(29, startDate, endDate, 1, OpGt)
	assert.NotNil(t, err, "expected error to not be nil")
	assert.Equal(t, expectedError, err)
}

func TestFetchGenrePeriodDetailsWithRevenueFilterUsingCache(t *testing.T) {
	t.Parallel()
	mockClient := new(mocks.TmdbClient)
	mockCache := new(mocks.MovieCache)

	var nilmap map[string]string

	mockCache.On("GetMovieDetails", int64(1)).Return(cachedMovieDetails, nil)
	mockCache.On("SaveMovieDetails", mock.Anything).Return(nil)

	mockClient.On("GetGenreMovieList", nilmap).Return(genreList, nil)
	mockClient.On("GetDiscoverMovie", map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
	}).Return(allMoviesDiscover, nil)
	mockClient.On("GetDiscoverMovie", map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
		"with_genres":      "28",
	}).Return(actionMoviesDiscover, nil)
	mockClient.On("GetDiscoverMovie", map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
		"with_genres":      "28",
		"page":             "1",
	}).Return(actionMoviesDiscover, nil)
	mockClient.On("GetMovieDetails", 1, nilmap).Return(someMovieDetails, nil)

	expected := GenrePeriodDetails{
		Id:     28,
		Name:   "Action",
		Pct:    25,
		Movies: []*tmdb.MovieDetails{cachedMovieDetails},
	}

	svc := NewMovieService(mockClient, mockCache)

	result, err := svc.FetchGenrePeriodDetailsWithRevenueFilter(expected.Id, startDate, endDate, 9999, OpGt)
	assert.Nilf(t, err, "expected error to be nil")
	assert.Equal(t, expected, result)
}

func TestFetchGenrePeriodDetailsWithRevenueFilterUpdatesCache(t *testing.T) {
	t.Parallel()
	mockClient := new(mocks.TmdbClient)
	mockCache := new(mocks.MovieCache)

	var nilmap map[string]string

	mockCache.On("GetMovieDetails", int64(1)).Return(nil, nil)
	mockCache.On("SaveMovieDetails", mock.Anything).Return(nil)

	mockClient.On("GetGenreMovieList", nilmap).Return(genreList, nil)
	mockClient.On("GetDiscoverMovie", map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
	}).Return(allMoviesDiscover, nil)
	mockClient.On("GetDiscoverMovie", map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
		"with_genres":      "28",
	}).Return(actionMoviesDiscover, nil)
	mockClient.On("GetDiscoverMovie", map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
		"with_genres":      "28",
		"page":             "1",
	}).Return(actionMoviesDiscover, nil)
	mockClient.On("GetMovieDetails", 1, nilmap).Return(someMovieDetails, nil)

	expected := GenrePeriodDetails{
		Id:     28,
		Name:   "Action",
		Pct:    25,
		Movies: []*tmdb.MovieDetails{someMovieDetails},
	}

	svc := NewMovieService(mockClient, mockCache)

	result, err := svc.FetchGenrePeriodDetailsWithRevenueFilter(expected.Id, startDate, endDate, 1, OpGt)
	assert.Nilf(t, err, "expected error to be nil")
	assert.Equal(t, expected, result)
	mockCache.AssertCalled(t, "SaveMovieDetails", someMovieDetails)
}

func TestFetchGenrePeriodDetailsWithRevenueFilterErrorsWhenCacheGetFails(t *testing.T) {
	t.Parallel()
	mockClient := new(mocks.TmdbClient)
	mockCache := new(mocks.MovieCache)
	var nilmap map[string]string

	expectedError := errors.New("some error occurred")

	mockCache.On("GetMovieDetails", mock.Anything).Return(nil, expectedError)
	mockCache.On("SaveMovieDetails", mock.Anything).Return(nil)

	mockClient.On("GetGenreMovieList", nilmap).Return(genreList, nil)
	mockClient.On("GetDiscoverMovie", map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
	}).Return(allMoviesDiscover, nil)
	mockClient.On("GetDiscoverMovie", map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
		"with_genres":      "29",
	}).Return(scifiMoviesDiscover, nil)
	mockClient.On("GetDiscoverMovie", map[string]string{
		"release_date.gte": startDate.Format(timeFormat),
		"release_date.lte": endDate.Format(timeFormat),
		"with_genres":      "29",
		"page":             "1",
	}).Return(scifiMoviesDiscover, nil)
	mockClient.On("GetMovieDetails", 2, nilmap).Return(someMovie1Details, nil)
	mockClient.On("GetMovieDetails", 3, nilmap).Return(someMovie2Details, nil)
	mockClient.On("GetMovieDetails", 4, nilmap).Return(someMovie3Details, nil)

	svc := NewMovieService(mockClient, mockCache)

	_, err := svc.FetchGenrePeriodDetailsWithRevenueFilter(29, startDate, endDate, 1, OpGt)
	assert.NotNil(t, err, "expected error to not be nil")
	assert.Equal(t, expectedError, err)
}
