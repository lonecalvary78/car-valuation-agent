package middleware

import "strings"

// sanitizeForLog strips CR/LF characters from untrusted values before they reach a log
// sink, preventing log injection/forging via crafted request method or URI values.
func sanitizeForLog(s string) string {
	return strings.NewReplacer("\r", "", "\n", "").Replace(s)
}
