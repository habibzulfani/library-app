package validators

import (
	userM "project/models/user_models"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func tenDigitsValidator(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	matched, _ := regexp.MatchString(`^\d{1,12}$`, field)
	return matched
}

// User validator
func ValidateUser(user *userM.User) error {
	validate := validator.New()
	return validate.Struct(user)
}

// Register validator

func ValidateUserRegister(user *userM.UserRegisterInput) error {
	validate := validator.New()

	// Register the custom validation
	// validate.RegisterValidation("ten_digits", tenDigitsValidator)

	// Validate the struct
	return validate.Struct(user)
}

// Login validator
// Register validator
func ValidateUserLogin(user *userM.UserLoginInput) error {
	validate := validator.New()
	return validate.Struct(user)
}

// Update validator
func ValidateUserUpdate(user *userM.UserUpdateInput) error {
	validate := validator.New()

	// Register the custom validation
	// validate.RegisterValidation("ten_digits", tenDigitsValidator)

	// Validate the struct
	return validate.Struct(user)
}

// Request validator
func ValidateUserResponse(user *userM.UserRequestResponse) error {
	validate := validator.New()
	return validate.Struct(user)
}
