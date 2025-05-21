package logger

import "log"

func Debug(msg string) {
	log.Println("DEBUG:", msg)
}

func Info(msg string) {
	log.Println("INFO:", msg)
}

func Warn(msg string) {
	log.Println("WARN:", msg)
}

func Error(err error) {
	log.Println("ERROR:", err)
}

func Fatal(err error) {
	log.Fatalln("FATAL:", err)
}
