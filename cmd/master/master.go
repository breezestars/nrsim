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
	"flag"
	"github.com/cmingou/nrsim/internal/logger"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	infoLog  *log.Logger
	errLog   *log.Logger
	debugLog *log.Logger

	masterServerPort int
	cliServerPort    int
)

func init() {
	infoLog = logger.InfoLog
	errLog = logger.ErrorLog
	debugLog = logger.DebugLog

	flag.IntVar(&cliServerPort, "cliSrvPort", 50050, "port for CLI gRPC server")
	flag.IntVar(&masterServerPort, "masterSrvPort", 50051, "port for master gRPC server")
	flag.Parse()
}
func main() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		eg     errgroup.Group
	)

	ctx, cancel = context.WithCancel(context.Background())

	// Listen ctrl+c to terminate all gRPC server
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		cancel()
	}()

	// Start gRPC server and listen for register
	eg.Go(
		func() error {
			if err := StartMasterGrpcServer(ctx, masterServerPort); err != nil {
				return errors.Wrapf(err, "Start master gRPC server failed")
			}
			return nil
		},
	)

	// Start gRPC server to listen CLI
	eg.Go(
		func() error {
			if err := StartCLIGrpcServer(ctx, cliServerPort); err != nil {
				return errors.Wrapf(err, "Start CLI gRPC server failed")
			}
			return nil
		},
	)

	if err := eg.Wait(); err != nil {
		errLog.Printf(err.Error())
		signalChannel <- syscall.SIGTERM
	}
}
