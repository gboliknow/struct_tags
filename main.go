package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type User struct {
	Name  string `validate:"min=2,max=32"`
	Email string `validate:"required,email"`
}

func validate(s interface{}) error {
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := v.Type().Field(i).Name
		tag := v.Type().Field(i).Tag.Get("validate")
		if tag == "" {
			continue
		}

		rules := strings.Split(tag, ",")
		for _, rule := range rules {
			if err := applyRule(rule, field, fieldName); err != nil {
				return err
			}
		}
	}
	return nil
}

func applyRule(rule string, field reflect.Value, fieldName string) error {
	switch {
	case strings.HasPrefix(rule, "min="):
		return validateMinLength(rule, field, fieldName)
	case strings.HasPrefix(rule, "max="):
		return validateMaxLength(rule, field, fieldName)
	case rule == "required":
		return validateRequired(field, fieldName)
	case rule == "email":
		return validateEmail(field, fieldName)
	}
	return nil
}

func validateMinLength(rule string, field reflect.Value, fieldName string) error {
	min, _ := strconv.Atoi(strings.TrimPrefix(rule, "min="))
	if len(field.String()) < min {
		return fmt.Errorf("%s must be at least %d characters long", fieldName, min)
	}
	return nil
}

func validateMaxLength(rule string, field reflect.Value, fieldName string) error {
	max, _ := strconv.Atoi(strings.TrimPrefix(rule, "max="))
	if len(field.String()) > max {
		return fmt.Errorf("%s must be at most %d characters long", fieldName, max)
	}
	return nil
}

func validateRequired(field reflect.Value, fieldName string) error {
	if field.String() == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

func validateEmail(field reflect.Value, fieldName string) error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(field.String()) {
		return fmt.Errorf("%s must be a valid email address", fieldName)
	}
	return nil
}

func main() {
	validUser := User{
		Name:  "Alice",
		Email: "alice@example.com",
	}

	invalidUser := User{
		Name:  "A",
		Email: "aliceexample.com",
	}

	printValidationResult(validUser)
	printValidationResult(invalidUser)
}

func printValidationResult(user User) {
	err := validate(user)
	if err != nil {
		fmt.Println("Validation error:", err)
	} else {
		fmt.Println("User is valid")
	}
}