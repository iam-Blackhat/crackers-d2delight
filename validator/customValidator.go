package validator

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func CustomValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("notnumeric", func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			// Reject if it contains only digits
			match, _ := regexp.MatchString(`^\d+$`, value)
			return !match
		})
	}
}
