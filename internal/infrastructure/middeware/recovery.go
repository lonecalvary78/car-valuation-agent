package middeware

import (
	"fmt"
	"net/http"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := recover()
		if err != nil {
			error := fmt.Errorf("%v", err)
			http.Error(w, error.Error(), http.StatusInternalServerError)
		}
		next.ServeHTTP(w, r)
	})
}
