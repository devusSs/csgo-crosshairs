package logging

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"gorm.io/gorm/logger"
)

const (
	defaultLogPath = "./logs"
)

var (
	InfSign  = color.WhiteString("i")
	WarnSign = color.YellowString("!")
	ErrSign  = color.RedString("x")
	SucSign  = color.GreenString("✓")

	day, month, year = time.Now().Date()

	appLogFile   *os.File
	errorLogFile *os.File
	gormLogFile  *os.File
)

func WriteInfo(message interface{}) {
	log.Printf("[%s] %v\n", InfSign, message)

	_, err := appLogFile.WriteString(fmt.Sprintf("%v\n", message))
	if err != nil {
		log.Printf("[%s] Error writing to log file: %s\n", ErrSign, err.Error())
	}
}

func WriteWarning(message interface{}) {
	log.Printf("[%s] %v\n", WarnSign, message)

	_, err := appLogFile.WriteString(fmt.Sprintf("%v\n", message))
	if err != nil {
		log.Printf("[%s] Error writing to log file: %s\n", ErrSign, err.Error())
	}
}

func WriteError(message interface{}) {
	log.Printf("[%s] %v\n", ErrSign, message)

	_, err := errorLogFile.WriteString(fmt.Sprintf("%v\n", message))
	if err != nil {
		log.Printf("[%s] Error writing to log file: %s\n", ErrSign, err.Error())
	}
}

func WriteSuccess(message interface{}) {
	log.Printf("[%s] %v\n", SucSign, message)

	_, err := appLogFile.WriteString(fmt.Sprintf("%v\n", message))
	if err != nil {
		log.Printf("[%s] Error writing to log file: %s\n", ErrSign, err.Error())
	}
}

func CreateAppLogFile() error {
	logFileName := fmt.Sprintf("%s/app_%d_%d_%d.log", defaultLogPath, year, int(month), day)
	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	appLogFile = f
	return nil
}

func CreateErrorLogFile() error {
	logFileName := fmt.Sprintf("%s/error_%d_%d_%d.log", defaultLogPath, year, int(month), day)
	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	errorLogFile = f
	return nil
}

func CreateDefaultLogsDirectory() error {
	if _, err := os.Stat(defaultLogPath); os.IsNotExist(err) {
		if err := os.Mkdir(defaultLogPath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func CloseLogFiles() error {
	if err := appLogFile.Close(); err != nil {
		return err
	}
	if err := errorLogFile.Close(); err != nil {
		return err
	}
	return gormLogFile.Close()
}

func CreateGormLogger() (logger.Interface, error) {
	logFileName := fmt.Sprintf("%s/gorm_%d_%d_%d.log", defaultLogPath, year, int(month), day)
	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	gormLogFile = f

	return logger.New(
		log.New(gormLogFile, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	), nil
}

func CreateAPILogFiles() (*os.File, *os.File, error) {
	logFileName := fmt.Sprintf("%s/api_%d_%d_%d.log", defaultLogPath, year, int(month), day)
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}

	logFileName = fmt.Sprintf("%s/api_error_%d_%d_%d.log", defaultLogPath, year, int(month), day)
	errorLogFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}

	return logFile, errorLogFile, nil
}
