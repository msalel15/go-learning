package user

import (
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
)

type User struct {
	server  *echo.Echo
	service UserService
}

func New(echo *echo.Echo, cb *gocb.Cluster) *User {
	return &User{
		server:  echo,
		service: NewService(cb),
	}
}

func (u *User) Init() error {

	g := u.server.Group("/user")

	g.POST("", u.signUpHandler)

	return nil
}
