package routes

import (
	"fmt"
	m "project/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	m.LogMiddleware(e)
	m.CorsMiddleware(e)
	e.Use(middleware.BodyLimit("10M"))

	e.Use(middleware.Static("public"))
	e.File("/dashboard", "public/views/dashboard/dashboard.html") // admin dashboard
	e.File("/user", "public/views/user/landing_page(user).html")  //homepage for registered user

	e.File("/register", "public/views/dashboard/user_register.html")                   // admin access : registering user
	e.File("/search-user", "public/views/dashboard/search_user.html")                  // admin access : search user
	e.File("/bibliography-list", "public/views/dashboard/bibliography_list.html")      // admin access : search Bibliography
	e.File("/add-bibliography", "public/views/dashboard/bibliography_input_form.html") // admin access : add Bibliography
	e.GET("/update-bibliography/:id", func(c echo.Context) error {
		id := c.Param("id")
		fmt.Printf("Serving /update-bibliography for item %s\n", id)
		return c.File("public/views/dashboard/bibliography_update_form.html")
	})

	e.File("/*", "public/index.html")
	e.File("/home", "public/index.html")
	e.File("/login", "public/login.html")
	e.File("/search-results", "public/search_results.html")
	e.GET("/items", func(c echo.Context) error {
		fmt.Println("Serving /items")
		return c.File("public/search_results_items.html")
	})
	e.File("/product-details", "public/product_details.html")
	e.GET("/product-details/*", func(c echo.Context) error {
		return c.File("public/product_details.html")
	})

	e.File("/view-univ-logo", "public/views")

	SetupPublicRoutes(e)
	SetupRegisteredRoutes(e)
	return e

}
