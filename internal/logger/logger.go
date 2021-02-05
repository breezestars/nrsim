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

package logger

import (
	"log"
	"os"
)

var (
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	DebugLog *log.Logger
)

func init() {
	InfoLog = log.New(os.Stdout, "[Info] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	ErrorLog = log.New(os.Stdout, "[Error] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	DebugLog = log.New(os.Stdout, "[Debug] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
}
