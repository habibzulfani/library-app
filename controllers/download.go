package controllers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"project/configs"
	"project/middlewares"
	"project/models/book_models"
	"project/models/paper_models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func getMIMEType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		return "application/octet-stream"
	}
}

func DownloadItem(c echo.Context) error {
	fmt.Printf("Download request received. Type: %s, ID: %s\n", c.Param("type"), c.Param("id"))
	if err := middlewares.JWTChecksRoleUser(c); err != nil {
		return err
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*middlewares.JWTCustomClaims)
	userID := claims.UserID
	if userID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid user",
		})
	}

	itemType := c.Param("type")
	itemIDStr := c.Param("id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid id",
			"error":   err.Error(),
		})
	}

	var filePath string
	var fileName string

	switch itemType {
	case "book":
		var book book_models.Book
		if err := configs.DB.First(&book, itemID).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Book not found"})
		}
		filePath = book.FileURL
		fileName = book.Title

	case "paper":
		var paper paper_models.Paper
		if err := configs.DB.First(&paper, itemID).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Paper not found"})
		}
		filePath = paper.FileURL
		fileName = paper.Title

	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid item type"})
	}

	// Determine the file extension
	fileExt := filepath.Ext(filePath)
	if fileExt == "" {
		// If no extension, default to .pdf
		fileExt = ".pdf"
	}

	fileName = strings.TrimSpace(fileName)
	fileName = strings.ReplaceAll(fileName, "/", "_")
	fileName = strings.ReplaceAll(fileName, "\\", "_")
	fileName = fmt.Sprintf("%s%s", fileName, fileExt)

	// Ensure the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "File not found"})
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to open file"})
	}
	defer file.Close()

	// Get file information
	fileInfo, err := file.Stat()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to get file info"})
	}

	// Log the download activity
	err = middlewares.LogDownloadActivity(configs.DB, uint(userID), uint(itemID))
	if err != nil {
		println("Failed to log download activity:", err.Error())
	}

	// Increment the download counter
	err = middlewares.IncrementDownloadCounter(configs.DB)
	if err != nil {
		println("Failed to increment download counter:", err.Error())
	}

	// Determine the MIME type
	mimeType := getMIMEType(fileName)

	// Set headers for file download
	encodedFileName := url.PathEscape(fileName)
	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"; filename*=UTF-8''%s", fileName, encodedFileName))
	c.Response().Header().Set(echo.HeaderContentType, mimeType)
	c.Response().Header().Set(echo.HeaderContentLength, fmt.Sprintf("%d", fileInfo.Size()))

	// Stream the file to the client
	_, err = io.Copy(c.Response().Writer, file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to stream file"})
	}

	return nil
}
