package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type appSettings struct {
	Debug     bool   `json:"debug"`
	Port      string `json:"port"`
	Secret    string `json:"secret"`
	MysqlHost string `json:"mysql_host"`
	MysqlDB   string `json:"mysql_db"`
	MysqlUser string `json:"mysql_user"`
	MysqlPass string `json:"mysql_pass"`
}

const filenameSettings = "settings.json"

var settings appSettings

func init() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dir, _ := filepath.Split(ex)
	dat, err := ioutil.ReadFile(dir + filenameSettings)
	if err != nil {
		data, _ := json.Marshal(settings)
		ioutil.WriteFile(dir+filenameSettings, data, 0664)
		log.Fatal("settings.json missing, " + err.Error())
	}

	if err := json.Unmarshal(dat, &settings); err != nil {
		log.Fatal(err)
	}
}
