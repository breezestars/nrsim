/*
 * Copyright © 2021 Jimmy Ou <breezestars@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"context"
	"github.com/cmingou/nrsim/internal/api"
	"google.golang.org/grpc"
	"net"
)

func GetCliServerClient() api.SimCliClient {
	if cliClient == nil {
		newCliServerClient()
	}
	return cliClient
}

func newCliServerClient() {
	serverAddr := net.TCPAddr{
		IP:   net.ParseIP(cliServerIp),
		Port: cliServerPort,
	}

	infoLog.Printf("Connecting to CLI server: %v", serverAddr.String())

	ctxCliSrv, ctxCliSrvCancel = context.WithTimeout(context.Background(), GrpcConnectTimeout)

	conn, err := grpc.DialContext(ctxCliSrv, serverAddr.String(), grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		errLog.Printf("Failed to dial: %v", err)
	}

	//TODO: deal with conn.Close when close connection

	cliClient = api.NewSimCliClient(conn)
}
