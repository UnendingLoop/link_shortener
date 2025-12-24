package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"time"

	"shortener/internal/cache"
	"shortener/internal/model"
	"shortener/internal/repository"
)

var (
	ErrNoRedirLink  error = errors.New("empty link is provided for redirect")                      // 422
	ErrKeyGen       error = errors.New("failed to generate valid shortkey. Try again later")       // 500
	ErrKeyCheck     error = errors.New("failed to check shortkey validity. Try again later")       // 500
	ErrBusyKey      error = errors.New("provided shortkey is not avalable. Try using another one") // 409
	ErrCommon500    error = errors.New("something went wrong. Try again later")                    // 500
	ErrAnalytics422 error = errors.New("incorreect query parameters for analytics request")        // 422
	ErrKeyNotFound  error = errors.New("specified shortkey doesn't exist")                         // 404

)

type ShortService struct {
	repo  repository.ShortRepository
	cache cache.ShortCache
}

func NewKeyService(keyrep repository.ShortRepository, keycache cache.ShortCache) *ShortService {
	return &ShortService{repo: keyrep, cache: keycache}
}

func (k ShortService) CreateKey(ctx context.Context, link *model.Link) (*model.Link, error) {
	if link.Redirect == "" {
		return nil, ErrNoRedirLink
	}

	// если юзер не предоставил свой вариант шортки, то генерируем сами
	custom := false
	for link.ShortKey == "" {
		link.ShortKey, _ = generateShortKey(8)
		custom = true
	}

	// проверка уникальности в БД
	idx, err := k.repo.GetLIDByKey(ctx, link.ShortKey)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		log.Println("Failed to check shortkey uniqueness:", err)
		return nil, ErrKeyCheck
	}
	if idx != -1 {
		switch custom {
		case true:
			return nil, ErrBusyKey
		case false:
			return nil, ErrKeyGen
		}
	}

	// создаем новую запись в БД
	if err := k.repo.Create(ctx, link); err != nil {
		log.Printf("Failed to insert shortkey %q into DB: %v \n", link.ShortKey, err)
		return nil, ErrCommon500
	}

	// кладем в кеш
	if err := k.cache.SetByShortkey(ctx, link.ShortKey, link.Redirect); err != nil {
		log.Printf("Failed to cache shortkey %q: %v\n", link.ShortKey, err)
	}

	return link, nil
}

func (k ShortService) GetAll(ctx context.Context, limit, offset int) ([]model.Link, error) {
	links, err := k.repo.GetAll(ctx, limit, offset)
	if err != nil {
		log.Println("Failed to get []links from DB:", err)
		return nil, ErrCommon500
	}
	return links, nil
}

func (k ShortService) GetAnalytics(ctx context.Context, grouping string, key string, start, end *time.Time, limit, offset int) (*model.AnalyticsResponse, error) {
	if grouping == "" || key == "" || start.After(*end) || limit <= 0 || offset < 0 {
		return nil, ErrAnalytics422
	}
	// узнаем lid для запроса
	lid, err := k.repo.GetLIDByKey(ctx, key)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrKeyNotFound
	}
	if err != nil {
		log.Println("Failed to get LID by Key for analytics:", err)
		return nil, ErrCommon500
	}

	// выбираем нужный тип аналитики
	switch grouping {
	case "day":
		data, err := k.repo.GetGroupByDay(ctx, lid, start, end, limit, offset)
		if err != nil {
			log.Printf("Failed to get daily analytics: %v", err)
			return nil, ErrCommon500
		}
		return data, nil
	case "month":
		data, err := k.repo.GetGroupByMonth(ctx, lid, start, end, limit, offset)
		if err != nil {
			log.Printf("Failed to get monthly analytics: %v", err)
			return nil, ErrCommon500
		}
		return data, nil
	case "useragent":
		data, err := k.repo.GetGroupByUserAgent(ctx, lid, start, end, limit, offset)
		if err != nil {
			log.Printf("Failed to get analytics by user-agent: %v", err)
			return nil, ErrCommon500
		}
		return data, nil
	default:
		return nil, fmt.Errorf("incorrect grouping specified: %q", grouping)
	}
}

func (k ShortService) GetRedirLinkByKey(ctx context.Context, key string) (string, error) {
	// проверяем кеш
	link, err := k.cache.GetByShortkey(ctx, key)
	if err == nil {
		return link, nil
	} else {
		if !errors.Is(err, repository.ErrNotFound) {
			log.Println("Failed to read data from cache:", err)
		}
	}

	// идем в базу
	link, err = k.GetRedirLinkByKey(ctx, key)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return "", ErrKeyNotFound
		default:
			log.Println("Failed to get redirect link by shortkey from DB:", err)
			return "", ErrCommon500
		}
	}

	// кладем в кеш
	if err := k.cache.SetByShortkey(ctx, key, link); err != nil {
		log.Printf("Failed to cache shortkey %q: %v\n", key, err)
	}

	return link, nil
}

func generateShortKey(length int) (string, error) {
	b := make([]byte, length)

	// криптостойкая генерация случайных байтов
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		b[i] = model.Alphabet[int(b[i])%len(model.Alphabet)]
	}

	return string(b), nil
}
