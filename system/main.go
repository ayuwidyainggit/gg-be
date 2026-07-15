package main

import (
	"system/adapter"
	"system/controller"
	"system/pkg/config"
	"system/pkg/config/env"
	"system/pkg/validation"
	"system/repository"
	"system/service"

	"system/pkg/middleware"
	"system/pkg/server"

	"github.com/gofiber/fiber/v2"
)

func main() {
	var Redis config.RedisConfig
	// load config env
	envCfg := env.NewCfgEnv()
	validatorPkg := validation.NewValiditor()

	postgreDB := config.PostgreSQLConnection(envCfg)
	redisClient, errRedisCon := Redis.SetConfig().RedisInstance()
	if errRedisCon != nil {
		panic(errRedisCon)
	}

	// Setup Repository
	transactionDB := repository.NewDbtransactionRepo(postgreDB)
	smcMCustomerRepository := repository.NewSmcMCustomerRepository(postgreDB)
	userRepository := repository.NewUserRepository(postgreDB)
	mConfigRepository := repository.NewMConfigRepository(postgreDB)
	mMenuRepository := repository.NewMMenuRepository(postgreDB)

	mDayRepository := repository.NewMDayRepository(postgreDB)
	cacheDB := repository.NewCache(redisClient)

	// setup adapter
	obsAdapter, err := adapter.InitObsAdapter(envCfg.Get("OBS_HUAWEI_AK"), envCfg.Get("OBS_HUAWEI_SK"), envCfg.Get("OBS_HUAWEI_ENDPOINT"), envCfg.Get("OBS_HUAWEI_BUCKET"))
	if err != nil {
		panic(err)
	}

	httpClient := adapter.HttpClientInfo{}

	// Setup Service
	userService := service.NewUserService(envCfg, smcMCustomerRepository, userRepository, mMenuRepository, transactionDB, cacheDB)
	configService := service.NewMConfigService(mConfigRepository, transactionDB)
	mDayService := service.NewMDayService(mDayRepository)
	filesService := service.NewFilesService(envCfg, obsAdapter)
	notificationService := service.NewNotificationService(envCfg, httpClient)

	// Setup Controller
	userController := controller.NewUserController(userService, validatorPkg)
	configController := controller.NewMConfigController(configService, validatorPkg)
	mDayController := controller.NewMDayController(mDayService, validatorPkg)
	filesController := controller.NewFilesController(filesService, validatorPkg)
	notificationController := controller.NewNotificationController(notificationService, validatorPkg)

	// define Fiber Framework config
	fiberCfg := config.NewFiberConfig(envCfg)
	app := fiber.New(fiberCfg)

	// middleware
	middleware.AppMiddleware(app) // Register Fiber's middleware for app.

	// This route path for test service is running "/"
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("It works")
	})

	// route.PublicRoutes(app)
	userController.Route(app)
	configController.Route(app)
	mDayController.Route(app)
	filesController.Route(app)
	notificationController.Route(app)

	// itemController.Route(app)

	// start fiber server
	server.FiberServerWithGracefulShutdown(app)
}
