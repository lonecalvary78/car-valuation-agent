package config

import "testing"

func TestGetServerAddr(t *testing.T) {
	serverConfig := Server{
		Host: "0.0.0.0",
		Port: 8080,
	}

	expected := "0.0.0.0:8080"

	if actual := serverConfig.GetAddr(); actual != expected {
		t.Errorf("actual does not match with expected result[actual: %v, expected: %v]", actual, expected)
	}
}
