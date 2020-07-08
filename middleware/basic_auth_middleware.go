package middleware

import (
	"encoding/base64"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	goLearningErrors "go-learning/errors"
	"strings"
)

func BasicAuthMiddleware(u, p string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
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

				e := goLearningErrors.GoLearningError{
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

			return next(ctx)
		}
	}
}