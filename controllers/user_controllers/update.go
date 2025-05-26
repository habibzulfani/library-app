package user_controllers

import (
	"net/http"
	"project/configs"
	"project/middlewares"
	userM "project/models/user_models"
	univM "project/models/univ_models"
	"project/validators"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Update(c echo.Context) error {
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

	var updateInput userM.UserUpdateInput
	if err := c.Bind(&updateInput); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := validators.ValidateUserUpdate(&updateInput); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	err = configs.DB.Transaction(func(tx *gorm.DB) error {
		if updateInput.Email != "" && updateInput.Email != user.Email {
			var existingUser userM.User
			if err := tx.Where("email = ? AND id != ?", updateInput.Email, inputID).First(&existingUser).Error; err == nil {
				return echo.NewHTTPError(http.StatusConflict, "Email is already in use")
			} else if err != gorm.ErrRecordNotFound {
				return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
			}
		}

		if updateInput.IDNumber != "" && updateInput.IDNumber != user.IDNumber {
			var existingUser userM.User
			if err := tx.Where("id_number = ? AND id != ?", updateInput.IDNumber, inputID).First(&existingUser).Error; err == nil {
				return echo.NewHTTPError(http.StatusConflict, "ID Number is already in use")
			} else if err != gorm.ErrRecordNotFound {
				return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
			}
		}

		if updateInput.Email != "" {
			user.Email = updateInput.Email
		}
		if updateInput.Name != "" {
			user.Name = updateInput.Name
		}
		if updateInput.UserType != "" {
			user.UserType = updateInput.UserType
		}
		if updateInput.IDNumber != "" {
			user.IDNumber = updateInput.IDNumber
		}
		if updateInput.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateInput.Password), bcrypt.DefaultCost)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process password")
			}
			user.Password = string(hashedPassword)
		}
		if updateInput.JurusanID != 0 {
			var jurusan univM.Jurusan
			if err := tx.First(&jurusan, updateInput.JurusanID).Error; err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid Jurusan ID")
			}
			user.JurusanID = updateInput.JurusanID
		}
		if updateInput.Address != "" {
			user.Address = updateInput.Address
		}
		if updateInput.Role != "" {
			user.Role = updateInput.Role
		}

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to update user data",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully updated user data",
		"user":   user,
	})
}

func UpdateProfile(c echo.Context) error {
	userProfile := c.Get("user").(*jwt.Token)
	claims := userProfile.Claims.(*middlewares.JWTCustomClaims)
	userID := claims.UserID

	if userID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid user",
		})
	}

	var user userM.User
	err := configs.DB.Preload("Jurusan.Fakultas").Where("id = ?", userID).First(&user).Error
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

	var updateInput userM.UserUpdateInput
	if err := c.Bind(&updateInput); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := validators.ValidateUserUpdate(&updateInput); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	err = configs.DB.Transaction(func(tx *gorm.DB) error {
		if updateInput.Email != "" && updateInput.Email != user.Email {
			var existingUser userM.User
			if err := tx.Where("email = ? AND id != ?", updateInput.Email, userID).First(&existingUser).Error; err == nil {
				return echo.NewHTTPError(http.StatusConflict, "Email is already in use")
			} else if err != gorm.ErrRecordNotFound {
				return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
			}
		}

		if updateInput.Email != "" {
			user.Email = updateInput.Email
		}
		if updateInput.Name != "" {
			user.Name = updateInput.Name
		}
		if updateInput.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateInput.Password), bcrypt.DefaultCost)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process password")
			}
			user.Password = string(hashedPassword)
		}
		if updateInput.JurusanID != 0 {
			var jurusan univM.Jurusan
			if err := tx.First(&jurusan, updateInput.JurusanID).Error; err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid Jurusan ID")
			}
			user.JurusanID = updateInput.JurusanID
		}
		if updateInput.Address != "" {
			user.Address = updateInput.Address
		}

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to update user data",
			"error":   err.Error(),
		})
	}

	userResponse := userM.UserUpdateResponse{
		ID:        int(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		UserType:  user.UserType,
		IDNumber:  user.IDNumber,
		// Fakultas:  user.Jurusan.Fakultas.Name,
		// Jurusan:   user.Jurusan.Name,
		Address:   user.Address,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully updated user data",
		"user":    userResponse,
	})
}