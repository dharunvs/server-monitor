package logger

import (
    "log"
    "os"
)

var (
    infoLogger  *log.Logger
    errorLogger *log.Logger
)

func init() {
    file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err)
    }

    infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(v ...interface{}) {
    infoLogger.Println(v...)
}

func Error(v ...interface{}) {
    errorLogger.Println(v...)
}
