package logging

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
)

var (
	InfSign  = color.CyanString("[INFO]")
	WarnSign = color.YellowString("[WARN]")
	ErrSign  = color.RedString("[ERROR]")
	SucSign  = color.GreenString("[SUCCESS]")

	day, month, year = time.Now().Date()

	gormLogFile *os.File
)

func WriteInfo(message interface{}) {
	log.Printf("%s %v\n", InfSign, message)
}

func WriteWarning(message interface{}) {
	log.Printf("%s %v\n", WarnSign, message)
}

func WriteError(message interface{}) {
	log.Printf("%s %v\n", ErrSign, message)
}

func WriteSuccess(message interface{}) {
	log.Printf("%s %v\n", SucSign, message)
}

func CreateDefaultLogsDirectory(logsDir string) error {
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		if err := os.Mkdir(logsDir, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func CreateGormLogger(logsPath string) (logger.Interface, error) {
	logFileName := fmt.Sprintf("%s/gorm_%d_%d_%d.log", logsPath, year, int(month), day)
	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	gormLogFile = f

	return logger.New(
		log.New(gormLogFile, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	), nil
}

func CloseLogFiles() error {
	return gormLogFile.Close()
}

func InitZapAPILogger(logsDir string, debug bool) *zap.Logger {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/api_zap.log", logsDir),
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     28,
	})

	var core zapcore.Core

	core = zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)

	if debug {
		core = zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), w, zap.DebugLevel),
			zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()), os.Stderr, zap.DebugLevel),
		)
	}

	logger := zap.New(core)

	return logger
}
