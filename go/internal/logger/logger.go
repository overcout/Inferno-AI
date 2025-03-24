package logger

import (
	"log"
	"os"
)

var (
	Info    *log.Logger
	Error   *log.Logger
	Warning *log.Logger
)

func Init() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Can't open log file:", err)
	}

	Info = log.New(logFile, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(logFile, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(logFile, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func InitConsole() {
	Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}
