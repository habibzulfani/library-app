package validators

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

func validateDOI(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr && !field.IsNil() {
		doi := field.Elem().String()
		return isValidDOI(doi)
	}
	return true // Field is nil, which is valid for an optional field
}

func isValidDOI(doi string) bool {
	// Remove URL prefix if present
	doi = strings.TrimPrefix(doi, "https://doi.org/")

	// Basic DOI format validation
	doiRegex := regexp.MustCompile(`^10\.\d{4,9}/[-._;()/:A-Z0-9]+$`)
	return doiRegex.MatchString(doi)
}
