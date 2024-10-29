package server

import (
	"log"
	"os"
)

var errorLog *log.Logger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
var warningLog *log.Logger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime)
var infoLog *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

func LogInfo(message ...any) {
	infoLog.Println(message)
}

func LogWarning(message ...any) {
	warningLog.Println(message)
}

func LogError(message ...any) {
	errorLog.Println(message)
}
