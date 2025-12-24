package cmd

import (
	"log"
	"time"

	"shortener/internal/api"
	"shortener/internal/cache"
	"shortener/internal/repository"
	"shortener/internal/service"

	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/redis"
)

func StartApp() {
	// init конфига
	appConfig := config.New()
	appConfig.EnableEnv("")
	if err := appConfig.LoadEnvFiles("./.env"); err != nil {
		log.Fatalf("Failed to load envs: %s\nExiting app...", err)
	}

	time.Sleep(25 * time.Second)

	// подключение к редису
	redisAddr := appConfig.GetString("REDIS_ADDR")
	redisPwd := appConfig.GetString("REDIS_PASSWORD")
	rawRedis := redis.New(redisAddr, redisPwd, 0)
	defer rawRedis.Close()
	redisClient := cache.NewRedisCache(rawRedis, 1*time.Hour)

	// подключение к БД
	dbConn := repository.ConnectWithRetries(appConfig, 5, 10*time.Second)
	defer func() {
		if err := dbConn.Master.Close(); err != nil {
			log.Println("Failed to close DB-conn correctly:", err)
		}
	}()

	// Накатываем миграцию
	if err := repository.Migrate(dbConn.Master, "./migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %s", err)
	}

	// Создание экземпляра репозитория
	repo := repository.NewPostgresRepo(dbConn)

	// Создание экземпляра сервиса
	svc := service.NewKeyService(repo, redisClient)

	// инит сервера
	server := ginext.New("") // empty - debug mode, release - prod mode
	handlers := api.NewHandler(svc)

	server.GET("/ping", handlers.SimplePinger)
	server.POST("/shorten", handlers.CreateKey)
	server.GET("/s/:short_url", handlers.Redirect)
	server.GET("/shorten/all", handlers.GetAll)
	server.GET("/analytics/:short_url", handlers.GetAnalytics)
	server.Static("/web", "./internal/web")

	// Server launch
	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
