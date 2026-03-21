// Package auth provides Google OAuth2 authentication and session management.
package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// User represents an authenticated user.
type User struct {
	ID      string
	Email   string
	Name    string
	Picture string
}

type contextKey struct{}

// UserFromContext returns the authenticated user, or nil.
func UserFromContext(ctx context.Context) *User {
	u, _ := ctx.Value(contextKey{}).(*User)
	return u
}

// ContextWithUser stores a user in the context.
func ContextWithUser(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, contextKey{}, u)
}

// Auth handles Google OAuth2 and session management.
type Auth struct {
	db     *sql.DB
	oauth  *oauth2.Config
	secure bool
}

// New creates an Auth service backed by the given database.
func New(db *sql.DB, clientID, clientSecret, redirectURL string, secure bool) (*Auth, error) {
	a := &Auth{
		db: db,
		oauth: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
		secure: secure,
	}
	if err := a.migrate(); err != nil {
		return nil, fmt.Errorf("auth migrate: %w", err)
	}
	return a, nil
}

func (a *Auth) migrate() error {
	_, err := a.db.Exec(`
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    google_id TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL DEFAULT '',
    picture TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
`)
	return err
}

// Register adds auth routes to the mux.
func (a *Auth) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /auth/login", a.handleLogin)
	mux.HandleFunc("GET /auth/callback", a.handleCallback)
	mux.HandleFunc("POST /auth/logout", a.handleLogout)
}

// RequireAuth wraps a handler to require authentication.
// Redirects to /auth/login if no valid session exists.
func (a *Auth) RequireAuth(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		user, err := a.userBySession(r.Context(), cookie.Value)
		if err != nil {
			a.clearCookie(w)
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		ctx := ContextWithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ────────────────────────────────────────────────────────────────────
// Handlers
// ────────────────────────────────────────────────────────────────────

func (a *Auth) handleLogin(w http.ResponseWriter, r *http.Request) {
	state := newID()
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/auth",
		MaxAge:   300,
		HttpOnly: true,
		Secure:   a.secure,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, a.oauth.AuthCodeURL(state), http.StatusSeeOther)
}

func (a *Auth) handleCallback(w http.ResponseWriter, r *http.Request) {
	// Verify CSRF state.
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "oauth_state", Path: "/auth", MaxAge: -1})

	// Exchange code for token.
	token, err := a.oauth.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		log.Printf("auth: oauth exchange: %v", err)
		http.Error(w, "authentication failed", http.StatusInternalServerError)
		return
	}

	// Fetch Google user info.
	client := a.oauth.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("auth: fetch userinfo: %v", err)
		http.Error(w, "authentication failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var info struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Printf("auth: decode userinfo: %v", err)
		http.Error(w, "authentication failed", http.StatusInternalServerError)
		return
	}

	// Upsert user.
	user, err := a.upsertUser(r.Context(), info.ID, info.Email, info.Name, info.Picture)
	if err != nil {
		log.Printf("auth: upsert user: %v", err)
		http.Error(w, "authentication failed", http.StatusInternalServerError)
		return
	}

	// Create session (30 days).
	sessionID := newID()
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	if _, err := a.db.ExecContext(r.Context(),
		`INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)`,
		sessionID, user.ID, expiresAt,
	); err != nil {
		log.Printf("auth: create session: %v", err)
		http.Error(w, "authentication failed", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   a.secure,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/work", http.StatusSeeOther)
}

func (a *Auth) handleLogout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("session"); err == nil {
		a.db.ExecContext(r.Context(), `DELETE FROM sessions WHERE id = $1`, cookie.Value)
	}
	a.clearCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ────────────────────────────────────────────────────────────────────
// Internal
// ────────────────────────────────────────────────────────────────────

func (a *Auth) upsertUser(ctx context.Context, googleID, email, name, picture string) (*User, error) {
	var u User
	err := a.db.QueryRowContext(ctx, `
		INSERT INTO users (id, google_id, email, name, picture)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (google_id) DO UPDATE SET
			email = EXCLUDED.email,
			name = EXCLUDED.name,
			picture = EXCLUDED.picture
		RETURNING id, email, name, picture`,
		newID(), googleID, email, name, picture,
	).Scan(&u.ID, &u.Email, &u.Name, &u.Picture)
	return &u, err
}

func (a *Auth) userBySession(ctx context.Context, sessionID string) (*User, error) {
	var u User
	err := a.db.QueryRowContext(ctx, `
		SELECT u.id, u.email, u.name, u.picture
		FROM users u
		JOIN sessions s ON s.user_id = u.id
		WHERE s.id = $1 AND s.expires_at > NOW()`,
		sessionID,
	).Scan(&u.ID, &u.Email, &u.Name, &u.Picture)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (a *Auth) clearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand: " + err.Error())
	}
	return hex.EncodeToString(b)
}
