package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidPassword = regexp.MustCompile(`^[a-zA-Z0-9!@#$&*]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z]{2,}\s[a-zA-Z]{1,}'?-?[a-zA-Z]{2,}\s?([a-zA-Z]{1,})?$`).MatchString
)

func ValidateStringField(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("length must be between %d and %d", minLength, maxLength)
	}

	return nil
}

func ValidateUsername(username string) error {
	err := ValidateStringField(username, 3, 50)
	if err != nil {
		return err
	}
	if !isValidUsername(username) {
		return fmt.Errorf("username must not contain special characters")
	}
	return nil
}

func ValidatePassword(password string) error {
	err := ValidateStringField(password, 6, 50)
	if err != nil {
		return err
	}
	if !isValidPassword(password) {
		return fmt.Errorf("password must contain special characters")
	}
	return nil
}

func ValidateEmail(emails string) error {
	err := ValidateStringField(emails, 6, 50)
	if err != nil {
		return err
	}
	if _, err = mail.ParseAddress(emails); err != nil {
		return fmt.Errorf("invalid email address")
	}
	return nil
}

func ValidateFullName(fullName string) error {
	err := ValidateStringField(fullName, 6, 50)
	if err != nil {
		return err
	}
	if !isValidFullName(fullName) {
		return fmt.Errorf("fullName must contain only alphabets and spaces")
	}
	return nil
}
