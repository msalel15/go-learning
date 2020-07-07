package main

import (
	"encoding/base64"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"strings"
	"time"
)

func BasicAuthMiddleware(u, p string) echo.MiddlewareFunc {
	return func(handler echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			authorization := ctx.Request().Header.Get("Authorization")

			encodedAuthorization, err := base64.StdEncoding.DecodeString(authorization)

			if err != nil {
				return errors.Wrap(err, "Base64 encoding fail")
			}

			basicAuthValues := strings.Split(string(encodedAuthorization), ":")

			if len(basicAuthValues) != 2 {
				return errors.New("basic auth fail")
			}

			username := basicAuthValues[0]
			password := basicAuthValues[1]

			if !strings.EqualFold(username, u) || !strings.EqualFold(password, p) {

				e := GoLearningError{
					Message: "Unauthorized",
					Code:    40100,
				}

				jsonErr := ctx.JSON(401, e)

				if jsonErr != nil {
					return jsonErr
				}

				ctx.Logger().Error(e)

				return e
			}

			return nil
		}
	}
}

type rateLimiting struct {
	c     int
	open  bool
	time  time.Time
	timer time.Timer
}

//var ignoreUsers = []string{
//	"",
//}
var ratePaths = map[string]rateLimiting{}

func RateLimitingMiddleware(limit int) echo.MiddlewareFunc {
	return func(handler echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			if v, ok := ratePaths[ctx.Path()]; ok {
				if v.c >= limit {

					e := GoLearningError{
						Message: "Limit exceed",
						Code:    40300,
					}

					jsonErr := ctx.JSON(403, e)

					if jsonErr != nil {
						return jsonErr
					}

					ctx.Logger().Error(e)

					if !v.open {
						v.open = true
						v.time = time.Now()
					}

					ratePaths[ctx.Path()] = v

					return e
				}
				v.c++

				ratePaths[ctx.Path()] = v
			} else {
				ratePaths[ctx.Path()] = rateLimiting{
					c: 0,
				}
			}

			return nil
		}
	}
}
