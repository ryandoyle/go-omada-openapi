package omada

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockValidTokenResponse(t *testing.T, w http.ResponseWriter, r *http.Request) {
	assert.Equal(t, "my-client-id", r.URL.Query().Get("client_id"))
	assert.Equal(t, "my-client-secret", r.URL.Query().Get("client_secret"))
	bytes, err := io.ReadAll(r.Body)
	assert.NoError(t, err)

	rawMap := &TestTokenRequest{}
	err = json.Unmarshal(bytes, rawMap)
	assert.NoError(t, err)
	assert.Equal(t, "my-cid", rawMap.OmadacId)

	_, err = w.Write([]byte(`{
			"errorCode": 0,
			"msg":  "hello",
			"result":  {
				"accessToken": "my-token",
				"tokenType": "Bearer",
				"expiresIn": 3600,
				"refreshToken": "my-refresh"
			}}`))
	assert.NoError(t, err)
}

func mockTokenExpiredResponse(t *testing.T, w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(`
		{
			"errorCode":-44112,
			"msg":"The Access Token has expired."
		}
	`))
	assert.NoError(t, err)
}

func TestOmadaClient_GetToken_ReturnsAValidToken(t *testing.T) {
	mockMux := http.NewServeMux()
	mockMux.HandleFunc("/openapi/authorize/token", func(w http.ResponseWriter, r *http.Request) { mockValidTokenResponse(t, w, r) })
	server := httptest.NewServer(mockMux)
	defer server.Close()

	c := NewClient(server.URL, "my-cid", "my-client-id", "my-client-secret", true)
	token, err := c.GetToken()

	assert.NoError(t, err)
	assert.Equal(t, "my-token", token.Result.AccessToken)
}

func TestNewClient_TokenRefreshingLogic_RefreshesWhenTheTokenIsExpired(t *testing.T) {
	mockMux := http.NewServeMux()
	mockMux.HandleFunc("/openapi/authorize/token", func(w http.ResponseWriter, r *http.Request) { mockValidTokenResponse(t, w, r) })
	handlerCalledTimes := 0
	mockMux.HandleFunc("/openapi/v1/my-cid/sites", func(w http.ResponseWriter, r *http.Request) {
		if handlerCalledTimes == 0 {
			mockTokenExpiredResponse(t, w, r)
		} else if handlerCalledTimes == 1 {
			_, err := w.Write([]byte(`{"errorCode": 0, "msg": "Success."}`))
			assert.NoError(t, err)
		} else {
			assert.Fail(t, "no mock response setup for this invocation")
		}
		handlerCalledTimes++
	})
	server := httptest.NewServer(mockMux)
	defer server.Close()

	c := NewClient(server.URL, "my-cid", "my-client-id", "my-client-secret", true)
	siteList, err := c.GetSiteList(1)

	assert.NoError(t, err)
	assert.Equal(t, siteList.ErrorCode, 0)
	assert.Equal(t, siteList.Message, "Success.")
}

func TestNewClient_TokenRefreshingLogic_WillPropagateTheErrorIfItCannotRefresh(t *testing.T) {
	mockMux := http.NewServeMux()
	mockMux.HandleFunc("/openapi/authorize/token", func(w http.ResponseWriter, r *http.Request) { mockValidTokenResponse(t, w, r) })
	mockMux.HandleFunc("/openapi/v1/my-cid/sites", func(w http.ResponseWriter, r *http.Request) { mockTokenExpiredResponse(t, w, r) })
	server := httptest.NewServer(mockMux)
	defer server.Close()

	c := NewClient(server.URL, "my-cid", "my-client-id", "my-client-secret", true)
	siteList, err := c.GetSiteList(1)

	assert.EqualError(t, err, "could not perform request after refreshing token")
	assert.Nil(t, siteList)
}

type TestTokenRequest struct {
	OmadacId string `json:"omadacId"`
}
