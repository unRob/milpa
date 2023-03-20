// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package logger

import "github.com/sirupsen/logrus"

var Main = logrus.WithField("milpa", "compa")

func Sub(name string) *logrus.Entry {
	return Main.WithField("milpa", name)
}

type Level int

const (
	LevelPanic Level = iota
	LevelFatal
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
	LevelTrace
)

func Configure(timestamps bool, colors bool, silent bool, level Level) {
	Main.Logger.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
		DisableTimestamp:       !timestamps,
		ForceColors:            colors,
	})

	if silent {
		Main.Logger.SetLevel(logrus.ErrorLevel)
	} else {
		Main.Logger.SetLevel(logrus.AllLevels[level])
	}
}

func Debug(args ...any) {
	Main.Debug(args...)
}

func Debugf(format string, args ...any) {
	Main.Debugf(format, args...)
}

func Info(args ...any) {
	Main.Info(args...)
}

func Infof(format string, args ...any) {
	Main.Infof(format, args...)
}

func Warn(args ...any) {
	Main.Warn(args...)
}

func Warnf(format string, args ...any) {
	Main.Warnf(format, args...)
}

func Error(args ...any) {
	Main.Error(args...)
}

func Errorf(format string, args ...any) {
	Main.Errorf(format, args...)
}

func Fatal(args ...any) {
	Main.Fatal(args...)
}

func Fatalf(format string, args ...any) {
	Main.Fatalf(format, args...)
}

func Trace(args ...any) {
	Main.Trace(args...)
}

func Tracef(format string, args ...any) {
	Main.Tracef(format, args...)
}
