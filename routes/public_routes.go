package routes

import (
	"project/controllers/book_controllers"
	bookC "project/controllers/book_controllers"
	"project/controllers/paper_controllers"
	paperC "project/controllers/paper_controllers"
	resC "project/controllers/resource_controllers"
	userC "project/controllers/user_controllers"

	"github.com/labstack/echo/v4"
)

func SetupPublicRoutes(e *echo.Echo) {
	// Login user
	e.POST("/login", userC.LoginController)

	// Get books as public
	e.GET("/public/books", bookC.GetBooksPublic)

	// Get papers as public
	e.GET("/public/papers", paperC.GetPapersPublic)

	// Search data as public
	e.GET("/public/search", resC.SearchBar)

	// advanced search data as public
	e.GET("/public/advanced-search", resC.AdvancedSearch)

	// showing single product details
	e.GET("/public/advanced-search", resC.AdvancedSearch)
	
	// Get books by id
	e.GET("public/books/:id", book_controllers.GetBook)
	
	// Get papers by id
	e.GET("public/papers/:id", paper_controllers.GetPaper)
	
	// Get all papers
	e.GET("public/papers", paper_controllers.GetPapersPublic)
	
	// Get all books
	e.GET("public/books", book_controllers.GetBooksPublic)


	// e.GET("/search-results", func(c echo.Context) error {
	// 	log.Printf("Received request for search results page: %s", c.Request().URL)
	// 	return c.File("public/search_results.html")
	// })

}
