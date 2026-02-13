package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/creatrid/creatrid/internal/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const testSecret = "test-secret-key-for-unit-tests"

func TestRequireAuth_MissingCookie(t *testing.T) {
	jwtSvc := auth.NewJWTService(testSecret)

	handler := RequireAuth(jwtSvc, nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called when cookie is missing")
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if body["error"] != "Not authenticated" {
		t.Errorf("expected error 'Not authenticated', got %q", body["error"])
	}
}

func TestRequireAuth_InvalidJWT(t *testing.T) {
	jwtSvc := auth.NewJWTService(testSecret)

	handler := RequireAuth(jwtSvc, nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called with invalid JWT")
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "this-is-not-a-valid-jwt"})
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if body["error"] != "Not authenticated" {
		t.Errorf("expected error 'Not authenticated', got %q", body["error"])
	}
}

func TestRequireAuth_ExpiredJWT(t *testing.T) {
	jwtSvc := auth.NewJWTService(testSecret)

	// Create an expired token manually
	claims := jwt.RegisteredClaims{
		Subject:   "user-123",
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredToken, err := token.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("failed to create expired token: %v", err)
	}

	handler := RequireAuth(jwtSvc, nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called with expired JWT")
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: expiredToken})
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if body["error"] != "Not authenticated" {
		t.Errorf("expected error 'Not authenticated', got %q", body["error"])
	}
}

func TestRequireAuth_WrongSigningKey(t *testing.T) {
	jwtSvc := auth.NewJWTService(testSecret)

	// Create a token signed with a different secret
	wrongSvc := auth.NewJWTService("wrong-secret-key")
	wrongToken, err := wrongSvc.Generate("user-123")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	handler := RequireAuth(jwtSvc, nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called with wrongly-signed JWT")
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: wrongToken})
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func TestRequireAuth_ValidJWT_NilStore(t *testing.T) {
	// When the JWT is valid but store is nil, the middleware will panic or return 401.
	// With a nil store, FindUserByID call will fail, resulting in 401.
	// This test verifies that a valid JWT alone is not sufficient without a valid store lookup.
	jwtSvc := auth.NewJWTService(testSecret)

	validToken, err := jwtSvc.Generate("user-123")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// The middleware will attempt to call st.FindUserByID on a nil store,
	// which will cause a panic. We recover from it to verify the JWT was valid
	// but the store dependency is required.
	handler := RequireAuth(jwtSvc, nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: validToken})
	rr := httptest.NewRecorder()

	panicked := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		handler.ServeHTTP(rr, req)
	}()

	// A valid JWT with nil store should panic when trying to call FindUserByID
	// This proves the JWT was accepted and the middleware proceeded past validation
	if !panicked {
		t.Error("expected panic when store is nil but JWT is valid (proves JWT validation passed)")
	}
}

func TestJWTService_GenerateAndValidate(t *testing.T) {
	jwtSvc := auth.NewJWTService(testSecret)

	t.Run("generate and validate round trip", func(t *testing.T) {
		token, err := jwtSvc.Generate("user-abc-123")
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		subject, err := jwtSvc.Validate(token)
		if err != nil {
			t.Fatalf("failed to validate token: %v", err)
		}

		if subject != "user-abc-123" {
			t.Errorf("expected subject 'user-abc-123', got %q", subject)
		}
	})

	t.Run("validate returns error for garbage input", func(t *testing.T) {
		_, err := jwtSvc.Validate("not-a-jwt")
		if err == nil {
			t.Error("expected error for invalid token")
		}
	})

	t.Run("validate returns error for expired token", func(t *testing.T) {
		claims := jwt.RegisteredClaims{
			Subject:   "user-expired",
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-48 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-24 * time.Hour)),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, _ := token.SignedString([]byte(testSecret))

		_, err := jwtSvc.Validate(signed)
		if err == nil {
			t.Error("expected error for expired token")
		}
	})

	t.Run("validate returns error for wrong signing key", func(t *testing.T) {
		otherSvc := auth.NewJWTService("other-secret")
		token, _ := otherSvc.Generate("user-123")

		_, err := jwtSvc.Validate(token)
		if err == nil {
			t.Error("expected error for token signed with different key")
		}
	})
}

// Additional tests using testify/assert

func TestRequireAuth_MissingCookie_Assert(t *testing.T) {
	jwtSvc := auth.NewJWTService(testSecret)
	handler := RequireAuth(jwtSvc, nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	var body map[string]string
	json.NewDecoder(rr.Body).Decode(&body)
	assert.Equal(t, "Not authenticated", body["error"])
}

func TestRequireAuth_InvalidJWT_Assert(t *testing.T) {
	jwtSvc := auth.NewJWTService(testSecret)
	handler := RequireAuth(jwtSvc, nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "invalid-jwt-garbage"})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestRequireAuth_ExpiredJWT_Assert(t *testing.T) {
	jwtSvc := auth.NewJWTService(testSecret)

	claims := jwt.RegisteredClaims{
		Subject:   "user-expired",
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredToken, err := token.SignedString([]byte(testSecret))
	assert.NoError(t, err)

	handler := RequireAuth(jwtSvc, nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: expiredToken})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestRequireAuth_ValidJWT_PassesThrough_Assert(t *testing.T) {
	jwtSvc := auth.NewJWTService(testSecret)

	validToken, err := jwtSvc.Generate("user-123")
	assert.NoError(t, err)

	// A valid JWT with nil store will panic when the middleware tries to call
	// st.FindUserByID. The panic proves the JWT was accepted and the middleware
	// progressed past the token validation step.
	handler := RequireAuth(jwtSvc, nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: validToken})
	rr := httptest.NewRecorder()

	panicked := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		handler.ServeHTTP(rr, req)
	}()

	assert.True(t, panicked, "valid JWT should cause the middleware to proceed to store lookup, which panics with nil store")
}
