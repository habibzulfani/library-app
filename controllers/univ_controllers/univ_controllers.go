package univ_controllers

import (
	"net/http"
	"project/configs"
	"project/models/univ_models"

	"github.com/labstack/echo/v4"
)

func GetFakultasAndJurusan(c echo.Context) error {
	var fakultases []univ_models.Fakultas
	if err := configs.DB.Preload("Jurusans").Find(&fakultases).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch Fakultas and Jurusan data"})
	}
	return c.JSON(http.StatusOK, fakultases)
}
