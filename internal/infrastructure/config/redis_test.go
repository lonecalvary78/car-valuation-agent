package config

import "testing"

func TestGetAddrForRedis(t *testing.T) {
	redisConfig := Redis{
		Host: "0.0.0.0",
		Port: 6379,
	}

	expected := "0.0.0.0:6379"
	if actual := redisConfig.GetAddr(); actual != expected {
		t.Errorf("actual does not match with expected result[actual: %v, expected: %v]", actual, expected)
	}
}
