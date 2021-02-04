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

	// Start gRPC server and listen for register
	wg.Add(1)
	go func() {
		if err := StartMasterGrpcServer(ctx, masterServerPort); err != nil {
			errLog.Printf(err.Error())
		}
		wg.Done()
	}()

	// Start gRPC server to listen CLI
	wg.Add(1)
	go func() {
		if err := StartCLIGrpcServer(ctx, cliServerPort); err != nil {
			errLog.Printf(err.Error())
		}
		wg.Done()
	}()

	wg.Wait()
}
