package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/UnendingLoop/link_shortener/internal/cache"
	"github.com/UnendingLoop/link_shortener/internal/model"
	"github.com/UnendingLoop/link_shortener/internal/repository"
)

var (
	ErrNoRedirLink error = errors.New("empty link is provided for redirect")                       // 422
	ErrKeyGen      error = errors.New("failed to generate valid shortkey. Try again later")        // 500
	ErrKeyCheck    error = errors.New("failed to check shortkey validity. Try again later")        // 500
	ErrBusyKey     error = errors.New("provided shortkey is not available. Try using another one") // 409
	ErrCommon500   error = errors.New("something went wrong. Try again later")                     // 500
	ErrAnalytics   error = errors.New("incorreect query parameters for analytics request")         // 422
	ErrKeyNotFound error = errors.New("specified shortkey doesn't exist")                          // 404

)

type ShortService interface {
	CreateKey(ctx context.Context, link *model.Link) (*model.Link, error)
	GetRedirLinkByKey(ctx context.Context, key, ua string) (string, error)
	GetAnalytics(ctx context.Context, req *model.AnalyticsRequest, limit, offset int) (*model.AnalyticsResponse, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.Link, error)
}

type SKService struct {
	repo  repository.ShortRepository
	cache cache.ShortCache
}

func NewKeyService(keyrep repository.ShortRepository, keycache cache.ShortCache) ShortService {
	return &SKService{repo: keyrep, cache: keycache}
}

func (s SKService) CreateKey(ctx context.Context, link *model.Link) (*model.Link, error) {
	if link.Redirect == "" {
		return nil, ErrNoRedirLink
	}

	// если юзер не предоставил свой вариант шортки, то генерируем сами
	custom := true
	for link.ShortKey == "" {
		var err error
		link.ShortKey, err = generateShortKey(8)
		if err != nil {
			return nil, ErrKeyGen
		}
		custom = false
	}

	// проверка уникальности в БД
	free, err := s.repo.CheckKeyIsFree(ctx, link.ShortKey)
	if err != nil {
		log.Printf("Failed to check uniqueness of a shortkey %q: %v", link.ShortKey, err)
		return nil, ErrKeyCheck
	}
	if !free {
		switch custom {
		case true:
			return nil, ErrBusyKey
		case false:
			return nil, ErrKeyGen
		}
	}

	// создаем новую запись в БД
	if err := s.repo.Create(ctx, link); err != nil {
		log.Printf("Failed to insert shortkey %q into DB: %v \n", link.ShortKey, err)
		return nil, ErrCommon500
	}

	// кладем в кеш
	if err := s.cache.SetByShortkey(ctx, link.ShortKey, link.Redirect); err != nil {
		log.Printf("Failed to cache shortkey %q: %v\n", link.ShortKey, err)
	}

	return link, nil
}

func (s SKService) GetAll(ctx context.Context, limit, offset int) ([]model.Link, error) {
	links, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		log.Println("Failed to get []links from DB:", err)
		return nil, ErrCommon500
	}
	return links, nil
}

func (s SKService) GetAnalytics(ctx context.Context, req *model.AnalyticsRequest, limit, offset int) (*model.AnalyticsResponse, error) {
	if req.GroupBy == "" || req.Shortkey == "" {
		return nil, fmt.Errorf("%w: grouping or shortkey is empty", ErrAnalytics)
	}
	if req.Start != nil && req.End != nil && req.Start.After(*req.End) {
		return nil, fmt.Errorf("%w: start-date must be before end-date", ErrAnalytics)
	}

	// узнаем lid для запроса
	lid, err := s.repo.GetLIDByKey(ctx, req.Shortkey)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrKeyNotFound
	}
	if err != nil {
		log.Println("Failed to get LID by Key for analytics:", err)
		return nil, ErrCommon500
	}

	// выбираем нужный тип аналитики
	switch req.GroupBy {
	case model.GroupByDay:
		data, err := s.repo.GetGroupByDay(ctx, lid, req.Start, req.End, limit, offset)
		if err != nil {
			log.Printf("Failed to get daily analytics: %v", err)
			return nil, ErrCommon500
		}
		return data, nil
	case model.GroupByMonth:
		data, err := s.repo.GetGroupByMonth(ctx, lid, req.Start, req.End, limit, offset)
		if err != nil {
			log.Printf("Failed to get monthly analytics: %v", err)
			return nil, ErrCommon500
		}
		return data, nil
	case model.GroupByUserAgent:
		data, err := s.repo.GetGroupByUserAgent(ctx, lid, req.Start, req.End, limit, offset)
		if err != nil {
			log.Printf("Failed to get analytics by user-agent: %v", err)
			return nil, ErrCommon500
		}
		return data, nil
	default:
		return nil, fmt.Errorf("incorrect grouping specified: %q", req.GroupBy)
	}
}

func (s SKService) GetRedirLinkByKey(ctx context.Context, key string, useragent string) (string, error) {
	// проверяем кеш
	link, err := s.cache.GetByShortkey(ctx, key)
	if err == nil {
		// логируем вызов
		if err := s.repo.AddReferralByKey(ctx, key, truncUA(useragent)); err != nil {
			log.Printf("Failed to add referral to shortkey %q: %v\n", key, err)
		}
		return link, nil
	} else {
		if !errors.Is(err, repository.ErrNotFound) {
			log.Println("Failed to read data from cache:", err)
		}
	}

	// идем в базу
	link, err = s.repo.GetRedirLinkByKey(ctx, key)
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
	if err := s.cache.SetByShortkey(ctx, key, link); err != nil {
		log.Printf("Failed to cache shortkey %q: %v\n", key, err)
	}

	// логируем вызов
	if err := s.repo.AddReferralByKey(ctx, key, truncUA(useragent)); err != nil {
		log.Printf("Failed to add referral to shortkey %q: %v\n", key, err)
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

func truncUA(ua string) string {
	browser := detectBrowser(ua)
	platform := detectPlatform(ua)
	return fmt.Sprintf("%s/%s", browser, platform)
}

func detectBrowser(ua string) string {
	ua = strings.ToLower(ua)

	switch {
	case strings.Contains(ua, "edg/"):
		return "Edge"
	case strings.Contains(ua, "opr/") || strings.Contains(ua, "opera"):
		return "Opera"
	case strings.Contains(ua, "chrome/"):
		return "Chrome"
	case strings.Contains(ua, "safari/"):
		return "Safari"
	case strings.Contains(ua, "firefox/"):
		return "Firefox"
	default:
		return "Unknown"
	}
}

func detectPlatform(ua string) string {
	ua = strings.ToLower(ua)

	switch {
	case strings.Contains(ua, "windows nt"):
		return "Windows"
	case strings.Contains(ua, "android"):
		return "Android"
	case strings.Contains(ua, "iphone"), strings.Contains(ua, "ipad"):
		return "iOS"
	case strings.Contains(ua, "mac os x"):
		return "macOS"
	case strings.Contains(ua, "linux"):
		return "Linux"
	default:
		return "Unknown"
	}
}
