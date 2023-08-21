package util

import (
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func InitLoggers() {
	infoFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening info.log:", err)
	}
	errorFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening error.log:", err)
	}

	Info = log.New(infoFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
