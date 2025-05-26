package middlewares

import (
	"project/configs"
	userM "project/models/user_models"

	"github.com/labstack/echo/v4"
)

func BasicAuthMiddleware(username, password string, c echo.Context) (bool, error) {
	var user userM.User

	err := configs.DB.Where("email = ? AND password = ?", username, password).First(&user).Error
	if err != nil {
		return false, err
	}

	return true, err
}
