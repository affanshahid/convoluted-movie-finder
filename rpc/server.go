//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pb/server.proto

package rpc

import (
	"fmt"
	"net"

	"github.com/affanshahid/configo"
	"github.com/affanshahid/convoluted-movie-finder/core"
	"github.com/affanshahid/convoluted-movie-finder/rpc/pb"
	"google.golang.org/grpc"
)

func Serve(movieService *core.MovieService) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", configo.MustGetInt("grpc_port")))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterMovieServer(server, &movieServer{service: movieService})

	fmt.Printf("Server running on: %s\n", listener.Addr())
	if err := server.Serve(listener); err != nil {
		return err
	}

	return nil
}
