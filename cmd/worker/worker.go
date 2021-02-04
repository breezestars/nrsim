package main

import (
	"context"
	"flag"
	"github.com/ng5gc/uegnbsim/internal/logger"
	"log"
	"os"
	"os/signal"
	"sync"
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
		wg  sync.WaitGroup
		ctx context.Context
	)

	ctx, cancel := context.WithCancel(context.Background())

	// Listen ctrl+c to terminate all gRPC server
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		cancel()
	}()

	// Start gRPC server and prepare to be configure
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		if err := StartWorkerGrpcServer(ctx, workerServerPort); err != nil {
			errLog.Printf(err.Error())
			return
		}
	}()

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

	wg.Wait()
}
