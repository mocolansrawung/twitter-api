package middleware

import (
	"net/http"

	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
)

func (a *Authentication) ApiKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-Api-Key")
		viper.AutomaticEnv()
		validApiKey := viper.GetString("API_KEY")

		requestID, _ := uuid.NewV4()
		requestIDStr := requestID.String()

		w.Header().Set("X-Request-Id", requestIDStr)

		if apiKey != validApiKey {
			response.WithMessage(w, http.StatusUnauthorized, "You need valid API Key to perform this action")
			return
		}

		next.ServeHTTP(w, r)
	})
}
