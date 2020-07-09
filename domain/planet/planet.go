package planet

import (
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
	"go-learning/domain/user"
	goLearningMiddleware "go-learning/middleware"
	"time"
)

type Planet struct {
	server *echo.Echo
	db     *gocb.Cluster
}

func New(echo *echo.Echo, cb *gocb.Cluster) *Planet {
	bucket := cb.Bucket("user")

	err := bucket.WaitUntilReady(15*time.Second, nil)
	if err != nil {
		panic(err)
	}

	return &Planet{
		server: echo,
		db:     cb,
	}
}

func (p *Planet) Init() error {
	g := p.server.Group("/planet")

	rateLimiterConfig := goLearningMiddleware.RateLimiterConfig{
		Limit: 2,
		Identifier: func(context echo.Context) string {
			return context.Path()
		},
		BanDurationInSeconds: 5,
	}

	rateLimiter := goLearningMiddleware.NewRateLimiter(rateLimiterConfig)

	userService := user.NewService(p.db)

	g.GET("/:name", p.getPlanetHandler, goLearningMiddleware.BasicAuthMiddleware(userService), rateLimiter.RateLimiterMiddleware())
	g.POST("/:name", p.postPlanetHandler, goLearningMiddleware.BasicAuthMiddleware(userService))

	return nil
}
