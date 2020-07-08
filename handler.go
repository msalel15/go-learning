package main

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"github.com/google/uuid"
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

	collection := bucket.DefaultCollection()


	request := HelloRequest{}

	bindErr := ctx.Bind(&request)

	if bindErr != nil {
		return bindErr
	}

	response := HelloResponse{
		Message: fmt.Sprintf("hello %s, message: %s", request.Name, request.Message),
	}

	upsertResult, upsertErr := collection.Upsert(uuid.New().String(), response, &gocb.UpsertOptions{})

	if upsertErr != nil {
		return upsertErr
	}

	fmt.Println(upsertResult.Cas())

	err := ctx.JSON(200, response)

	if err != nil {
		return err
	}

	return nil
}
