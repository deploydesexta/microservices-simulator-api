package log

import (
	std "log"
)

type NoopLogger struct{}

func (*NoopLogger) Fatal(v ...interface{})                 { std.Fatal(v...) }
func (*NoopLogger) Fatalf(format string, v ...interface{}) { std.Fatalf(format, v...) }
func (*NoopLogger) Print(...interface{})                   {}
func (*NoopLogger) Println(...interface{})                 {}
func (*NoopLogger) Printf(string, ...interface{})          {}
