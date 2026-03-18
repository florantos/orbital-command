package config_test

import (
	"os"
	"testing"

	"github.com/florantos/orbital-command/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_ValidConfig(t *testing.T) {
	envVars := map[string]string{
		"PORT":         "8080",
		"APP_ENV":      "development",
		"DATABASE_URL": "postgres://user:password@localhost:5432/orbital_command",
	}

	// defensive: unset vars to account for shell env vars
	for _, key := range []string{"PORT", "APP_ENV", "DATABASE_URL", "LOG_LEVEL"} {
		err := os.Unsetenv(key)
		require.NoError(t, err)
	}

	for key, val := range envVars {
		t.Setenv(key, val)
	}

	cfg, err := config.Load()

	require.NoError(t, err)

	assert.Equal(t, envVars["APP_ENV"], cfg.Env)
	assert.Equal(t, envVars["PORT"], cfg.Port)
	assert.Equal(t, envVars["DATABASE_URL"], cfg.DatabaseURL)
}

func TestLoad_MissingRequiredVars(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectedErr string
	}{
		{
			name: "missing APP_ENV returns error",
			envVars: map[string]string{
				"PORT":         "8080",
				"DATABASE_URL": "postgres://user:password@localhost:5432/orbital_command",
			},
			expectedErr: "cannot find env variable APP_ENV",
		},
		{
			name: "missing PORT returns error",
			envVars: map[string]string{
				"APP_ENV":      "development",
				"DATABASE_URL": "postgres://user:password@localhost:5432/orbital_command",
			},
			expectedErr: "cannot find env variable PORT",
		},
		{
			name: "missing DATABASE_URL returns error",
			envVars: map[string]string{
				"PORT":    "8080",
				"APP_ENV": "development",
			},
			expectedErr: "cannot find env variable DATABASE_URL",
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

			_, err := config.Load()

			assert.EqualError(t, err, tt.expectedErr)
		})
	}
}

func TestLoad_EmptyRequiredVars(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectedErr string
	}{
		{

			name: "empty APP_ENV returns error",
			envVars: map[string]string{
				"PORT":         "8080",
				"APP_ENV":      "",
				"DATABASE_URL": "postgres://user:password@localhost:5432/orbital_command",
			},
			expectedErr: "env variable APP_ENV is empty",
		},
		{
			name: "empty PORT returns error",
			envVars: map[string]string{
				"PORT":         "",
				"APP_ENV":      "development",
				"DATABASE_URL": "postgres://user:password@localhost:5432/orbital_command",
			},
			expectedErr: "env variable PORT is empty",
		},
		{
			name: "empty DATABASE_URL returns error",
			envVars: map[string]string{
				"PORT":         "8080",
				"APP_ENV":      "development",
				"DATABASE_URL": "",
			},
			expectedErr: "env variable DATABASE_URL is empty",
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

			_, err := config.Load()

			assert.EqualError(t, err, tt.expectedErr)
		})
	}

}

func TestLoad_DefaultLogLevel(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
	}{

		{
			name: "LOG_LEVEL defaults to info when not set",
			envVars: map[string]string{
				"PORT":         "8080",
				"APP_ENV":      "development",
				"DATABASE_URL": "postgres://user:password@localhost:5432/orbital_command",
			},
		},
		{
			name: "LOG_LEVEL defaults to info when empty",
			envVars: map[string]string{
				"PORT":         "8080",
				"APP_ENV":      "development",
				"DATABASE_URL": "postgres://user:password@localhost:5432/orbital_command",
				"LOG_LEVEL":    "",
			},
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

			require.NoError(t, err)

			assert.Equal(t, "info", cfg.LogLevel)

		})
	}
}
