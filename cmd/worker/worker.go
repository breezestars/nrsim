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

	workerServerPort int
	masterServerIp   string
)

func init() {
	infoLog = logger.InfoLog
	errLog = logger.ErrorLog
	debugLog = logger.DebugLog

	flag.IntVar(&workerServerPort, "workerSrvPort", 50052, "port for worker gRPC server")
	flag.StringVar(&masterServerIp, "masterSrvIp", "", "IP for master gRPC server")
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

	// Start gRPC server and prepare to be configure
	eg.Go(
		func() error {
			defer func() {
			}()
			if err := StartWorkerGrpcServer(ctx, workerServerPort); err != nil {
				return errors.Wrapf(err, "Start worker gRPC server failed")
			}
			return nil
		},
	)

	// Create master gRPC client
	client, connClose, err := createMasterGrpcClient()
	if err != nil {
		errLog.Printf(err.Error())
		return
	}
	defer connClose()

	// Register to master
	if err = register(ctx, client); err != nil {
		errLog.Printf(err.Error())
		return
	}

	if err := eg.Wait(); err != nil {
		errLog.Printf(err.Error())
		signalChannel <- syscall.SIGTERM
	}
}
