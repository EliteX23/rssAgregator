package logs

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"log"
	"os"
)


const (
	logFolder = "log"
	logFile = "log.txt"
)
var settings LogSettings

type Settings struct {
	Log LogSettings `json:"logSetting"`
}
type LogSettings struct {
	LogLevel string    `json:"logLevel"`
}
func InitConfig(filePath string) {
	cfg, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("File error: ", err.Error())
		os.Exit(1)
	}
	var settingIn Settings
	err = json.Unmarshal(cfg, &settingIn)
	if err != nil {
		log.Println("Произошла ошибка при попытке преобразовать конфиг файл из JSON: ", err.Error())
	} else {
		log.Println("Файл конфигурации успешно считан")
	}
	settings = settingIn.Log
}
func CreateLogger() *logrus.Logger {
	logger := logrus.New()
	if _, err := os.Stat(logFolder); os.IsNotExist(err) {
		os.Mkdir(logFolder, 0777)
	}
	level, err := logrus.ParseLevel(settings.LogLevel)
	if err!=nil{
		fmt.Println("ошибка при определении минимального уровня логирования:", err)
		level = logrus.DebugLevel
	}else{
		logger.SetLevel(level)	}

	writer := createWriter(logFolder+"/"+logFile, 0777)
	logger.Out = *writer

	return logger
}
func createWriter(fileName string, permissionCode os.FileMode) *io.Writer {
	var logWriters []io.Writer
	if len(fileName) == 0 {
		fmt.Println("[WARNING]. Log filename is empty! File binding is skipped! Use output to stdout only")
	} else {
		file, err := getFileWriter(fileName, permissionCode)
		if err != nil {
			fmt.Printf("[ERROR!] Cannot create/open file (filename:%v) for logger! Use output to stdout only... Reason: %v\n", fileName, err.Error())
		} else {
			logWriters = append(logWriters, file)
		}
	}
	logWriters = append(logWriters, os.Stdout)

	var writer = io.MultiWriter(logWriters...)
	return &writer
}

func getFileWriter(fileName string, permissionCode os.FileMode) (file *os.File, err error) {
	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, permissionCode)
	return
}