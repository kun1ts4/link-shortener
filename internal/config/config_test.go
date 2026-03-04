package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	path := makeTempConfig(t, `
    server:
        host: "127.0.0.1"
        port: 8888
    `)

	t.Setenv("POSTGRES_USER", "testuser")
	t.Setenv("POSTGRES_PASSWORD", "testpass")
	t.Setenv("POSTGRES_HOST", "testhost")
	t.Setenv("POSTGRES_PORT", "2345")
	t.Setenv("POSTGRES_DB", "testdb")

	cfg, err := LoadConfig(path)
	require.NoError(t, err)

	assert.Equal(t, 8888, cfg.Server.Port)
	assert.Equal(t, "127.0.0.1", cfg.Server.Host)

	expectedDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		"testuser", "testpass", "testhost", "2345", "testdb",
	)
	assert.Equal(t, expectedDSN, cfg.Database.DSN)
}

func TestLoadConfig_DefaultDSN(t *testing.T) {
	path := makeTempConfig(t, `
    server:
        host: "127.0.0.1"
        port: 8888
    `)

	for _, key := range []string{"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_DB"} {
		old, ok := os.LookupEnv(key)
		os.Unsetenv(key)
		if ok {
			t.Cleanup(func() { os.Setenv(key, old) })
		}
	}

	cfg, err := LoadConfig(path)
	require.NoError(t, err)

	assert.Equal(t, "postgres://postgres:postgres@localhost:5432/link_shortener?sslmode=disable", cfg.Database.DSN)
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("/not_config/config.yaml")
	require.Error(t, err)
}

func makeTempConfig(t *testing.T, content string) string {
	t.Helper()
	tmp, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	t.Cleanup(func() { os.Remove(tmp.Name()) })
	require.NoError(t, os.WriteFile(tmp.Name(), []byte(content), 0o644))
	return tmp.Name()
}
