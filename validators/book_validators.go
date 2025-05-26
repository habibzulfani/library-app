package validators

import (
	bookM "project/models/book_models"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Book validator
func ValidateBook(book *bookM.Book) error {
	validate := validator.New()
	return validate.Struct(book)
}

// Update validator
// func ValidateBookUpdate(book *bookM.BookUpdateInput) error {
// 	validate := validator.New()
// 	return validate.Struct(book)
// }

// Create book validator

func ValidateCreateBook(book *bookM.CreateBookInput) error {
	validate := validator.New()
	_ = validate.RegisterValidation("isbn", validateISBN)
	_ = validate.RegisterValidation("doi", validateDOI)
	return validate.Struct(book)
}
func ValidateBookUpdate(book *bookM.BookUpdateInput) error {
	validate := validator.New()
	_ = validate.RegisterValidation("isbn", validateUpdateISBN)
	_ = validate.RegisterValidation("doi", validateDOI)
	return validate.Struct(book)
}

func validateISBN(fl validator.FieldLevel) bool {
	isbn := fl.Field().String()
	value := strings.ReplaceAll(isbn, "-", "") // Remove hyphens for validation
	if len(value) == 10 {
		return isValidISBN10(value)
	}
	if len(value) == 13 {
		return isValidISBN13(value)
	}
	return false
}

func validateUpdateISBN(fl validator.FieldLevel) bool {
	isbn := fl.Field().String()
	if isbn == "" {
		return true // Optional for updates
	}
	value := strings.ReplaceAll(isbn, "-", "") // Remove hyphens for validation
	if len(value) == 10 {
		return isValidISBN10(value)
	}
	if len(value) == 13 {
		return isValidISBN13(value)
	}
	return false
}


func isValidISBN10(isbn string) bool {
	if len(isbn) != 10 {
		return false
	}
	sum := 0
	for i := 0; i < 9; i++ {
		digit, err := strconv.Atoi(string(isbn[i]))
		if err != nil {
			return false
		}
		sum += digit * (10 - i)
	}
	checkDigit := 11 - (sum % 11)
	if checkDigit == 10 {
		return isbn[9] == 'X'
	}
	if checkDigit == 11 {
		return isbn[9] == '0'
	}
	return isbn[9] == byte(checkDigit+'0')
}

func isValidISBN13(isbn string) bool {
	if len(isbn) != 13 {
		return false
	}
	sum := 0
	for i := 0; i < 12; i++ {
		digit, err := strconv.Atoi(string(isbn[i]))
		if err != nil {
			return false
		}
		if i%2 == 0 {
			sum += digit
		} else {
			sum += digit * 3
		}
	}
	checkDigit := 10 - (sum % 10)
	if checkDigit == 10 {
		checkDigit = 0
	}
	return isbn[12] == byte(checkDigit+'0')
}
