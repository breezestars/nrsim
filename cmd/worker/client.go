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
	"github.com/vishvananda/netlink"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

func createMasterGrpcClient() (api.SimMasterClient, func() error, error) {
	if masterServerIp == "" {
		return nil, nil, errors.Errorf("The parameter masterSrvIp is required")
	}

	serverIp, serverPort, err := net.SplitHostPort(masterServerIp)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "The format of masterSrvIp: %v is wrong: %v", masterServerIp, err)
	}

	port, err := strconv.Atoi(serverPort)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Parsed master server port failed")
	}

	hostAddr := net.TCPAddr{
		IP:   net.ParseIP(serverIp),
		Port: port,
	}

	infoLog.Printf("Connecting to master server: %v", hostAddr.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	conn, err := grpc.DialContext(ctx, hostAddr.String(), grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Failed to dial: %v", err)
	}
	return api.NewSimMasterClient(conn), conn.Close, nil
}

func register(ctx context.Context, client api.SimMasterClient) error {
	var (
		// The N1/N2 interface name
		ifName = "eth0"
		wg     sync.WaitGroup
		ctx2   context.Context
		cancel context.CancelFunc
	)

	ctx2, cancel = context.WithCancel(context.Background())

	link, err := netlink.LinkByName(ifName)
	if err != nil {
		return errors.Wrapf(err, "Get interface: %v failed", ifName)
	}

	ips, err := netlink.AddrList(link, netlink.FAMILY_V4)
	if err != nil {
		return errors.Wrapf(err, "Failed to get IPs in interface: %v with family type: %v", ifName, netlink.FAMILY_V4)
	}

	if len(ips) != 1 {
		return errors.Errorf("Interface: %v have more than 1 IP, size: %v", ifName, len(ips))
	}

	// Create metadata with IP and register
	md := metadata.Pairs("IP", ips[0].IP.String())
	stream, err := client.StreamChannel(metadata.NewOutgoingContext(context.Background(), md))
	if err != nil {
		return errors.Wrapf(err, "Get stream client failed")
	}

	// Control register cycle, trigger by ctx or connection error
	go func() {
		defer func() {
			if err := stream.CloseSend(); err != nil {
				infoLog.Printf("Close send failed, %v", err)
			}
		}()

		select {
		// Trigger due to ctx
		case <-ctx.Done():
			cancel()
		// Trigger due to connection error
		case <-ctx2.Done():
		}
	}()

	// Send heartbeat
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx2.Done():
				return
			default:
				err = stream.Send(&emptypb.Empty{})
				time.Sleep(time.Second * 2)
			}

		}
	}()

	// Monitor connection
	wg.Add(1)
	go func() {
		defer func() {
			cancel()
			wg.Done()
		}()

		for {
			select {
			case <-ctx2.Done():
				infoLog.Printf("Close connection monitor")
				return
			default:
				_, err = stream.Recv()
				if err != nil {
					if errors.Is(err, io.EOF) {
						// The connection is EOF
						infoLog.Printf("Connection between master node status: %v", err)
						return
					} else {
						errLog.Printf("Connection receive error: %v", err)
						return
					}
				}
				debugLog.Printf("Conn is normal")
			}
		}
	}()

	wg.Wait()
	return nil
}
