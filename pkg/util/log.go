package util

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var (
	DebugLog  *logger
	InfoLog   *logger
	ErrorLog  *logger
	debugmode = flag.Bool("dev", false, "是否启用debug日志")
)

type logger struct {
	sync.Mutex
	out        *log.Logger
	fileprefix string
	filename   string
	logfile    *os.File
	dev        bool
}

func (lg *logger) Println(data ...any) {
	lg.Lock()
	defer lg.Unlock()
	if lg.dev {
		return
	}
	lg.checkdate()
	lg.out.Println(data...)
}

func (lg *logger) Printf(str string, data ...any) {
	lg.Lock()
	defer lg.Unlock()
	if lg.dev {
		return
	}
	lg.checkdate()
	lg.out.Printf(str, data...)
}

func NewLogger(out io.Writer, prefix string, flag int, fileprefix string) *logger {
	return &logger{
		out:        log.New(out, prefix, flag),
		fileprefix: fileprefix,
	}
}

func init() {
	T := time.Now()
	os.MkdirAll("log", os.ModePerm)
	logfile := fmt.Sprintf("log/%04d-%02d-%02d.log", T.Year(), T.Month(), T.Day())
	errorfile := fmt.Sprintf("log/error%04d-%02d-%02d.log", T.Year(), T.Month(), T.Day())
	debugfile := fmt.Sprintf("log/debug%04d-%02d-%02d.log", T.Year(), T.Month(), T.Day())
	logf, _ := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	debugf, _ := os.OpenFile(debugfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	errorf, _ := os.OpenFile(errorfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	InfoLog = NewLogger(logf, "[Info]  - ", log.Ltime, "")
	// InfoLog.Mutex = sync.Mutex{}
	DebugLog = NewLogger(debugf, "[Debug] - ", log.Ltime, "debug")
	DebugLog.dev = !*debugmode
	// DebugLog.Mutex = sync.Mutex{}
	ErrorLog = NewLogger(errorf, "[Error] - ", log.Ltime, "error")
	// ErrorLog.Mutex = sync.Mutex{}
	InfoLog.logfile = logf
	DebugLog.logfile = debugf
	ErrorLog.logfile = errorf
	InfoLog.filename = logfile
	ErrorLog.filename = errorfile
	DebugLog.filename = debugfile
	ErrorLog.Println("Error Start logging", T.String())
	InfoLog.Println("Error Start logging", T.String())
	DebugLog.Println("Error Start logging", T.String())
}

func (lg *logger) checkdate() {
	//T := time.Now()
	//tmp := fmt.Sprintf("log/%s%04d-%02d-%02d.log", lg.fileprefix, T.Year(), T.Month(), T.Day())
	//if tmp != lg.filename {
	//	lg.filename = tmp
	//	lg.logfile.Close()
	//	lg.logfile, _ = os.OpenFile(lg.filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	//	lg.out.SetOutput(lg.logfile)
	//}
}
