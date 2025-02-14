package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/kardianos/service"
)

const (
	appVersionStr = "v1.0.1"
	nameOfService = "sqm-merge-pdf-service"
)

var db *sql.DB

var (
	routes = Routes{
		Route{
			"Index",
			"GET",
			"/",
			defaultHandler,
			false,
		},
		Route{
			"Index",
			"GET",
			"/mergepdf",
			defaultHandler,
			false,
		},
	}
	router *mux.Router
)

func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	svcConfig := &service.Config{
		Name:        nameOfService,
		DisplayName: nameOfService,
		Description: nameOfService,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open("mysql", settings.MysqlUser+":"+settings.MysqlPass+"@tcp("+settings.MysqlHost+":3306)/"+settings.MysqlDB+"?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<html><body>We are up and running "+nameOfService+" version "+appVersionStr+" ;)</body></html>")
}
