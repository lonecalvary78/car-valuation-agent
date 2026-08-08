package validator

import (
	"errors"
	"testing"
	"time"
)

type TestCase[K, V any] struct {
	Code           string
	AttributeName  K
	AttributeValue V
	ExpectedError  error
}

func TestValidateForRequiredOfString(t *testing.T) {
	testcases := []TestCase[string, string]{
		{
			Code:           "TC_VALIDATOR-STRING-001",
			AttributeName:  "Field1",
			AttributeValue: "ABC",
			ExpectedError:  nil,
		},
		{
			Code:           "TC_VALIDATOR-STRING-002",
			AttributeName:  "Field1",
			AttributeValue: "",
			ExpectedError:  errors.New("Field1 is required"),
		},
	}

	for _, testcase := range testcases {
		expectedError := testcase.ExpectedError
		actual := ValidateForRequiredOfString(testcase.AttributeName, testcase.AttributeValue)
		if (actual == nil) != (expectedError == nil) || (actual != nil && expectedError != nil && actual.Error() != expectedError.Error()) {
			t.Errorf("code: %v failed since the actual is not motch with expected[actual=%v, expected=%v]", testcase.Code, actual, expectedError)
		}
	}
}

func TestValidateForRequiredOfNumeric(t *testing.T) {
	testcases := []TestCase[string, int]{
		{
			Code:           "TC-VALIDATOR-INT-001",
			AttributeName:  "Field 1",
			AttributeValue: 12,
			ExpectedError:  nil,
		},
		{
			Code:           "TC-VALIDATOR-INT-001",
			AttributeName:  "Field 1",
			AttributeValue: 0,
			ExpectedError:  errors.New("Field 1 is required"),
		},
	}

	for _, testcase := range testcases {
		expectedError := testcase.ExpectedError
		actual := ValidateForRequiredOfNumeric(testcase.AttributeName, testcase.AttributeValue)
		if (actual == nil) != (expectedError == nil) || (actual != nil && expectedError != nil && actual.Error() != expectedError.Error()) {
			t.Errorf("code: %v failed since the actual is not motch with expected[actual=%v, expected=%v]", testcase.Code, actual, testcase.ExpectedError)
		}
	}
}

func TestValidateForRequiredOfDuration(t *testing.T) {
	emptyDuration, _ := time.ParseDuration("0s")
	testcases := []TestCase[string, time.Duration]{
		{
			Code:           "TC-VALIDATOR-INT-001",
			AttributeName:  "Field 1",
			AttributeValue: 30 * time.Second,
			ExpectedError:  nil,
		},
		{
			Code:           "TC-VALIDATOR-INT-001",
			AttributeName:  "Field 1",
			AttributeValue: emptyDuration,
			ExpectedError:  errors.New("Field 1 is required"),
		},
	}

	for _, testcase := range testcases {
		expectedError := testcase.ExpectedError

		actual := ValidateForRequiredOfDuration(testcase.AttributeName, testcase.AttributeValue)
		if (actual == nil) != (expectedError == nil) || (actual != nil && expectedError != nil && actual.Error() != expectedError.Error()) {
			t.Errorf("code: %v failed since the actual is not motch with expected[actual=%v, expected=%v]", testcase.Code, actual, testcase.ExpectedError)
		}
	}
}
