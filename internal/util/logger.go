package util

import (
	"io"
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func Init(logToFile bool, filename string) error {
	var logOutput io.Writer = os.Stdout
	if logToFile && filename != "" {
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		logOutput = file
	}

	Info = log.New(logOutput, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(logOutput, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}
