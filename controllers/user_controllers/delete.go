package user_controllers

import (
	"net/http"
	"project/configs"
	"project/middlewares"
	userM "project/models/user_models"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type DeleteRequest struct {
	UserIDs []int `json:"userIds"`
}

func DeleteSingle(c echo.Context) error {
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
	if err := configs.DB.First(&user, inputID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "User not found",
		})
	}

	// Soft delete
	if err := configs.DB.Delete(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to delete user",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully deleted user data",
	})

}

// func DeleteMultiple(c echo.Context) error {
// 	if err := middlewares.JWTChecksRoleAdmin(c); err != nil {
// 		return err
// 	}

// 	var req DeleteRequest
// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]interface{}{
// 			"message": "Invalid request body",
// 			"error":   err.Error(),
// 		})
// 	}

// 	if len(req.UserIDs) == 0 {
// 		return c.JSON(http.StatusBadRequest, map[string]interface{}{
// 			"message": "No user IDs provided",
// 		})
// 	}

// 	tx := configs.DB.Begin()
//     defer func() {
//         if r := recover(); r != nil {
//             tx.Rollback()
//         }
//     }()

//     var deletedCount int64
//     for _, id := range req.UserIDs {
//         var user userM.User
//         if err := tx.First(&user, id).Error; err != nil {
//             continue // Skip if user not found
//         }
//         if err := tx.Delete(&user).Error; err != nil {
//             tx.Rollback()
//             return c.JSON(http.StatusInternalServerError, map[string]interface{}{
//                 "message": "Failed to delete users",
//                 "error":   err.Error(),
//             })
//         }
//         deletedCount++
//     }

// 	if err := tx.Commit().Error; err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
// 			"message": "Failed to commit transaction",
// 			"error":   err.Error(),
// 		})
// 	}

// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"message":       "Successfully deleted users",
// 		"deleted_count": deletedCount,
// 	})
// }

// // In your main.go or database initialization file
// func init() {
// 	// Enable soft delete
// 	configs.DB.Use(gorm.Clause(&gorm.Config{
// 		PrepareStmt: true,
// 	}))
// }

// DeleteMultiple function with soft delete
func DeleteMultiple(c echo.Context) error {
	if err := middlewares.JWTChecksRoleAdmin(c); err != nil {
		return err
	}

	var req DeleteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if len(req.UserIDs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "No user IDs provided",
		})
	}

	tx := configs.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var deletedCount int64
	result := tx.Model(&userM.User{}).Where("id IN ?", req.UserIDs).Update("deleted_at", time.Now())
	if result.Error != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to delete users",
			"error":   result.Error.Error(),
		})
	}
	deletedCount = result.RowsAffected

	if err := tx.Commit().Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to commit transaction",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "Successfully soft-deleted users",
		"deleted_count": deletedCount,
	})
}

// If you need to force delete (hard delete)
func ForceDeleteUser(db *gorm.DB, userID uint) error {
	return db.Unscoped().Delete(&userM.User{}, userID).Error
}

// To query including soft-deleted records
func GetAllUsersIncludingDeleted(db *gorm.DB) ([]userM.User, error) {
	var users []userM.User
	err := db.Unscoped().Find(&users).Error
	return users, err
}

func DeleteProfile(c echo.Context) error {
	// Retrieve id from JWT token
	userProfile := c.Get("user").(*jwt.Token)
	claims := userProfile.Claims.(*middlewares.JWTCustomClaims)
	UserID := claims.UserID

	// Check if UserID is valid
	if UserID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid user",
		})
	}

	var user userM.User
	result := configs.DB.Delete(&user, UserID)

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "User not found",
		})
	}

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to delete user",
			"error":   result.Error.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"search": "Successfully deleted user data",
	})
}
