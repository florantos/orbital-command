package config_test

import (
	"os"
	"testing"

	"github.com/florantos/orbital-command/internal/config"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name             string
		envVars          map[string]string
		expectError      bool
		expectedErrMsg   string
		expectedLogLevel string
	}{
		{
			name: "valid config loads successfully",
			envVars: map[string]string{
				"PORT":         "8080",
				"APP_ENV":      "development",
				"DATABASE_URL": "postgres://user:password@localhost:5432/orbital_command",
			},
			expectError: false,
		},
		{
			name: "missing ENV_VAR returns error",
			envVars: map[string]string{
				"PORT":         "8080",
				"DATABASE_URL": "postgres://user:password@localhost:5432/orbital_command",
			},
			expectError:    true,
			expectedErrMsg: "cannot find env variable APP_ENV",
		},
		{
			name: "empty ENV_VAR returns error",
			envVars: map[string]string{
				"PORT":         "8080",
				"APP_ENV":      "",
				"DATABASE_URL": "postgres://user:password@localhost:5432/orbital_command",
			},
			expectError:    true,
			expectedErrMsg: "env variable APP_ENV is empty",
		},
		{
			name: "missing PORT returns error",
			envVars: map[string]string{
				"APP_ENV":      "development",
				"DATABASE_URL": "postgres://user:password@localhost:5432/orbital_command",
			},
			expectError:    true,
			expectedErrMsg: "cannot find env variable PORT",
		},
		{
			name: "empty PORT returns error",
			envVars: map[string]string{
				"PORT":         "",
				"APP_ENV":      "development",
				"DATABASE_URL": "postgres://user:password@localhost:5432/orbital_command",
			},
			expectError:    true,
			expectedErrMsg: "env variable PORT is empty",
		},
		{
			name: "missing DATABASE_URL returns error",
			envVars: map[string]string{
				"PORT":    "8080",
				"APP_ENV": "development",
			},
			expectError:    true,
			expectedErrMsg: "cannot find env variable DATABASE_URL",
		},
		{
			name: "empty DATABASE_URL returns error",
			envVars: map[string]string{
				"PORT":         "8080",
				"APP_ENV":      "development",
				"DATABASE_URL": "",
			},
			expectError:    true,
			expectedErrMsg: "env variable DATABASE_URL is empty",
		},
		{
			name: "LOG_LEVEL defaults to info when not set",
			envVars: map[string]string{
				"PORT":         "8080",
				"APP_ENV":      "development",
				"DATABASE_URL": "postgres://user:password@localhost:5432/orbital_command",
			},
			expectError:      false,
			expectedLogLevel: "info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// defensive: unset vars to account for shell env vars
			for _, key := range []string{"PORT", "APP_ENV", "DATABASE_URL", "LOG_LEVEL"} {
				if err := os.Unsetenv(key); err != nil {
					t.Fatalf("failed to unset env var %s: %v", key, err)
				}
			}

			for key, val := range tt.envVars {
				t.Setenv(key, val)
			}

			cfg, err := config.Load()

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
					return
				}
				if tt.expectedErrMsg != "" && err.Error() != tt.expectedErrMsg {
					t.Errorf("expected error message %q, got %q", tt.expectedErrMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cfg.Env != tt.envVars["APP_ENV"] {
				t.Errorf("expected app env %s, got %s", tt.envVars["APP_ENV"], cfg.Env)
			}
			if cfg.Port != tt.envVars["PORT"] {
				t.Errorf("expected port %s, got %s", tt.envVars["PORT"], cfg.Port)
			}
			if cfg.DatabaseURL != tt.envVars["DATABASE_URL"] {
				t.Errorf("expected database url %s, got %s", tt.envVars["DATABASE_URL"], cfg.DatabaseURL)
			}

			if tt.expectedLogLevel != "" && cfg.LogLevel != tt.expectedLogLevel {
				t.Errorf("expected log level %s, got %s", tt.expectedLogLevel, cfg.LogLevel)
			}

		})
	}
}
