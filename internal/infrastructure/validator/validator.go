package validator

import (
	"fmt"
	"time"
)

func ValidateForRequiredOfString(attributeName string, atrributeValue string) error {
	if len(atrributeValue) == 0 {
		return fmt.Errorf("%v is required", attributeName)
	}
	return nil
}

func ValidateForRequiredOfNumeric(attributeName string, atrributeValue int) error {
	if atrributeValue <= 0 {
		return fmt.Errorf("%v is required", attributeName)
	}
	return nil
}

func ValidateForRequiredOfDuration(attributeName string, atrributeValue time.Duration) error {
	if atrributeValue <= 0 {
		return fmt.Errorf("%v is required", attributeName)
	}
	return nil
}
