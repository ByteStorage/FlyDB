package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path/filepath"
)

var logger *zap.Logger

const (
	EnvLogLevel      = "zap_level"
	LogLocation      = "./logs/runtime.log"
	ErrorLogLocation = "./logs/runtime_err.log"
	Release          = "release"
)

func init() {
	locations := []string{LogLocation, ErrorLogLocation}
	for _, location := range locations {
		filePathExist, err := pathExist(location)
		if err != nil {
			panic(err)
		}
		if !filePathExist {
			dir, _ := filepath.Split(location)
			err := os.MkdirAll(dir, os.ModeDir|os.ModePerm)
			if err != nil {
				panic(err)
			}
			_, err = os.Create(location)
			if err != nil {
				panic(err)
			}
		}
	}
	Init()
}

func Init() {
	encoder := getEncoder()

	// default DebugLevel
	level := zapcore.DebugLevel
	l := os.Getenv(EnvLogLevel)
	var c1, c2 zapcore.Core
	if l == Release {
		level = zapcore.InfoLevel
		// FILE ONLY
		c1 = zapcore.NewCore(encoder, getFileWriter(LogLocation), level)
		c2 = zapcore.NewCore(encoder, getFileWriter(ErrorLogLocation), zap.ErrorLevel)
	} else {
		// STD ONLY
		// c1 = zapcore.NewCore(encoder, getWriter(), level)
		c1 = zapcore.NewCore(encoder, getFileWriter(LogLocation), level)
		// STD and FILE
		c2 = zapcore.NewCore(encoder, getWriterWithFile(ErrorLogLocation), zap.ErrorLevel)
	}

	core := zapcore.NewTee(c1, c2)
	logger = zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(logger)
}

// getEncoder get encoder
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getLogWriterWithFile STD and FILE
func getWriterWithFile(location string) zapcore.WriteSyncer {
	file, err := os.OpenFile(location, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	mr := io.MultiWriter(file, os.Stdout)
	return zapcore.AddSync(mr)
}

// getLogWriter STD only
func getWriter() zapcore.WriteSyncer {
	writer := io.Writer(os.Stdout)
	return zapcore.AddSync(writer)
}

// getFileLogWriter FILE only
func getFileWriter(location string) zapcore.WriteSyncer {
	file, err := os.OpenFile(location, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	writer := io.Writer(file)

	return zapcore.AddSync(writer)
}

// PathExist check if the path exist
func pathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
