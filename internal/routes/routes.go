package routes

import (
	"errors"
	"fmt"
	"github.com/dw-parameter-service/configs"
	"github.com/dw-parameter-service/pkg/xlogger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"os"
)

func setupRoutes(app *fiber.App) {

	api := app.Group("/api/v1/params")
	initVoucherRoutes(api)

	xlogger.Log.Println("| routes >> initialized")
}

func getLogFile() *os.File {
	// Define file to logs
	f, err := os.OpenFile(configs.MainConfig.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	return f
}

func getFiberLogConfig(a *fiber.App) {
	var conf logger.Config

	if configs.MainConfig.DebugMode == true {
		conf.Format = "[REST-API] ${time} | ${status} | ${latency} | ${method} | ${path}\n"
		if configs.MainConfig.VerboseAPIResponse {
			conf.Format = "[REST-API] ${time} | ${status} | ${latency} | ${method} | ${path} \n\tResponse: ${resBody}\n"
		}

		conf.TimeFormat = "2006/01/02 15:04:05"

		if configs.MainConfig.LogOutput == "file" {
			conf.Output = getLogFile()
		}

		a.Use(logger.New(conf))
	}

}

func Initialize() error {

	app := fiber.New()
	getFiberLogConfig(app)
	setupRoutes(app)

	err := app.Listen(fmt.Sprintf(":%s", configs.MainConfig.APIServer.Port))
	if err != nil {
		return errors.New(fmt.Sprintf("error on starting service: %s", err.Error()))
	}

	return nil

}
