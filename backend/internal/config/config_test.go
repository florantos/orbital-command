package config_test

import (
	"os"
	"testing"

	"github.com/florantos/orbital-command/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				err := os.Unsetenv(key)
				require.NoError(t, err)
			}

			for key, val := range tt.envVars {
				t.Setenv(key, val)
			}

			cfg, err := config.Load()

			if tt.expectError {
				require.Error(t, err)

				assert.EqualError(t, err, tt.expectedErrMsg)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tt.envVars["APP_ENV"], cfg.Env)

			assert.Equal(t, tt.envVars["PORT"], cfg.Port)

			assert.Equal(t, tt.envVars["DATABASE_URL"], cfg.DatabaseURL)

			if tt.expectedLogLevel != "" {
				assert.Equal(t, tt.expectedLogLevel, cfg.LogLevel)

			}

		})
	}
}
