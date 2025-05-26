package validators

import (
	paperM "project/models/paper_models"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Paper validator
func ValidatePaper(paper *paperM.Paper) error {
	validate := validator.New()
	return validate.Struct(paper)
}
func ValidateCreatePaper(paper *paperM.CreatePaperInput) error {
	validate := validator.New()
	_ = validate.RegisterValidation("issn", validateISSN)
	_ = validate.RegisterValidation("doi", validateDOI)
	return validate.Struct(paper)
}
func ValidatePaperUpdate(paper *paperM.PaperUpdateInput) error {
	validate := validator.New()
	_ = validate.RegisterValidation("issn", validateUpdateISSN)
	_ = validate.RegisterValidation("doi", validateDOI)

	return validate.Struct(paper)
}

func validateISSN(fl validator.FieldLevel) bool {
	issn := fl.Field().String()
	value := strings.ReplaceAll(issn, "-", "") // Remove hyphens for validation
	return isValidISSN(value)
}
func validateUpdateISSN(fl validator.FieldLevel) bool {
	issn := fl.Field().String()
	if issn == "" {
		return true // Optional for updates
	}
	value := strings.ReplaceAll(issn, "-", "") // Remove hyphens for validation
	return isValidISSN(value)
}

func isValidISSN(issn string) bool {
	if len(issn) != 8 {
		return false
	}
	sum := 0
	for i := 0; i < 7; i++ {
		digit, err := strconv.Atoi(string(issn[i]))
		if err != nil {
			return false
		}
		sum += digit * (8 - i)
	}
	checkDigit := 11 - (sum % 11)
	if checkDigit == 10 {
		return issn[7] == 'X'
	}
	if checkDigit == 11 {
		return issn[7] == '0'
	}
	return issn[7] == byte(checkDigit+'0')
}
