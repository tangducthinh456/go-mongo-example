package customlog

import (
	"io"
	"os"
	log "github.com/sirupsen/logrus"
)

// var (
// 	file *os.File
// )

func Init() {
	// var err error
	file, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	// defer file.Close()
	mw := io.MultiWriter(file, os.Stdout)
	log.SetOutput(mw)
}

func Err(err error)	{
	// log.SetLevel(logrus.ErrorLevel)
	log.Error(err)
}

func Info(info string) {
	log.Info(info)
}

func Warn(warn string){
	log.Warn(warn)
}