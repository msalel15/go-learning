package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

func getPlanetHandler(ctx echo.Context) error {

	response := HelloResponse{
		Message: fmt.Sprintf("hello world"),
	}

	err := ctx.JSON(200, response)

	if err != nil {
		return err
	}

	return nil
}

func postPlanetHandler(ctx echo.Context) error {

	request := HelloRequest{}

	bindErr := ctx.Bind(&request)

	if bindErr != nil {
		return bindErr
	}

	response := HelloResponse{
		Message: fmt.Sprintf("hello %s, message: %s", request.Name, request.Message),
	}

	err := ctx.JSON(200, response)

	if err != nil {
		return err
	}

	return nil
}
