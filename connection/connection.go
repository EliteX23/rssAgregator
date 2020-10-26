package connection

import (
	"encoding/json"
	"fmt"
	"github.com/gocraft/dbr"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"io/ioutil"
	"log"
	"os"
)

var settings DataBase

type Settings struct {
	DataBase DataBase`json:"dataBase"`
}
type DataBase struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	User         string `json:"user"`
	DBName       string `json:"dbName"`
	Password     string `json:"password"`
	MaxOpenConns int    `json:"maxOpenConns"`
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
	settings = settingIn.DataBase
}

func InitDBRConnectionPG() *dbr.Connection {
	connectionString := "host=" + settings.Host + " port=" + settings.Port + " user=" +
		settings.User + " dbname=" + settings.DBName + " sslmode=disable password=" +
		settings.Password
	conn, err := dbr.Open("postgres", connectionString, nil)
	if err != nil {
		fmt.Printf("Error dbr: %v \n", err)
	}
	conn.SetMaxOpenConns(6)
	conn.SetMaxIdleConns(2)
	return conn
}
