package main

import (
	"cronjob/controller"
	"cronjob/pkg/config"
	"cronjob/pkg/config/env"
	"cronjob/pkg/middleware"
	"cronjob/pkg/server"
	"cronjob/pkg/validation"
	"cronjob/repository"
	"cronjob/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// load config env
	envCfg := env.NewCfgEnv()
	validatorPkg := validation.NewValiditor()
	postgreDB := config.PostgreSQLConnection(envCfg)

	// Setup Repository
	transactionDB := repository.NewDbtransactionRepo(postgreDB)
	jobRepository := repository.NewJobRepository(postgreDB)

	// Setup Service
	jobService := service.NewJobService(envCfg, jobRepository, transactionDB)

	// Setup Controller
	jobController := controller.NewJobController(envCfg, jobService, validatorPkg)

	// define Fiber Framework config
	fiberCfg := config.NewFiberConfig(envCfg)
	app := fiber.New(fiberCfg)

	// middleware
	middleware.AppMiddleware(app) // Register Fiber's middleware for app.

	// This route path for test service is running "/"
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("It works")
	})

	jobController.Route(app)

	// start fiber server
	server.FiberServerWithGracefulShutdown(app)

}


