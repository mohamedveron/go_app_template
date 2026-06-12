package testutils

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mohamedveron/go_app_template/internal/transport/http/middleware"
	"github.com/stretchr/testify/require"
)

// TestServer represents a test server instance with shared test utilities
type TestServer struct {
	server  *httptest.Server
	testDir string
}

// TestServerOptions configures the test server setup
type TestServerOptions struct {
	WorkspaceDir    string // If provided, uses this directory instead of creating a temp one
	UseRandomTokens bool   // If true, generates random tokens; otherwise uses static test tokens
}

// SetupTestEnvironment creates the test environment and returns the components for server creation
func SetupTestEnvironment(t *testing.T, opts TestServerOptions) (string, *middleware.AuthMiddleware) {
	gin.SetMode(gin.TestMode)

	// Create or use workspace directory
	var testDir string
	var err error
	if opts.WorkspaceDir != "" {
		testDir = opts.WorkspaceDir
	} else {
		testDir, err = os.MkdirTemp("", "mohamedveron-go_app_template-test-*")
		require.NoError(t, err)
	}

	// Initialize services for testing
	authConfig := middleware.AuthConfig{
		Workspace: testDir,
	}
	authMiddleware, err := middleware.NewAuthMiddleware(authConfig)
	require.NoError(t, err)

	require.NoError(t, err)

	return testDir, authMiddleware
}

// GenerateRandomToken generates a random token for testing (exported for use in main server tests)
func GenerateRandomToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func generateRandomToken() string {
	return GenerateRandomToken()
}

// NewTestServer creates a TestServer instance with the given parameters (for main server tests)
func NewTestServer(server *httptest.Server, testDir string) *TestServer {
	return &TestServer{
		server:  server,
		testDir: testDir,
	}
}

// Cleanup cleans up test server resources
func (ts *TestServer) Cleanup() {
	ts.server.Close()
	_ = os.RemoveAll(ts.testDir)
}

// TestFilePath returns a file path within the test directory
func (ts *TestServer) TestFilePath(relativePath string) string {
	return filepath.Join(ts.testDir, relativePath)
}

// MakeRequest makes an HTTP request without authentication
func (ts *TestServer) MakeRequest(method, path string, body []byte) (*http.Response, []byte) {
	return ts.MakeRequestWithToken(method, path, body, "")
}

// MakeRequestWithToken makes an HTTP request with optional authentication token
func (ts *TestServer) MakeRequestWithToken(method, path string, body []byte, token string) (*http.Response, []byte) {
	url := ts.server.URL + path
	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	require.NoError(nil, err)

	// Add authentication token if provided
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(nil, err)

	defer func() { _ = resp.Body.Close() }()
	responseBody, err := io.ReadAll(resp.Body)
	require.NoError(nil, err)

	return resp, responseBody
}

func (ts *TestServer) GetTestDir() string {
	return ts.testDir
}

func (ts *TestServer) GetServerURL() string {
	return ts.server.URL
}
