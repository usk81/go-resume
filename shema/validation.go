package shema

import (
	"regexp"
	"time"

	validator "gopkg.in/go-playground/validator.v9"
)

// SemVerRegex is the regular expression used to parse a semantic version.
const SemVerRegex string = `^v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`

// Validation validates Resume struct or a part of Resume struct
func Validation(v interface{}) error {
	vl := validator.New()
	vl.RegisterValidation("date", validateDateISO8601)
	vl.RegisterValidation("semver", validateSemVer)
	return vl.Struct(v)
}
func validateDateISO8601(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().String())
	return err == nil
}

func validateSemVer(fl validator.FieldLevel) bool {
	r := regexp.MustCompile(SemVerRegex)
	return r.MatchString(fl.Field().String())
}
