package auth

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func testAuth(t *testing.T) (*Auth, *sql.DB) {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	a, err := New(db, "test-client-id", "test-client-secret", "http://localhost/auth/callback", false)
	if err != nil {
		t.Fatalf("new auth: %v", err)
	}
	return a, db
}

func TestAPIKeyAuth(t *testing.T) {
	a, db := testAuth(t)
	ctx := context.Background()

	// Create a test user.
	userID := "auth-test-user-1"
	db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, userID)
	_, err := db.ExecContext(ctx,
		`INSERT INTO users (id, google_id, email, name, kind) VALUES ($1, $2, $3, $4, 'human')`,
		userID, "google:auth-test", "authtest@test.com", "Auth Tester")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	t.Cleanup(func() {
		db.ExecContext(ctx, `DELETE FROM api_keys WHERE user_id = $1`, userID)
		db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, userID)
	})

	t.Run("create_and_auth", func(t *testing.T) {
		rawKey, err := a.createAPIKey(ctx, userID, "test-key", "")
		if err != nil {
			t.Fatalf("create key: %v", err)
		}
		if rawKey[:3] != "lv_" {
			t.Errorf("key should start with lv_, got %s", rawKey[:3])
		}

		// Auth with the key.
		user, err := a.userByAPIKey(ctx, rawKey)
		if err != nil {
			t.Fatalf("auth: %v", err)
		}
		if user.ID != userID {
			t.Errorf("user ID = %q, want %q", user.ID, userID)
		}
		if user.Name != "Auth Tester" {
			t.Errorf("name = %q, want %q", user.Name, "Auth Tester")
		}
		if user.Kind != "human" {
			t.Errorf("kind = %q, want %q", user.Kind, "human")
		}
	})

	t.Run("invalid_key", func(t *testing.T) {
		_, err := a.userByAPIKey(ctx, "lv_invalid_key_that_doesnt_exist")
		if err == nil {
			t.Error("should fail with invalid key")
		}
	})

	t.Run("agent_key", func(t *testing.T) {
		// Clean up any prior test agent.
		db.ExecContext(ctx, `DELETE FROM users WHERE name = 'AuthTestBot'`)

		rawKey, err := a.createAPIKey(ctx, userID, "agent-key", "AuthTestBot")
		if err != nil {
			t.Fatalf("create agent key: %v", err)
		}

		t.Cleanup(func() {
			db.ExecContext(ctx, `DELETE FROM users WHERE name = 'AuthTestBot'`)
		})

		// Auth should resolve to the agent, not the human.
		user, err := a.userByAPIKey(ctx, rawKey)
		if err != nil {
			t.Fatalf("auth: %v", err)
		}
		if user.Kind != "agent" {
			t.Errorf("kind = %q, want %q", user.Kind, "agent")
		}
		if user.Name != "AuthTestBot" {
			t.Errorf("name = %q, want %q", user.Name, "AuthTestBot")
		}
		// The user ID should be the agent's, not the sponsor's.
		if user.ID == userID {
			t.Errorf("should resolve to agent user ID, not sponsor ID")
		}
	})

	t.Run("bearer_header", func(t *testing.T) {
		rawKey, err := a.createAPIKey(ctx, userID, "bearer-test", "")
		if err != nil {
			t.Fatalf("create key: %v", err)
		}

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+rawKey)
		user := a.userFromBearer(req)
		if user == nil {
			t.Fatal("should resolve user from bearer token")
		}
		if user.ID != userID {
			t.Errorf("user ID = %q, want %q", user.ID, userID)
		}
	})

	t.Run("no_bearer", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		user := a.userFromBearer(req)
		if user != nil {
			t.Error("should return nil without bearer header")
		}
	})

	t.Run("list_keys", func(t *testing.T) {
		keys, err := a.ListAPIKeys(ctx, userID)
		if err != nil {
			t.Fatalf("list keys: %v", err)
		}
		if len(keys) < 1 {
			t.Errorf("should have at least 1 key, got %d", len(keys))
		}
	})
}

func TestRequireAuth(t *testing.T) {
	a, db := testAuth(t)
	ctx := context.Background()

	// Create user + key.
	userID := "auth-middleware-test"
	db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, userID)
	db.ExecContext(ctx,
		`INSERT INTO users (id, google_id, email, name, kind) VALUES ($1, $2, $3, $4, 'human')`,
		userID, "google:mw-test", "mwtest@test.com", "MW Tester")
	rawKey, _ := a.createAPIKey(ctx, userID, "mw-key", "")
	t.Cleanup(func() {
		db.ExecContext(ctx, `DELETE FROM api_keys WHERE user_id = $1`, userID)
		db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, userID)
	})

	handler := a.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		user := UserFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(user.Name))
	})

	t.Run("with_api_key", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+rawKey)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}
		if w.Body.String() != "MW Tester" {
			t.Errorf("body = %q, want %q", w.Body.String(), "MW Tester")
		}
	})

	t.Run("without_auth", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		// Should redirect to login.
		if w.Code != http.StatusSeeOther {
			t.Errorf("status = %d, want %d (redirect to login)", w.Code, http.StatusSeeOther)
		}
	})
}
