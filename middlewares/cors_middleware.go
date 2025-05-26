package middlewares

import (
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CorsMiddleware(e *echo.Echo) {
	// Get allowed origins from environment variable
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			// Check if the origin is in the list of allowed origins
            for _, allowedOrigin := range allowedOrigins {
                if origin == strings.TrimSpace(allowedOrigin) {
                    return true, nil
                }
            }

			// For development, allow localhost origins
			if strings.HasPrefix(origin, "http://localhost:") || strings.HasPrefix(origin, "http://127.0.0.1:") {
				return true, nil
			}
			return false, nil
		},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
}
