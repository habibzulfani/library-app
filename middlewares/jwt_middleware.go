package middlewares

import (
	"net/http"
	"project/constants"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JWTCustomClaims struct {
	UserID int    `json:"id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type UserData struct {
	UserID int
	Name   string
	Role   string
}

func CreateToken(user UserData) (string, error) {
	claims := &JWTCustomClaims{
		UserID: user.UserID,
		Name:   user.Name,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(constants.SECRET_JWT))
}

func JWTChecksRoleAdmin(c echo.Context) error {
	// Retrieve role from JWT token
	userData := c.Get("user").(*jwt.Token)
	claims := userData.Claims.(*JWTCustomClaims)
	role := claims.Role

	// Check if User's role is valid
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Request failed",
			"error":   "User unauthorized",
		})
	}

	return nil
}

func JWTChecksRoleUser(c echo.Context) error {
	// Retrieve role from JWT token
	userData := c.Get("user").(*jwt.Token)
	claims := userData.Claims.(*JWTCustomClaims)
	role := claims.Role

	// Check if User's role is valid
	if role != "admin" && role != "user"{
		return echo.NewHTTPError(http.StatusUnauthorized, "User unauthorized")
	}

	return nil
}
