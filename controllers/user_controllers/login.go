package user_controllers

import (
	"net/http"
	"project/configs"
	m "project/middlewares"
	userM "project/models/user_models"
	"project/validators"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginController(c echo.Context) error {
	var loginInput userM.UserLoginInput
	if err := c.Bind(&loginInput); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := validators.ValidateUserLogin(&loginInput); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	var user userM.User
	err := configs.DB.Preload("JurusanRef.Fakultas").Where("email = ?", loginInput.Email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "Login failed",
				"error":   "Invalid email or password",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Login failed",
			"error":   "Database error",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Login failed",
			"error":   "Invalid email or password",
		})
	}

	userData := m.UserData{
		UserID: int(user.ID),
		Name:   user.Name,
		Role:   user.Role,
	}

	token, err := m.CreateToken(userData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Login failed",
			"error":   "Failed to generate token",
		})
	}

	userResponse := userM.UserLoginResponse{
		Id:       int(user.ID),
		Email:    user.Email,
		Name:     user.Name,
		UserType: user.UserType,
		IDNumber: user.IDNumber,
		Fakultas: user.JurusanRef.Fakultas.Name,
		Jurusan:  user.JurusanRef.Name,
		Address:  user.Address,
		Role:     user.Role,
		Token:    token,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user":    userResponse,
	})
}