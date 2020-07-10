package user

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type mockUserService struct {
}

func (m *mockUserService) CreateUser(user UserDTO) error {
	mockDB[user.Username] = &user
	return nil
}

func (m *mockUserService) IsExistUsernameAndPassword(username, password string) (bool, error) {
	if val, ok := mockDB[username]; ok {
		if strings.EqualFold(val.Password, password) {
			return true, nil
		}
		return false, nil
	}
	return false, nil

}

func (m mockUserService) IsExistByUsername(username string) (bool, error) {
	_, ok := mockDB[username]
	return ok, nil
}

var (
	userId string = uuid.New().String()
	mockDB        = map[string]*UserDTO{
		"unittest": &UserDTO{
			ID:          userId,
			Name:        "test",
			Surname:     "golang",
			Username:    "unittest",
			Password:    "123456",
			CreatedDate: time.Now().String(),
		},
		"benchmarkTest": &UserDTO{
			ID:          userId,
			Name:        "testSon",
			Surname:     "golangSon",
			Username:    "benchmarkTest",
			Password:    "123456",
			CreatedDate: time.Now().String(),
		},
	}
	userExist = `{
	"name":"test",
	"surname":"golang",
	"username":"unittest",
	"password":"123456",
	"confirmPassword":"123456"
}`
	userNew = `{
	"name":"testNew",
	"surname":"golangNew",
	"username":"unittestnew",
	"password":"123456",
	"confirmPassword":"123456"
}`
	userBlind = `{
	"name":"test",
	"surname":"golangBlind",
	"username":"unittestBlind",
	"password":"123456",
	"confirmPassword":"12356"
}`
)

func TestSignUpHandler(t *testing.T) {
	e := echo.New()
	user := User{server: e, service: new(mockUserService)}

	// Test func with already exist user
	c, rec, _ := createContextWithRequest(e, userExist)
	// Assertions
	if assert.Errorf(t, user.signUpHandler(c), "user already exist") {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	//Test func with new user
	c, rec, _ = createContextWithRequest(e, userNew)
	// Assertions
	if assert.NoError(t, user.signUpHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	//Test func with user that doesnot matching Password and ConfirmPassword
	c, rec, _ = createContextWithRequest(e, userBlind)
	// Assertions
	if assert.Errorf(t, user.signUpHandler(c), "confirm password not equal to password") {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

}

func createContextWithRequest(e *echo.Echo, json string) (echo.Context, *httptest.ResponseRecorder, *http.Request) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec, req
}
