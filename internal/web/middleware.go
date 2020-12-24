package web

import (
	"context"
	"net/http"
)

func (c *Controller) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			sess, err := c.sessionManager.GetSession(r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else {
				ctx := context.WithValue(r.Context(), "sessionKey", sess)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		},
	)
}
