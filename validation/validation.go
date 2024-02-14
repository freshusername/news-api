package validation

import (
	"fmt"
	"reflect"
	"regexp"
)

// ValidationError wraps a validation rule error
type ValidationError struct {
	Field string
	Error string
}

func (e *ValidationError) PrintError() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Error)
}

// Rule defines a function signature for validation rules
type Rule func(value interface{}) *ValidationError

// Validator struct to hold validation rules
type Validator struct {
	rules map[string][]Rule
}

// NewValidator creates a new Validator instance
func NewValidator() *Validator {
	return &Validator{
		rules: make(map[string][]Rule),
	}
}

// AddRule adds a new validation rule for a field
func (v *Validator) AddRule(field string, rule Rule) {
	v.rules[field] = append(v.rules[field], rule)
}

// Validate executes the validation rules and returns errors
func (v *Validator) Validate(obj interface{}) []*ValidationError {
	var errors []*ValidationError

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i).Interface()

		rules, exists := v.rules[field.Name]
		if !exists {
			continue
		}

		for _, rule := range rules {
			if err := rule(value); err != nil {
				err.Field = field.Name
				errors = append(errors, err)
				break // Stop on the first error for each field
			}
		}
	}
	return errors
}

// Validation rules

// Required validates that the given value is not empty
func Required() Rule {
	return func(value interface{}) *ValidationError {
		if value == nil || value == "" {
			return &ValidationError{Error: "is required"}
		}
		return nil
	}
}

// Length validates the string length is within the specified range
func Length(min, max int) Rule {
	return func(value interface{}) *ValidationError {
		str, ok := value.(string)
		if !ok {
			return &ValidationError{Error: "is not a valid string"}
		}
		if len(str) < min || len(str) > max {
			return &ValidationError{Error: fmt.Sprintf("must be between %d and %d characters", min, max)}
		}
		return nil
	}
}

// Email validates the string is a valid email format
func Email() Rule {
	return func(value interface{}) *ValidationError {
		str, ok := value.(string)
		if !ok {
			return &ValidationError{Error: "is not a valid string"}
		}
		if match, _ := regexp.MatchString(`^\S+@\S+\.\S+$`, str); !match {
			return &ValidationError{Error: "is not a valid email address"}
		}
		return nil
	}
}
