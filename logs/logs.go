package logs

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	errLog  *log.Logger
	infoLog *log.Logger
)

func InitLog() error {
	errFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return err
	}
	errLog = log.New(io.MultiWriter(errFile, os.Stderr), "ERROR: ", log.LstdFlags)

	infoFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return err
	}
	infoLog = log.New(io.MultiWriter(infoFile, os.Stdout), "INFO: ", log.LstdFlags)

	return nil
}

func logWithCaller(logger *log.Logger, format string, v ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		file = filepath.Base(file) // 只保留文件名
		format = fmt.Sprintf("%s:%d: %s", file, line, format)
	}
	logger.Printf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	logWithCaller(errLog, format, v...)
}

func Error(v ...interface{}) {
	logWithCaller(errLog, "%v", v...)
}

func Infof(format string, v ...interface{}) {
	logWithCaller(infoLog, format, v...)
}

func Info(v ...interface{}) {
	logWithCaller(infoLog, "%v", v...)
}
