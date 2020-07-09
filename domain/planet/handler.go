package planet

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

func (p *Planet) getPlanetHandler(ctx echo.Context) error {

	response := PlanetDTO{
		Message: fmt.Sprintf("hello world"),
	}

	err := ctx.JSON(200, response)

	if err != nil {
		return err
	}

	return nil
}

func (p *Planet) postPlanetHandler(ctx echo.Context) error {

	request := HelloRequest{}

	bindErr := ctx.Bind(&request)

	if bindErr != nil {
		return bindErr
	}

	response := PlanetDTO{
		Message: fmt.Sprintf("hello %s, message: %s", request.Name, request.Message),
	}

	err := ctx.JSON(200, response)

	if err != nil {
		return err
	}

	return nil
}
