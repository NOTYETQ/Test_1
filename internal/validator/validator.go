package validator

import (
    "regexp"
    "strings"
    "unicode/utf8"
)

type Validator struct {
    Errors map[string]string
}

func NewValidator() *Validator {
    return &Validator{
        Errors: make(map[string]string),
    }
}

func (v *Validator) ValidData() bool {
    return len(v.Errors) == 0
}

func (v *Validator) AddError(field string, message string) {
    if _, exists := v.Errors[field]; !exists {
        v.Errors[field] = message
    }
}

func (v *Validator) Check(ok bool, field string, message string) {
    if !ok {
        v.AddError(field, message)
    }
}

// NotBlank returns true if data is present in the input box
func NotBlank(value string) bool {
    return strings.TrimSpace(value) != ""
}

// MaxLength returns true if the value contains no more than n characters
func MaxLength(value string, n int) bool {
    return utf8.RuneCountInString(value) <= n
}

// MinLength returns true if the value contains at least n characters
func MinLength(value string, n int) bool {
    return utf8.RuneCountInString(value) >= n
}

// IsValidNumeric returns true if the value is a valid number
func IsValidNumeric(value string) bool {
    numericRegex := regexp.MustCompile(`^-?\d+(\.\d+)?$`)
    return numericRegex.MatchString(value)
}

// IsValidDate returns true if the value is a valid date in format YYYY-MM-DD
func IsValidDate(value string) bool {
    dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
    return dateRegex.MatchString(value)
}