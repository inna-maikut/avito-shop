//go:build integration

package integration

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/inna-maikut/avito-shop/internal/api"
	"github.com/inna-maikut/avito-shop/internal/infrastructure/config"
)

const (
	usernameLen = 32
	password    = "password"
)

var noOut *error

func setUp() {
	_ = config.Load()
}

func makeUsername(t *testing.T) string {
	b := make([]byte, usernameLen)
	_, err := rand.Read(b)
	require.NoError(t, err)

	return base64.URLEncoding.EncodeToString(b)[:usernameLen]
}

func makeUserToken(t *testing.T, username string) string {
	resp := apiPost(t, "/api/auth", "", api.AuthRequest{
		Username: username,
		Password: password,
	})
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		require.Equal(t, http.StatusOK, resp.StatusCode, "body: "+string(body))
	}
	out := parseJSON[api.AuthResponse](t, resp)
	require.NotEmpty(t, out.Token)

	return *out.Token
}

func getInfo(t *testing.T, token string) api.InfoResponse {
	resp := apiGet(t, "/api/info", token)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	out := parseJSON[api.InfoResponse](t, resp)
	return out
}

// apiGet path should start with slash
func apiGet(t *testing.T, path, token string) *http.Response {
	t.Helper()

	url := "http://localhost:" + os.Getenv("SERVER_PORT") + path
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	if token != "" {
		req.Header.Set("Authorization", token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

// apiPost path should start with slash
func apiPost[In any](t *testing.T, path, token string, in In) *http.Response {
	t.Helper()

	inStr, err := json.Marshal(in)
	require.NoError(t, err)

	url := "http://localhost:" + os.Getenv("SERVER_PORT") + path
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(inStr))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

func parseJSON[Out any](t *testing.T, resp *http.Response) Out {
	var out Out

	bodyBytes, err := io.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()
	require.NoError(t, err)

	err = json.Unmarshal(bodyBytes, &out)
	require.NoError(t, err)

	return out
}

func assertResponseError(t *testing.T, resp *http.Response, statusCode int, errorText string) {
	require.Equal(t, statusCode, resp.StatusCode)

	out := parseJSON[api.ErrorResponse](t, resp)

	require.Equal(t, errorText, *out.Errors)
}

func assertResponsePlainError(t *testing.T, resp *http.Response, statusCode int, errorText string) {
	require.Equal(t, statusCode, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()
	require.NoError(t, err)

	require.Equal(t, errorText, string(bodyBytes))
}
