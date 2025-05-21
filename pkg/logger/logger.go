package logger

import "log"

func Debug(v ...any) {
	log.Println("DEBUG:", v)
}

func Info(v ...any) {
	log.Println("INFO:", v)
}

func Warn(v ...any) {
	log.Println("WARN:", v)
}

func Error(v ...any) {
	log.Println("ERROR:", v)
}

func Fatal(v ...any) {
	log.Fatalln("FATAL:", v)
}
