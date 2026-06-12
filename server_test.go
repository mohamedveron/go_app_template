package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mohamedveron/go_app_template/internal/transport/http/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupMainServer creates a main server for testing (to avoid circular dependency)
func setupMainServer(t *testing.T) *testutils.TestServer {
	return setupMainServerWithOptions(t, testutils.TestServerOptions{})
}

func setupMainServerWithOptions(t *testing.T, opts testutils.TestServerOptions) *testutils.TestServer {
	// Use testutils to set up the environment and get all components
	testDir, authMiddleware := testutils.SetupTestEnvironment(t, opts)

	// Create main server with dependencies
	serverInstance := NewServer(authMiddleware, testDir, 0, "test")
	server := httptest.NewServer(serverInstance.GetRouter())

	return testutils.NewTestServer(server, testDir)
}

func TestHealthEndpoint(t *testing.T) {
	ts := setupMainServer(t)
	defer ts.Cleanup()

	t.Run("HealthCheck", func(t *testing.T) {
		resp, responseBody := ts.MakeRequest("GET", "/health", nil)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err := json.Unmarshal(responseBody, &response)
		require.NoError(t, err)
		assert.Equal(t, "ok", response["status"])
		assert.Equal(t, "test", response["version"])
	})
}
