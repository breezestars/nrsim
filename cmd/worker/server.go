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
