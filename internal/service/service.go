package service

import (
	"context"

	"shortener/internal/cache"
	"shortener/internal/model"
	"shortener/internal/repository"
)

type KeyService struct {
	repo  repository.ShortRepository
	cache cache.KeyCache
}

func NewKeyService(keyrep repository.ShortRepository, keycache cache.KeyCache) *KeyService {
	return &KeyService{repo: keyrep, cache: keycache}
}

func (k KeyService) CreateKey(ctx context.Context, link string) (*model.Link, error) {
	return nil, nil
}

/*
создание(внутри нужна проверка занятости)
запрос всех существующих ссылок
запрос одной ссылки с аналитикой
деактивация ссылки по истечении аренды
*/
