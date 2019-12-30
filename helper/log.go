package helper

import (
	"log"
	"os"
	"time"
)

var LogBasePath = "."

func InitLog() {
	dir, _ := os.Getwd()
	LogBasePath = dir + "/../log"
}

func WriteLog(message string, tag string) {
	file := LogBasePath + "/" + time.Now().String()[0:10] + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if nil != err {
		panic(err)
	}
	loger := log.New(logFile, "", log.LstdFlags)
	loger.Println(tag + " | " + message)
}
