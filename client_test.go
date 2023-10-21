package omada

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOmadaClient_GetToken_ReturnsAValidToken(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/openapi/authorize/token", r.URL.Path)
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

	}))

	c := NewClient(server.URL, "my-cid", "my-client-id", "my-client-secret", true)
	token, err := c.GetToken()

	assert.NoError(t, err)
	assert.Equal(t, "my-token", token.Result.AccessToken)

	server.Close()
}

type TestTokenRequest struct {
	OmadacId string `json:"omadacId"`
}
