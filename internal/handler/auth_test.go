package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"blog/internal/db"
)

var testAuthDB *db.DB

func init() {
	var err error
	testAuthDB, err = db.NewLocal(context.Background(), ":memory:")
	if err != nil {
		panic(err)
	}

	if err := testAuthDB.InitSchema(context.Background()); err != nil {
		panic(err)
	}

	// Set test environment
	os.Setenv("ADMIN_PASSWORD", "testpassword")
	os.Setenv("JWT_SECRET", "test-secret-key")
}

func TestHandleLogin_Success(t *testing.T) {
	os.Setenv("PASSWORD_HASH", "")  // Use raw password
	os.Setenv("ADMIN_PASSWORD", "testpass123")
	os.Setenv("JWT_SECRET", "secret")

	body := LoginRequest{
		Password: "testpass123",
	}

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	HandleLogin(w, req, testAuthDB)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp LoginResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Token == "" {
		t.Error("Token should not be empty")
	}
}

func TestHandleLogin_InvalidPassword(t *testing.T) {
	os.Setenv("PASSWORD_HASH", "")
	os.Setenv("ADMIN_PASSWORD", "correctpass")
	os.Setenv("JWT_SECRET", "secret")

	body := LoginRequest{
		Password: "wrongpass",
	}

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	HandleLogin(w, req, testAuthDB)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestHandleLogin_InvalidRequest(t *testing.T) {
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	HandleLogin(w, req, testAuthDB)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandleLogout(t *testing.T) {
	req := httptest.NewRequest("POST", "/auth/logout", nil)
	w := httptest.NewRecorder()

	HandleLogout(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check for cleared cookie
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "auth_token" && cookie.MaxAge == -1 {
			found = true
			break
		}
	}

	if !found {
		t.Error("Cookie should be cleared")
	}
}

func TestHandleVerify_ValidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")

	// First login to get token
	os.Setenv("ADMIN_PASSWORD", "verifytest")
	loginBody := LoginRequest{Password: "verifytest"}
	loginBodyBytes, _ := json.Marshal(loginBody)
	loginReq := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(loginBodyBytes))
	loginW := httptest.NewRecorder()

	HandleLogin(loginW, loginReq, testAuthDB)

	var loginResp LoginResponse
	json.NewDecoder(loginW.Body).Decode(&loginResp)

	// Now verify the token
	verifyReq := httptest.NewRequest("GET", "/auth/verify", nil)
	verifyReq.Header.Set("Authorization", "Bearer "+loginResp.Token)
	verifyW := httptest.NewRecorder()

	HandleVerify(verifyW, verifyReq)

	if verifyW.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", verifyW.Code)
	}

	var verifyResp map[string]bool
	json.NewDecoder(verifyW.Body).Decode(&verifyResp)

	if !verifyResp["valid"] {
		t.Error("Token should be valid")
	}
}

func TestHandleVerify_InvalidToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/auth/verify", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()

	HandleVerify(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestHandleVerify_NoToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/auth/verify", nil)
	w := httptest.NewRecorder()

	HandleVerify(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestTokenFromCookie(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("ADMIN_PASSWORD", "cookietest")

	// Login to get token
	loginBody := LoginRequest{Password: "cookietest"}
	loginBodyBytes, _ := json.Marshal(loginBody)
	loginReq := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(loginBodyBytes))
	loginW := httptest.NewRecorder()

	HandleLogin(loginW, loginReq, testAuthDB)

	// Extract cookie from response
	var cookie *http.Cookie
	for _, c := range loginW.Result().Cookies() {
		if c.Name == "auth_token" {
			cookie = c
			break
		}
	}

	if cookie == nil {
		t.Fatal("Cookie not found in response")
	}

	// Verify using cookie
	verifyReq := httptest.NewRequest("GET", "/auth/verify", nil)
	verifyReq.AddCookie(cookie)
	verifyW := httptest.NewRecorder()

	HandleVerify(verifyW, verifyReq)

	if verifyW.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", verifyW.Code)
	}
}
