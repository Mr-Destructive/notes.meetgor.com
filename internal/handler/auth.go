package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"blog/internal/db"
	"blog/internal/util"

	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// HandleLogin authenticates user and returns JWT
func HandleLogin(w http.ResponseWriter, r *http.Request, database *db.DB) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	
	// Get password hash from env (in real app, store in DB)
	passwordHash := os.Getenv("PASSWORD_HASH")
	
	if passwordHash != "" {
		// Use hashed password
		if !util.CheckPassword(passwordHash, req.Password) {
			respondError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
	} else {
		// Try to use raw password from ADMIN_PASSWORD env var
		adminPassword := os.Getenv("ADMIN_PASSWORD")
		if adminPassword == "" || req.Password != adminPassword {
			respondError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
	}

	// Generate JWT token
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin": true,
		"exp":   expiresAt.Unix(),
		"iat":   time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Set secure cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	respondJSON(w, http.StatusOK, LoginResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	})
}

// HandleLogout clears auth cookie
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	})

	respondJSON(w, http.StatusOK, map[string]string{"status": "logged out"})
}

// HandleVerify checks token validity
func HandleVerify(w http.ResponseWriter, r *http.Request) {
	tokenString := getToken(r)
	if tokenString == "" {
		respondError(w, http.StatusUnauthorized, "No token provided")
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		respondError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"valid": true})
}

// getToken extracts JWT from Authorization header or cookie
func getToken(r *http.Request) string {
	// Try Authorization header first
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			return authHeader[7:]
		}
	}

	// Try cookie
	if cookie, err := r.Cookie("auth_token"); err == nil {
		return cookie.Value
	}

	return ""
}

// VerifyToken middleware to protect routes
func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := getToken(r)
		if tokenString == "" {
			respondError(w, http.StatusUnauthorized, "No token provided")
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			respondError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Helper response functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
