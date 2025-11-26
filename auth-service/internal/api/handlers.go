package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Bkgediya/feed_system/auth-service/internal/auth"
	"github.com/Bkgediya/feed_system/auth-service/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	svc       service.AuthService
	log       *zap.SugaredLogger
	jwtSecret string
}

func NewHandler(svc service.AuthService, l *zap.SugaredLogger) *Handler {
	return &Handler{svc: svc, log: l, jwtSecret: ""}
}

type signupReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Signup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signupReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		u, err := h.svc.SignUp(req.Username, req.Email, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res := map[string]interface{}{"id": u.ID, "username": u.Username, "email": u.Email}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		token, err := h.svc.Login(req.Email, req.Password)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		res := map[string]string{"token": token}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authz := r.Header.Get("Authorization")
		if authz == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
		}

		parts := strings.SplitN(authz, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		claims, err := auth.ParseJWT(parts[1], h.jwtSecret)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) Me() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("user_id").(int64)
		u, err := h.svc.GetUser(uid)
		if err != nil || u == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": u.ID, "email": u.Email, "username": u.Username})
	}
}
