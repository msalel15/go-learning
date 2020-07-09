package user

import (
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"time"
)

func (u *User) signUpHandler(ctx echo.Context) error {
	signUpRequest := SignUpRequest{}

	bindErr := ctx.Bind(&signUpRequest)

	if bindErr != nil {
		return bindErr
	}

	if signUpRequest.Password != signUpRequest.ConfirmPassword {
		return errors.New("confirm password not equal to password")
	}

	exist, existErr := u.service.IsExistByUsername(signUpRequest.Username)

	if existErr != nil {
		return existErr
	}

	if exist {
		return errors.New("user already exist")
	}

	user := UserDTO{
		ID:          uuid.New().String(),
		Name:        signUpRequest.Name,
		Surname:     signUpRequest.Surname,
		Username:    signUpRequest.Username,
		Password:    base64.StdEncoding.EncodeToString([]byte(signUpRequest.Password)),
		CreatedDate: time.Now().String(),
	}

	createErr := u.service.CreateUser(user)

	if createErr != nil {
		return createErr
	}

	ctxErr := ctx.NoContent(201)

	if ctxErr != nil {
		return ctxErr
	}

	return nil
}
