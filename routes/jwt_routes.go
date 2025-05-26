package routes

import (
	"project/constants"
	"project/controllers"
	bookC "project/controllers/book_controllers/registered_books_access"
	paperC "project/controllers/paper_controllers/registered_papers_access"
	resC "project/controllers/resource_controllers"
	univC "project/controllers/univ_controllers"
	userC "project/controllers/user_controllers"

	m "project/middlewares"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func SetupRegisteredRoutes(e *echo.Echo) {
	eJWT := e.Group("/api")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(m.JWTCustomClaims)
		},
		SigningKey: []byte(constants.SECRET_JWT),
	}
	eJWT.Use(echojwt.WithConfig(config))

	RegisteredUserGroups(eJWT)
	RegisteredProfileGroups(eJWT)
	RegisteredBooksGroups(eJWT)
	RegisteredPapersGroups(eJWT)
	RegisteredUnivGroups(eJWT)
	RegisteredResourcesGroups(eJWT)
	RegisteredDownloadGroups(eJWT)
}

func RegisteredUserGroups(eJWT *echo.Group) {
	// Get all users
	eJWT.GET("/users", userC.GetAllUserController)

	// Register user
	eJWT.POST("/register", userC.Register)

	// Get user by id
	eJWT.GET("/users/:id", userC.GetUser)
	// Update user
	eJWT.PUT("/users/:id", userC.Update)

	// Delete single user
	eJWT.DELETE("/users/:id", userC.DeleteSingle)

	// Delete multiple users
	eJWT.DELETE("/users", userC.DeleteMultiple)
}

func RegisteredUnivGroups(eJWT *echo.Group) {
	// Get all users
	eJWT.GET("/fakultas-jurusan", univC.GetFakultasAndJurusan)

}
func RegisteredResourcesGroups(eJWT *echo.Group) {
	// Get all users
	eJWT.GET("/bibliographies", resC.GetResources)

	// Delete single/multiple papers or books
	eJWT.DELETE("/bibliographies", resC.DeleteBibliographies)

}

func RegisteredProfileGroups(eJWT *echo.Group) {
	// Get me user
	eJWT.GET("/profile", userC.GetUserProfile)

	// Update profile
	eJWT.PUT("/profile", userC.UpdateProfile)

	// Delete profile
	eJWT.DELETE("/profile", userC.DeleteProfile)
}

func RegisteredBooksGroups(eJWT *echo.Group) {
	// Get books
	eJWT.GET("/books", bookC.GetBooksRegistered)

	// Get books by id
	eJWT.GET("/books/:id", bookC.GetBook)

	// Create book
	eJWT.POST("/books", bookC.CreateBook)

	// Update book
	eJWT.PUT("/books/:id", bookC.UpdateBook)

	// Delete book
	eJWT.DELETE("/books/id", bookC.DeleteBook)
}

func RegisteredPapersGroups(eJWT *echo.Group) {
	// Get papers
	eJWT.GET("/papers", paperC.GetPapersRegistered)

	// Get papers by id
	eJWT.GET("/papers/:id", paperC.GetPaper)

	// Create paper
	eJWT.POST("/papers", paperC.CreatePaper)

	// Update paper
	eJWT.PUT("/papers/:id", paperC.UpdatePaper)

	// Delete paper
	// eJWT.DELETE("/papers/id", paperC.DeletePaper)
}

func RegisteredDownloadGroups(eJWT *echo.Group) {
	eJWT.GET("/:type/:id/download", controllers.DownloadItem)
}
