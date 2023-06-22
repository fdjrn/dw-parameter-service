package main

import (
	"fmt"
	"github.com/dw-parameter-service/configs"
	"github.com/dw-parameter-service/internal"
	"github.com/dw-parameter-service/internal/db"
	"github.com/dw-parameter-service/internal/routes"
	"github.com/dw-parameter-service/pkg/xlogger"
)

func main() {
	var err error
	internal.SetupCloseHandler()

	defer internal.ExitGracefully()

	// Service Initialization
	err = configs.Initialize()
	if err != nil {
		xlogger.Log.Fatalln(fmt.Sprintf("error on config initialization: %s", err.Error()))
	}

	// DB Connection
	if err = db.Mongo.Connect(); err != nil {
		xlogger.Log.Fatalln(fmt.Sprintf("error on mongodb connection: %s", err.Error()))
	}

	// Start Rest API

	err = routes.Initialize()
	if err != nil {
		xlogger.Log.Fatalln(err)
	}

}
