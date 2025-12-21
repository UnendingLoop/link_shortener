package main

import (
	"log"
	"time"

	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/redis"
)

func main() {
	//init конфига
	appConfig := config.New()
	appConfig.EnableEnv("")
	if err := appConfig.LoadEnvFiles("./.env"); err != nil {
		log.Fatalf("Failed to load envs: %s\nExiting app...", err)
	}

	//подключение к редису
	redisAddr := appConfig.GetString("REDIS_ADDR")
	redisPwd := appConfig.GetString("REDIS_PASSWORD")
	rawRedis := redis.New(redisAddr, redisPwd, 0)
	defer rawRedis.Close()
	redisClient := cache.NewRedisCache(rawRedis, 1*time.Hour)

	//подключение к постгрес

	//создание экземпляра репозитория

	//создание экземпляра сервиса

	//инит сервера
	server := ginext.New("") //empty - debug mode, release - prod mode
	handlers := handler.NewHandler(svc)
	api := server.Group("/api")
	notify := api.Group("/notify")

	server.GET("/ping", handlers.SimplePinger)
	notify.POST("", handlers.CreateTask)
	notify.GET("/:uid", handlers.GetTask)
	notify.GET("/all", handlers.GetAll)
	notify.DELETE("/:uid", handlers.DeleteTask)
	server.Static("/web", "./internal/web")

	// Server launch
	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}

/*
Redis: key Shortlink, value OuterLink
Analytics: no cache, фильтр по юзерагенту и времени + их комбинация
DB: 2 related tables with cascade deletion
UI: 3 блока: форма для создания, список всех сокращенных ссылок в базе, просмотр аналитики по 1 короткой ссылке с фильтрацией по юзерагенту и времени(возможность задавать период)
*/
