package main

import (
	"fmt"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

//auth middleware yapalim     //x
//rate limiting middleware   //bunu channel timer ile nasil yapabilir?
//db baglantisi ve repo
//unit test yazalim
//benchmark test yazalim
//dockerize edelim
//profilinge bakalim go da (memory leak yaratalim)
//deployment ve service yaml lari yazalim
//pipline yazalim ve gitlab uzerinden deploy edelim

func main() {

	server := echo.New()

	// Middleware
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	// Routes
	g := server.Group("/planet")

	g.GET("/:name", getPlanetHandler, RateLimitingMiddleware(2))
	g.POST("/:name", postPlanetHandler, BasicAuthMiddleware("test", "sifre123"))

	go func() {
		for {
			time.Sleep(10 * time.Millisecond)

			for s, limiting := range ratePaths {
				fmt.Println(limiting.time.Add(10 * time.Second).Before(time.Now()))
				if limiting.time.Add(10 * time.Second).Before(time.Now()) {
					delete(ratePaths, s)
				}
			}

		}
	}()

	// Start server
	server.Logger.Fatal(server.Start(":8081"))
}
