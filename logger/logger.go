package logger

import (
	"log"
)

var (
    infLogger *log.Logger
    errLogger *log.Logger
    dbgLogger *log.Logger
)

func init() {
    // file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    // if err != nil {
    //     log.Fatalf("Failed to open log file: %v", err)
    // }

    // infLogger = log.New(file, "INF ", log.Ldate|log.Ltime)
    // errLogger = log.New(file, "ERR ", log.Ldate|log.Ltime)
    // dbgLogger = log.New(file, "DBG ", log.Ldate|log.Ltime)

}

func Info(v ...interface{}) {
    // infLogger.Println(v...)
    log.Println(v...)
}

func Error(v ...interface{}) {
    // errLogger.Println(v...)
    log.Println(v...)
}

func Debug(v ...interface{}) {
    // dbgLogger.Println(v...)
    log.Println(v...)
}
