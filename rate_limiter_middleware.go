package main

import (
	"errors"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	RateLimitIdentifier func(echo.Context) string

	RateLimiterConfig struct {
		Limit                int
		Identifier           RateLimitIdentifier
		BanDurationInSeconds time.Duration
	}

	RateLimiter struct {
		config      RateLimiterConfig
		records     map[string]int
		bannedUsers map[string]banRecord
		banMutex    *sync.Mutex
	}

	banRecord struct {
		identifier  string
		releaseTime time.Time
	}
)

func (rateLimiter *RateLimiter) isUserBanned(identifier string) bool {
	rateLimiter.banMutex.Lock()
	defer rateLimiter.banMutex.Unlock()

	_, ok := rateLimiter.bannedUsers[identifier]
	return ok
}

func (rateLimiter *RateLimiter) banUser(identifier string) {
	banRecord := banRecord{
		identifier:  identifier,
		releaseTime: time.Now().Add(rateLimiter.config.BanDurationInSeconds * time.Second),
	}

	rateLimiter.banMutex.Lock()
	defer rateLimiter.banMutex.Unlock()

	rateLimiter.bannedUsers[identifier] = banRecord
}

func NewRateLimiter(rlConfig RateLimiterConfig) *RateLimiter {
	rateLimiter := &RateLimiter{
		config: rlConfig,
		// TODO Access to records field must also be guarded with locks.
		records:     make(map[string]int),
		bannedUsers: make(map[string]banRecord),
		banMutex:    &sync.Mutex{},
	}

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {

			rateLimiter.banMutex.Lock()

			timeNow := time.Now()
			for id, banInfo := range rateLimiter.bannedUsers {
				if banInfo.releaseTime.Before(timeNow) {
					delete(rateLimiter.bannedUsers, id)
					delete(rateLimiter.records, id)
				}
			}

			rateLimiter.banMutex.Unlock()
		}
	}()

	return rateLimiter
}

func RateLimiterMiddleware(rateLimiter *RateLimiter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// Get the identifier for the user.
			identifier := rateLimiter.config.Identifier(ctx)

			// Check to see if the user is already in ban list.
			// Else it is safe to assume that the user is not banned.
			if banned := rateLimiter.isUserBanned(identifier); banned {
				return errors.New("TODO")
			}

			// Check to see if the user has accessed the path before.
			count, ok := rateLimiter.records[identifier]
			// If not, this is their first request, record him.
			if !ok {
				rateLimiter.records[identifier] = 1
				return next(ctx)
			}

			// User has accessed the path before. Increase the count.
			// If the count reaches the limit, ban the user.
			count++
			rateLimiter.records[identifier] = count

			if count == rateLimiter.config.Limit {
				rateLimiter.banUser(identifier)
			}

			return next(ctx)
		}
	}
}
