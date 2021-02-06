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
	"github.com/cmingou/nrsim/internal/logger"
	"log"
	"time"
)

const (
	GrpcConnectTimeout = time.Second * 3
)

var (
	// Logger
	infoLog  *log.Logger
	errLog   *log.Logger
	debugLog *log.Logger

	// CLI server information
	cliServerIp   = "127.0.0.1"
	cliServerPort = 50050

	// CLI server client and ctx
	cliClient       api.SimCliClient
	ctxCliSrvCancel context.CancelFunc
	ctxCliSrv       context.Context
)

func init() {
	infoLog = logger.InfoLog
	errLog = logger.ErrorLog
	debugLog = logger.DebugLog
}

func dealError(err error) {
	errLog.Printf("%v", err)
}
