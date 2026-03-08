package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTokenStoreValidation(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name           string
		token          TokenStore
		accessValid    bool
		refreshValid   bool
	}{
		{
			name: "both valid",
			token: TokenStore{
				AccessToken:      "access",
				RefreshToken:     "refresh",
				ExpiresAt:        now.Add(1 * time.Hour),
				RefreshExpiresAt: now.Add(7 * 24 * time.Hour),
			},
			accessValid:  true,
			refreshValid: true,
		},
		{
			name: "access expired, refresh valid",
			token: TokenStore{
				AccessToken:      "access",
				RefreshToken:     "refresh",
				ExpiresAt:        now.Add(-1 * time.Hour),
				RefreshExpiresAt: now.Add(7 * 24 * time.Hour),
			},
			accessValid:  false,
			refreshValid: true,
		},
		{
			name: "both expired",
			token: TokenStore{
				AccessToken:      "access",
				RefreshToken:     "refresh",
				ExpiresAt:        now.Add(-1 * time.Hour),
				RefreshExpiresAt: now.Add(-1 * time.Hour),
			},
			accessValid:  false,
			refreshValid: false,
		},
		{
			name: "access within 60s buffer",
			token: TokenStore{
				AccessToken: "access",
				ExpiresAt:   now.Add(30 * time.Second),
			},
			accessValid:  false,
			refreshValid: false,
		},
		{
			name:         "empty token",
			token:        TokenStore{},
			accessValid:  false,
			refreshValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.token.IsAccessTokenValid(); got != tt.accessValid {
				t.Errorf("IsAccessTokenValid() = %v, want %v", got, tt.accessValid)
			}
			if got := tt.token.IsRefreshTokenValid(); got != tt.refreshValid {
				t.Errorf("IsRefreshTokenValid() = %v, want %v", got, tt.refreshValid)
			}
		})
	}
}

func TestMaskToken(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"u-abcdefghijklmnopqrstuvwxyz", "u-abcd...uvwxyz"},
		{"short", "***"},
		{"exactly12ch", "***"},
		{"1234567890123", "123456...890123"},
	}

	for _, tt := range tests {
		if got := MaskToken(tt.input); got != tt.expected {
			t.Errorf("MaskToken(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestSaveAndLoadToken(t *testing.T) {
	tmpDir := t.TempDir()
	tokenFile := filepath.Join(tmpDir, "token.json")

	// Override TokenPath for testing
	original := tokenPathFunc
	tokenPathFunc = func() (string, error) { return tokenFile, nil }
	defer func() { tokenPathFunc = original }()

	now := time.Now().Truncate(time.Second)
	token := &TokenStore{
		AccessToken:      "test-access-token",
		RefreshToken:     "test-refresh-token",
		TokenType:        "Bearer",
		ExpiresAt:        now.Add(2 * time.Hour),
		RefreshExpiresAt: now.Add(7 * 24 * time.Hour),
		Scope:            "offline_access",
	}

	// Save
	if err := SaveToken(token); err != nil {
		t.Fatalf("SaveToken() error: %v", err)
	}

	// Check file permissions
	info, err := os.Stat(tokenFile)
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("file permissions = %o, want 0600", perm)
	}

	// Load
	loaded, err := LoadToken()
	if err != nil {
		t.Fatalf("LoadToken() error: %v", err)
	}
	if loaded.AccessToken != token.AccessToken {
		t.Errorf("AccessToken = %q, want %q", loaded.AccessToken, token.AccessToken)
	}
	if loaded.RefreshToken != token.RefreshToken {
		t.Errorf("RefreshToken = %q, want %q", loaded.RefreshToken, token.RefreshToken)
	}

	// Delete
	if err := DeleteToken(); err != nil {
		t.Fatalf("DeleteToken() error: %v", err)
	}
	loaded, err = LoadToken()
	if err != nil {
		t.Fatalf("LoadToken() after delete error: %v", err)
	}
	if loaded != nil {
		t.Error("expected nil after delete")
	}
}

func TestLoadTokenNotExist(t *testing.T) {
	tokenPathFunc = func() (string, error) { return filepath.Join(t.TempDir(), "nonexistent.json"), nil }
	defer func() { tokenPathFunc = originalTokenPath }()

	token, err := LoadToken()
	if err != nil {
		t.Fatalf("LoadToken() error: %v", err)
	}
	if token != nil {
		t.Error("expected nil for nonexistent file")
	}
}
