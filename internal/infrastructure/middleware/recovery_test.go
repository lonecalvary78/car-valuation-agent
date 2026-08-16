package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecovery(t *testing.T) {
	t.Run("logs the panic and hides it from the client", func(t *testing.T) {
		var logBuf bytes.Buffer
		log.SetOutput(&logBuf)
		t.Cleanup(func() { log.SetOutput(os.Stderr) })

		h := Recovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("secret internal detail")
		}))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequestWithContext(t.Context(), "GET", "/boom", nil))

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.NotContains(t, w.Body.String(), "secret internal detail")
		require.Contains(t, logBuf.String(), "secret internal detail")
		require.Contains(t, logBuf.String(), "recovery.go")
	})
}
