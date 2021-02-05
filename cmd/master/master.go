package main

import (
	"context"
	"flag"
	"github.com/ng5gc/uegnbsim/internal/logger"
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
