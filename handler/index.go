package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/sunchero/http/service"
)

type Handler struct {
	Svc *service.Service
	NC  *nats.Conn
}

func (h *Handler) New() http.Handler {
	router := mux.NewRouter()

	AuthRouter := router.PathPrefix("/streams").Subrouter()
	AuthRouter.Use(JWTMiddleware)
	AuthRouter.HandleFunc("/streams", h.createStream).Methods("POST")
	AuthRouter.HandleFunc("/streams", h.updateStreams).Methods("PATCH")
	AuthRouter.HandleFunc("/streams/list", h.listStreams).Methods("GET")
	AuthRouter.HandleFunc("/streams/delete/{id}", h.deleteStreams).Methods("DELETE")

	// router.HandleFunc("/consumers/create", h.createConsumer).Methods("POST")
	// router.HandleFunc("/consumers", h.deleteConsumer).Methods("DELETE")

	AuthRouter.HandleFunc("/streams/join/{stream}", h.joinStream)
	AuthRouter.HandleFunc("/streams/leave/{stream}", h.leaveStream)

	UnAuthRouter := router.PathPrefix("/ws").Subrouter()
	UnAuthRouter.HandleFunc("/listen", h.handleWebSocket)

	return router

}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Parse the token and validate the signature
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Provide the secret key used to sign the token
			secret := []byte("my-secret-key")
			return secret, nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Extract the "name" claim from the token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if name, ok := claims["name"].(string); ok {
				// Add the name to the HTTP context
				ctx := context.WithValue(r.Context(), "name", name)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		w.WriteHeader(http.StatusUnauthorized)
	})
}
