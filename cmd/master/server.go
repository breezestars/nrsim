package main

import (
	"context"
	"github.com/ng5gc/uegnbsim/internal/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

type MasterServer struct {
	api.UnimplementedSimMasterServer
	workers sync.Map
}

func (s *MasterServer) Register(ctx context.Context, cfg *api.RegisterConfig) (*emptypb.Empty, error) {
	// Get worker IP
	ip := net.ParseIP(cfg.IP)

	// Delete first and re-store
	if _, found := s.workers.Load(ip.String()); found {
		s.workers.Delete(ip.String())
		infoLog.Printf("Duplicated worker: %v found, will delete and re-store", ip.String())
	}
	s.workers.Store(ip.String(), ip)

	infoLog.Printf("IP: %v has finished registration", ip.String())
	return &emptypb.Empty{}, nil
}

func (s *MasterServer) StreamChannel(stream api.SimMaster_StreamChannelServer) error {
	var (
		wg     sync.WaitGroup
		err    error
		ctx    context.Context
		cancel context.CancelFunc
	)

	ctx, cancel = context.WithCancel(context.Background())

	// Parse metadata
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		errLog.Printf("Parsed incoming context failed")
	}

	// Get worker IP
	ip := net.ParseIP(md.Get("IP")[0])

	// Delete first and re-store
	if _, found := s.workers.Load(ip.String()); found {
		s.workers.Delete(ip.String())
		infoLog.Printf("Duplicated worker: %v found, will delete and re-store", ip.String())
	}
	s.workers.Store(ip.String(), ip)

	infoLog.Printf("IP: %v has finished registration", ip.String())

	// TODO: use errgroup to refactor.

	// Send heartbeat
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err = stream.Send(&emptypb.Empty{})
				time.Sleep(time.Second * 2)
			}

		}
	}()

	// Monitor connection
	wg.Add(1)
	go func(err error) {
		defer func() {
			cancel()
			wg.Done()
		}()

		for {
			_, err = stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					// TODO: The connection is EOF, print log and remove worker from workers pool.
					infoLog.Printf("Connection %v EOF: %v", ip, err)
					return
				} else {
					// TODO: Print log, remove worker from workers pool and return error.
					errLog.Printf("Connection receive error: %v", err)
					return
				}
			}
			debugLog.Printf("Conn is normal")
		}
	}(err)

	wg.Wait()
	return err
}

func StartMasterGrpcServer(ctx context.Context, srvPort int) error {
	addr := ":" + strconv.Itoa(srvPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrapf(err, "Failed to listen at address: %v", addr)
	}

	s := grpc.NewServer()
	api.RegisterSimMasterServer(s, &MasterServer{})
	reflection.Register(s)

	// Control gRPC server lifecycle
	go func() {
		<-ctx.Done()
		infoLog.Print("Got notification to graceful stop Master gRPC server.")
		s.Stop()
		return
	}()

	infoLog.Printf("%v started! Listen at %v", "Master gRPC Server", addr)
	if err := s.Serve(lis); err != nil {
		return errors.Wrap(err, "Failed to start master gRPC server")
	}
	return nil
}

type CLIServer struct {
	api.UnimplementedSimCliServer
}

func StartCLIGrpcServer(ctx context.Context, srvPort int) error {
	addr := ":" + strconv.Itoa(srvPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrapf(err, "Failed to listen at address: %v", addr)
	}

	s := grpc.NewServer()
	api.RegisterSimCliServer(s, &CLIServer{})
	reflection.Register(s)

	// Control gRPC server lifecycle
	go func() {
		<-ctx.Done()
		infoLog.Print("Got notification to graceful stop CLI gRPC server.")
		s.Stop()
		return
	}()

	infoLog.Printf("%v started! Listen at %v\n", "CLI gRPC Server", addr)
	if err := s.Serve(lis); err != nil {
		return errors.Wrap(err, "Failed to start CLI gRPC server")
	}
	return nil
}
