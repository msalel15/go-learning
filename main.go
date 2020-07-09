package main

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	planet "go-learning/domain/planet"
	"go-learning/domain/user"
	"os"
	"os/signal"
	"syscall"
)

//auth middleware yapalim    //x
//rate limiting middleware   //x
//db baglantisi ve repo      //x
//unit test yazalim          // todo
//benchmark test yazalim
//dockerize edelim
//profilinge bakalim go da (memory leak yaratalim)
//deployment ve service yaml lari yazalim
//pipline yazalim ve gitlab uzerinden deploy edelim

func main() {

	cluster := connectToCouchbase()
	server := startHttpServer()

	p := planet.New(server, cluster)
	u := user.New(server, cluster)

	pErr := p.Init()

	if pErr != nil {
		panic(pErr)
	}

	uErr := u.Init()

	if uErr != nil {
		panic(uErr)
	}

	go server.Logger.Fatal(server.Start(fmt.Sprintf(":%d", 8081)))

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func connectToCouchbase() *gocb.Cluster {
	cluster, err := gocb.Connect(
		"127.0.0.1",
		gocb.ClusterOptions{
			Username: "go-learning",
			Password: "123456",
		})

	if err != nil {
		panic(err)
	}

	return cluster
}

func startHttpServer() *echo.Echo {
	server := echo.New()

	// Middleware
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	return server
}
