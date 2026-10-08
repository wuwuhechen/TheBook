package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestHashPasswordAndVerifyPassword(t *testing.T) {
	hash, err := HashPassword("correct-password")
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}
	if hash == "correct-password" {
		t.Fatal("HashPassword returned plaintext password")
	}
	if !VerifyPassword(hash, "correct-password") {
		t.Fatal("VerifyPassword rejected the correct password")
	}
	if VerifyPassword(hash, "wrong-password") {
		t.Fatal("VerifyPassword accepted an incorrect password")
	}
}

func TestGenerateAndParseToken(t *testing.T) {
	InitJWT("auth-test-secret")

	token, err := GenerateToken(7, "alice", "user", "Alice", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken returned an error: %v", err)
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken returned an error: %v", err)
	}
	if claims.UserID != 7 || claims.UserName != "alice" || claims.Role != "user" || claims.Nickname != "Alice" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
	if claims.ExpiresAt == nil || !claims.ExpiresAt.After(time.Now()) {
		t.Fatal("token expiration was not set in the future")
	}
}

func TestParseTokenRejectsExpiredAndTamperedToken(t *testing.T) {
	InitJWT("auth-test-secret")

	expiredToken, err := GenerateToken(7, "alice", "user", "Alice", -time.Minute)
	if err != nil {
		t.Fatalf("GenerateToken returned an error: %v", err)
	}
	if _, err := ParseToken(expiredToken); err == nil {
		t.Fatal("ParseToken accepted an expired token")
	}

	validToken, err := GenerateToken(7, "alice", "user", "Alice", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken returned an error: %v", err)
	}
	if _, err := ParseToken(validToken + "tampered"); err == nil {
		t.Fatal("ParseToken accepted a tampered token")
	}
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	InitJWT("auth-test-secret")

	t.Run("rejects request without token", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthMiddleware())
		router.GET("/protected", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})

		request := httptest.NewRequest(http.MethodGet, "/protected", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		if response.Code != http.StatusFound {
			t.Fatalf("expected status %d, got %d", http.StatusFound, response.Code)
		}
		if location := response.Header().Get("Location"); location == "" {
			t.Fatal("expected redirect location for unauthenticated request")
		}
	})

	t.Run("adds claims to context for valid bearer token", func(t *testing.T) {
		token, err := GenerateToken(42, "bob", "user", "Bob", time.Hour)
		if err != nil {
			t.Fatalf("GenerateToken returned an error: %v", err)
		}

		router := gin.New()
		router.Use(AuthMiddleware())
		router.GET("/protected", func(c *gin.Context) {
			userID, exists := c.Get("userID")
			if !exists || userID != uint(42) {
				t.Fatalf("unexpected userID in context: %v", userID)
			}
			if username := c.GetString("username"); username != "bob" {
				t.Fatalf("unexpected username in context: %q", username)
			}
			c.Status(http.StatusNoContent)
		})

		request := httptest.NewRequest(http.MethodGet, "/protected", nil)
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		if response.Code != http.StatusNoContent {
			t.Fatalf("expected status %d, got %d", http.StatusNoContent, response.Code)
		}
	})
}
