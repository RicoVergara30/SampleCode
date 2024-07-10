package logger

import (
	"log"
	"os"
	"time"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	Separator     *log.Logger
)

func OpenLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, err
}

func SetLogs() {
	logFileName := time.Now().Format("2006-01-02") + ".log"
	// create a new log with file name xxx or more the existing one
	f, err := os.OpenFile("./logs/insta-"+logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	// set default log output to the 'new' file
	log.SetOutput(f)
	log.Println("This is a test log entry")
	defer f.Close()
}

func Login_Page() {

}
