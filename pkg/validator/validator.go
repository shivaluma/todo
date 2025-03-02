package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// Validator is a wrapper around validator.Validate
type Validator struct {
	validate *validator.Validate
	trans    ut.Translator
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	// Always create a new validator instead of trying to reuse Gin's
	validate := validator.New()

	// Register the validator with gin but disable validation during binding
	binding.Validator = &defaultValidator{
		validate: validate,
		disableValidation: true,
	}

	// Register custom validation tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Set up translator
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)

	// Register custom error messages
	registerCustomErrorMessages(validate, trans)

	return &Validator{
		validate: validate,
		trans:    trans,
	}
}

// defaultValidator is a copy of gin's defaultValidator to allow us to replace it
type defaultValidator struct {
	validate *validator.Validate
	disableValidation bool
}

// ValidateStruct implements the gin.StructValidator interface
func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	if obj == nil {
		return nil
	}

	if v.disableValidation {
		return nil
	}

	value := reflect.ValueOf(obj)
	if value.Kind() == reflect.Ptr && !value.IsNil() {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return nil
	}

	return v.validate.Struct(obj)
}

// Engine returns the underlying validator engine
func (v *defaultValidator) Engine() interface{} {
	return v.validate
}

// Validate validates the given struct and returns validation errors
func (v *Validator) Validate(i interface{}) []ValidationError {
	err := v.validate.Struct(i)

	// Add debug logging in development environments only
	// fmt.Printf("Validation error: %v\n", err)

	if err == nil {
		return nil
	}

	var errors []ValidationError
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// Handle case where type assertion fails
		// fmt.Printf("Type assertion failed for validation errors\n")
		return []ValidationError{
			{Field: "unknown", Message: "Invalid input data"},
		}
	}

	// fmt.Printf("Number of validation errors: %d\n", len(validationErrors))

	for _, err := range validationErrors {
		// fmt.Printf("Field: %s, Tag: %s, Param: %s\n", err.Field(), err.Tag(), err.Param())

		// Use the JSON tag name directly
		field := err.Field()

		errors = append(errors, ValidationError{
			Field:   field,
			Message: err.Translate(v.trans),
		})
	}

	// if len(errors) == 0 {
	//    fmt.Printf("No errors were added to the errors slice\n")
	// }

	return errors
}

// registerCustomErrorMessages registers custom error messages for validation tags
func registerCustomErrorMessages(validate *validator.Validate, trans ut.Translator) {
	// Custom error message for required fields
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is required", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	// Custom error message for email validation
	validate.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email address", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	// Custom error message for min length validation
	validate.RegisterTranslation("min", trans, func(ut ut.Translator) error {
		return ut.Add("min", "{0} must be at least {1} characters long", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("min", fe.Field(), fe.Param())
		return t
	})

	// Custom error message for max length validation
	validate.RegisterTranslation("max", trans, func(ut ut.Translator) error {
		return ut.Add("max", "{0} must be at most {1} characters long", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("max", fe.Field(), fe.Param())
		return t
	})
}

// RegisterCustomValidation registers a custom validation function
func (v *Validator) RegisterCustomValidation(tag string, fn validator.Func, errMsg string) error {
	err := v.validate.RegisterValidation(tag, fn)
	if err != nil {
		return err
	}

	return v.validate.RegisterTranslation(tag, v.trans, func(ut ut.Translator) error {
		return ut.Add(tag, errMsg, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	})
}

// ValidateVar validates a single variable
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	err := v.validate.Var(field, tag)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok && len(validationErrors) > 0 {
			return fmt.Errorf(validationErrors[0].Translate(v.trans))
		}
		return err
	}
	return nil
}
