package service

import (
	"shortener/internal/cache"
	"shortener/internal/repository"
)

type Service struct {
	repo  repository.ShortRepository
	cache cache.KeyCache
}
