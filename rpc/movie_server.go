package rpc

import (
	"context"

	"github.com/affanshahid/convoluted-movie-finder/core"
	"github.com/affanshahid/convoluted-movie-finder/rpc/pb"
)

type movieServer struct {
	pb.UnimplementedMovieServer
	service *core.MovieService
}

func (s *movieServer) FetchGenrePeriodDetails(
	ctx context.Context,
	in *pb.GenrePeriodDetailsRequest,
) (*pb.GenrePeriodDetailsReply, error) {
	resp, err := s.service.FetchGenrePeriodDetailsWithRevenueFilter(
		in.GenreId,
		in.StartDate.AsTime(),
		in.EndDate.AsTime(),
		in.Revenue,
		core.Operator(in.RevenueCheckOperator),
	)
	if err != nil {
		return nil, err
	}

	reply := pb.GenrePeriodDetailsReply{
		GenreId: resp.Id,
		Name:    resp.Name,
		Pct:     float32(resp.Pct),
		Movies:  []*pb.MovieMsg{},
	}

	for _, movie := range resp.Movies {
		reply.Movies = append(reply.Movies, &pb.MovieMsg{
			Id:          movie.ID,
			Title:       movie.Title,
			ReleaseDate: movie.ReleaseDate,
			Revenue:     movie.Revenue,
		})
	}

	return &reply, nil
}
