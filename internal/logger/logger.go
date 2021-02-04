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
