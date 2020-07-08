package main

import (
	"github.com/couchbase/gocb/v2"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	goLearningMiddleware "go-learning/middleware"
	"time"
)

//auth middleware yapalim    //x
//rate limiting middleware   //x
//db baglantisi ve repo      //rate limitingi db uzerinden yonetmeliyim
//unit test yazalim          // todo
//benchmark test yazalim
//dockerize edelim
//profilinge bakalim go da (memory leak yaratalim)
//deployment ve service yaml lari yazalim
//pipline yazalim ve gitlab uzerinden deploy edelim

var bucket *gocb.Bucket

func main() {

	server := echo.New()

	// Middleware
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	// Routes
	g := server.Group("/planet")

	rateLimiterConfig := goLearningMiddleware.RateLimiterConfig{
		Limit: 2,
		Identifier: func(context echo.Context) string {
			return context.Path()
		},
		BanDurationInSeconds: 5,
	}

	rateLimiter := goLearningMiddleware.NewRateLimiter(rateLimiterConfig)

	g.GET("/:name", getPlanetHandler, rateLimiter.RateLimiterMiddleware())
	g.POST("/:name", postPlanetHandler, goLearningMiddleware.BasicAuthMiddleware("test", "sifre123"))

	cluster, err := gocb.Connect(
		"127.0.0.1",
		gocb.ClusterOptions{
			Username: "go-learning",
			Password: "123456",
		})

	if err != nil {
		panic(err)
	}

	bucket = cluster.Bucket("go-learning")

	err = bucket.WaitUntilReady(15*time.Second, nil)
	if err != nil {
		panic(err)
	}

	// Start server
	server.Logger.Fatal(server.Start(":8081"))
}
