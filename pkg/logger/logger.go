package logger

import "log"

func Debug(v ...any) {
	log.Println(append([]any{"DEBUG:"}, v...)...)
}

func Info(v ...any) {
	log.Println(append([]any{"INFO:"}, v...)...)
}

func Warn(v ...any) {
	log.Println(append([]any{"WARN:"}, v...)...)
}

func Error(v ...any) {
	log.Println(append([]any{"ERROR:"}, v...)...)
}

func Fatal(v ...any) {
	log.Fatalln(append([]any{"FATAL:"}, v...)...)
}
