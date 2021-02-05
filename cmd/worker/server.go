package main

import (
	"context"
	"github.com/ng5gc/uegnbsim/internal/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
)

type Server struct {
	api.UnimplementedSimWorkerServer
}

func StartWorkerGrpcServer(ctx context.Context, srvPort int) error {
	addr := ":" + strconv.Itoa(srvPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrapf(err, "Failed to listen at address:%v", addr)
	}

	s := grpc.NewServer()
	api.RegisterSimWorkerServer(s, &Server{})
	reflection.Register(s)

	// Control gRPC server lifecycle
	go func() {
		<-ctx.Done()
		infoLog.Print("Got notification to stop CLI gRPC server.")
		s.Stop()
		return
	}()

	log.Printf("%v started! Listen at %v\n", "GRPC Server", addr)
	if err := s.Serve(lis); err != nil {
		return errors.Wrap(err, "Failed to start gRPC server")
	}
	return nil
}
