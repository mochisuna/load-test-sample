package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mochisuna/load-test-sample/application"
	"github.com/mochisuna/load-test-sample/config"
	"github.com/mochisuna/load-test-sample/handler"
	"github.com/mochisuna/load-test-sample/infrastructure"
	"github.com/mochisuna/load-test-sample/infrastructure/db"
)

func main() {
	// parse options
	path := flag.String("c", "_tools/local/config.toml", "config file")
	flag.Parse()

	// import config
	conf := &config.Config{}
	log.Println(*path)
	if err := config.New(conf, *path); err != nil {
		panic(err)
	}

	// init db connection
	// master
	dbmClient, err := db.NewMySQL(&conf.DBMaster)
	if err != nil {
		panic(err)
	}
	defer dbmClient.Close()
	// slave
	dbsClient, err := db.NewMySQL(&conf.DBSlave)
	if err != nil {
		panic(err)
	}
	defer dbsClient.Close()

	// initialize and injection relay
	// init repository
	userRepo := infrastructure.NewUserRepository(dbmClient, dbsClient)
	// init application service
	userService := application.NewUserService(userRepo)
	// inject all services
	services := &handler.Services{
		UserService: userService,
	}

	// Run App server
	server := handler.New(conf.Server.Port, services, conf.Server.RedirectURL)
	log.Println("Start server")
	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("Failed ListenAndServe. err: %v", err))
	}
}
