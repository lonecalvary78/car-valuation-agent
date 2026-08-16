package validator

import (
	"errors"
	"fmt"
	"time"
)

var ErrRequired = errors.New("is required")

func ValidateForRequiredOfString(attributeName string, atrributeValue string) error {
	if len(atrributeValue) == 0 {
		return fmt.Errorf("%v %w", attributeName, ErrRequired)
	}
	return nil
}

func ValidateForRequiredOfNumeric(attributeName string, atrributeValue int) error {
	if atrributeValue <= 0 {
		return fmt.Errorf("%v %w", attributeName, ErrRequired)
	}
	return nil
}

func ValidateForRequiredOfDuration(attributeName string, atrributeValue time.Duration) error {
	if atrributeValue <= 0 {
		return fmt.Errorf("%v %w", attributeName, ErrRequired)
	}
	return nil
}
