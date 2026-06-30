package auth_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lwmacct/260630-go-hsr-auth/pkg/auth"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func TestModulePasswordSessionFlow(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	defer func() { _ = db.Close() }()
	if err := auth.ApplySchema(ctx, db); err != nil {
		t.Fatal(err)
	}

	module, err := auth.New(auth.Options{
		DB: db,
		Config: auth.Config{
			Local: auth.LocalConfig{
				LoginEnabled:        true,
				RegistrationEnabled: true,
			},
			ChallengeProvider: passChallengeProvider{},
		},
		SessionTTL: time.Hour,
	})
	if err != nil {
		t.Fatal(err)
	}
	handler := withRequestContext(module.Handler())

	register := httptest.NewRecorder()
	handler.ServeHTTP(register, jsonRequest(http.MethodPost, "/auth/password/register", map[string]any{
		"username": "Alice",
		"password": "correct-password",
		"challenge": map[string]any{
			"provider": "pass",
		},
	}))
	if register.Code != http.StatusCreated {
		t.Fatalf("register status = %d, body = %s", register.Code, register.Body.String())
	}
	sessionCookie := register.Result().Cookies()[0]
	if sessionCookie.Name != "web_session" || sessionCookie.Value == "" {
		t.Fatalf("unexpected session cookie: %#v", sessionCookie)
	}

	me := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.AddCookie(sessionCookie)
	handler.ServeHTTP(me, req)
	if me.Code != http.StatusOK {
		t.Fatalf("me status = %d, body = %s", me.Code, me.Body.String())
	}
	if !strings.Contains(me.Body.String(), `"username":"alice"`) {
		t.Fatalf("expected normalized user in response, got %s", me.Body.String())
	}
}

func openTestDB(t *testing.T) *bun.DB {
	t.Helper()

	sqlDB, err := sql.Open(sqliteshim.ShimName, "file:auth-module-test?mode=memory&cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	return bun.NewDB(sqlDB, sqlitedialect.New())
}

func jsonRequest(method string, target string, body any) *http.Request {
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(method, target, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func withRequestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := auth.SessionRequest{
			IP:         "127.0.0.1",
			Host:       r.Host,
			UserAgent:  r.UserAgent(),
			Method:     r.Method,
			Path:       r.URL.Path,
			RemoteAddr: r.RemoteAddr,
		}
		next.ServeHTTP(w, r.WithContext(auth.ContextWithRequest(r.Context(), request)))
	})
}

type passChallengeProvider struct{}

func (passChallengeProvider) Name() string {
	return "pass"
}

func (passChallengeProvider) PublicConfig() auth.ChallengePublicConfig {
	return auth.ChallengePublicConfig{Provider: "pass"}
}

func (passChallengeProvider) Create(context.Context, auth.ChallengeInput) (*auth.Challenge, error) {
	return &auth.Challenge{Provider: "pass", ExpiresAt: time.Now().Add(time.Minute)}, nil
}

func (passChallengeProvider) Verify(context.Context, auth.ChallengeAnswer, auth.ChallengeInput) error {
	return nil
}
