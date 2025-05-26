package user_controllers

import (
	"net/http"
	"strconv"

	"project/configs"
	"project/middlewares"
	userM "project/models/user_models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetUser(c echo.Context) error {
	if err := middlewares.JWTChecksRoleAdmin(c); err != nil {
		return err
	}

	inputID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid id",
			"error":   err.Error(),
		})
	}

	var user userM.User
	err = configs.DB.Preload("Jurusan.Fakultas").Where("id = ?", inputID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to retrieve user data",
			"error":   err.Error(),
		})
	}

	userResponse := userM.UserRequestResponse{
		ID:        int(user.ID),
		Email:     user.Email,
		Name:      user.Name,
		UserType:  user.UserType,
		IDNumber:  user.IDNumber,
		// Fakultas:  user.Jurusan.Fakultas.Name,
		// Jurusan:   user.Jurusan.Name,
		Address:   user.Address,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully retrieved user data",
		"user":    userResponse,
	})
}

func GetUserProfile(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*middlewares.JWTCustomClaims)
	UserID := claims.UserID
	if UserID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid user",
		})
	}

	var userData userM.User
	err := configs.DB.Preload("Jurusan.Fakultas").Where("id = ?", UserID).First(&userData).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to retrieve user data",
			"error":   err.Error(),
		})
	}

	userResponse := userM.UserRequestResponse{
		ID:        int(userData.ID),
		Email:     userData.Email,
		Name:      userData.Name,
		UserType:  userData.UserType,
		IDNumber:  userData.IDNumber,
		// Fakultas:  userData.Jurusan.Fakultas.Name,
		// Jurusan:   userData.Jurusan.Name,
		Address:   userData.Address,
		Role:      userData.Role,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully retrieved user data",
		"user":    userResponse,
	})
}

func GetAllUserController(c echo.Context) error {
    if err := middlewares.JWTChecksRoleAdmin(c); err != nil {
        return err
    }

    var users []userM.User

    err := configs.DB.Preload("JurusanRef.Fakultas").Find(&users).Error
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": err.Error(),
        })
    }

    userResponse := make([]userM.UserRequestResponse, len(users))
    for i, user := range users {
        userResponse[i] = userM.UserRequestResponse{
            ID:        int(user.ID),
            Email:     user.Email,
            Name:      user.Name,
            UserType:  user.UserType,
            IDNumber:  user.IDNumber,
            Fakultas:  user.JurusanRef.Fakultas.Name,
            Jurusan:   user.JurusanRef.Name,
            Address:   user.Address,
            Role:      user.Role,
            CreatedAt: user.CreatedAt,
            UpdatedAt: user.UpdatedAt,
        }
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "message": "success",
        "data":    userResponse,
    })
}