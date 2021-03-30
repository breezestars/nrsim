/*
Copyright © 2021 Jimmy Ou <breezestars@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"github.com/cmingou/nrsim/internal/api"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	NrMap sync.Map
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

	// TODO: Redesign
	NrMap.Range(
		func(key, value interface{}) bool {
			if value.(*GnbConfig).Ip == ip.String() {
				value.(*GnbConfig).Registered = true
				return false
			}
			return true
		},
	)

	// Delete first and re-store
	if _, found := s.workers.Load(ip.String()); found {
		s.workers.Delete(ip.String())
		infoLog.Printf("Duplicated worker: %v found, will delete and re-store", ip.String())
	}
	s.workers.Store(ip.String(), ip)

	infoLog.Printf("IP: %v has finished registration", ip.String())

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
			debugLog.Printf("Connection between with %v is normal", ip.String())
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
		infoLog.Print("Got notification to stop Master gRPC server.")
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

	containerClient *client.Client
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
		infoLog.Print("Got notification to stop CLI gRPC server.")
		s.Stop()
		return
	}()

	infoLog.Printf("%v started! Listen at %v\n", "CLI gRPC Server", addr)
	if err := s.Serve(lis); err != nil {
		return errors.Wrap(err, "Failed to start CLI gRPC server")
	}
	return nil
}

func (s *CLIServer) CreateGnb(ctx context.Context, cfg *api.GnbConfig) (*emptypb.Empty, error) {
	contName := s.GenContainerName(cfg.GlobalGNBID.Gnbid)

	// Check if existed
	if _, found := NrMap.Load(contName); found {
		return &emptypb.Empty{}, status.New(codes.AlreadyExists, "NR already existed.").Err()
	}

	// Create new worker
	contId, ip, err := s.NewWorker(contName)
	if err != nil {
		return &emptypb.Empty{}, errors.Wrapf(err, "New worker failed.")
	}

	// TODO: Should check the integrity of cfg
	NrMap.Store(contName, &GnbConfig{
		ContainerId: contId,
		Config:      cfg,
		Ip:          ip,
	})

	debugLog.Printf("Created Gnb which IP:%v", ip)

	return &emptypb.Empty{}, nil
}

func (s *CLIServer) DelGnb(ctx context.Context, id *api.IdMessage) (*emptypb.Empty, error) {
	contName := s.GenContainerName(id.Id)

	// If found then delete, if not found then return nil.
	// TODO: Change to LoadAndDelete func when the issue been fix,
	// https://github.com/golang/go/issues/40999
	if v, found := NrMap.Load(contName); found {
		// TODO: Del container
		if err := s.DelWorker(v.(*GnbConfig).ContainerId, contName); err != nil {
			return &emptypb.Empty{}, status.New(codes.Canceled, "Del worker failed.").Err()
		}

		NrMap.Delete(contName)
	}

	return &emptypb.Empty{}, nil
}

func (s *CLIServer) ListGnb(ctx context.Context, empty *emptypb.Empty) (*api.GnbConfigList, error) {
	var (
		cfgList = &api.GnbConfigList{
			GnbConfig: make([]*api.GnbConfig, 0),
		}
	)

	NrMap.Range(
		func(key, value interface{}) bool {
			cfgList.GnbConfig = append(cfgList.GnbConfig, value.(*GnbConfig).Config)
			return true
		},
	)

	return cfgList, nil
}
