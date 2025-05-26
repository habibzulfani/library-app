package user_controllers

import (
	"fmt"
	"net/http"
	"project/configs"
	mid "project/middlewares"
	univM "project/models/univ_models"
	userM "project/models/user_models"
	"project/validators"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register1(c echo.Context) error {
	if err := mid.JWTChecksRoleAdmin(c); err != nil {
		return err
	}

	var userInput userM.UserRegisterInput
	if err := c.Bind(&userInput); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := validators.ValidateUserRegister(&userInput); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	err := configs.DB.Transaction(func(tx *gorm.DB) error {
		var existingUser userM.User
		if err := tx.Where("email = ?", userInput.Email).First(&existingUser).Error; err == nil {
			return echo.NewHTTPError(http.StatusConflict, "Email is already in use")
		} else if err != gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
		}

		if err := tx.Where("id_number = ?", userInput.IDNumber).First(&existingUser).Error; err == nil {
			return echo.NewHTTPError(http.StatusConflict, "ID Number is already in use")
		} else if err != gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
		}

		// var jurusan univM.Jurusan
		// if err := tx.Preload("Fakultas").First(&jurusan, userInput.JurusanID).Error; err != nil {
		// 	return echo.NewHTTPError(http.StatusBadRequest, "Invalid Jurusan ID")
		// }

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process password")
		}

        var jurusan univM.Jurusan
        if err := tx.First(&jurusan, userInput.JurusanID).Error; err != nil {
            return echo.NewHTTPError(http.StatusBadRequest, "Invalid Jurusan ID")
        }
    
        user := userM.User{
            Email:     userInput.Email,
            Name:      userInput.Name,
            UserType:  userInput.UserType,
            IDNumber:  userInput.IDNumber,
            JurusanID: userInput.JurusanID,
            Jurusan:   jurusan.Name,  // Set the Jurusan name
            Password:  string(hashedPassword),
            Address:   userInput.Address,
            Role:      userInput.Role,
            // NIM:       userInput.IDNumber,  // Assuming NIM is the same as IDNumber for students
        }
    
		// if err := tx.Create(&user).Error; err != nil {
		// 	return err
		// }

		fmt.Printf("User to be created: %+v\n", user)

		// Enable GORM debug mode
		if err := tx.Debug().Create(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		httpError, ok := err.(*echo.HTTPError)
		if ok {
			return c.JSON(httpError.Code, map[string]interface{}{
				"message": httpError.Message,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to create user",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
	})
}
