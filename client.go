package omada

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type OmadaClient struct {
	// For paginated requests
	PageSize int

	httpClient     *http.Client
	omadaCId       string
	baseUrl        string
	clientId       string
	clientSecret   string
	accessTokenCtx *accessTokenCtx
}

type tokenState int64

const (
	TokenStateUninitialised tokenState = 0
	TokenStateActive                   = 1
)

const defaultPageSize = 100

type accessTokenCtx struct {
	token      string
	tokenState tokenState
	//ttl 		int
	mu *sync.Mutex
}

func NewClient(baseUrl, omadaCId, clientId, clientSecret string, disableCertVerification bool) *OmadaClient {
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: disableCertVerification},
	}}

	c := OmadaClient{
		httpClient:   client,
		omadaCId:     omadaCId,
		baseUrl:      baseUrl,
		clientSecret: clientSecret,
		clientId:     clientId,
		accessTokenCtx: &accessTokenCtx{
			mu: &sync.Mutex{},
		},
		PageSize: defaultPageSize,
	}
	return &c
}

type EnvelopeResponse struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"msg"`
}

type AccessTokenResponse struct {
	EnvelopeResponse
	Result struct {
		AccessToken  string `json:"accessToken"`
		TokenType    string `json:"tokenType"`
		ExpiresIn    int    `json:"expiresIn"`
		RefreshToken string `json:"refreshToken"`
	} `json:"result"`
}

func (c *OmadaClient) GetToken() (*AccessTokenResponse, error) {
	path := fmt.Sprintf("%s/openapi/authorize/token?grant_type=client_credentials&client_id=%s&client_secret=%s", c.baseUrl, c.clientId, c.clientSecret)
	payload := map[string]string{
		"omadacId": c.omadaCId,
	}
	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	response, err := c.httpClient.Post(path, "application/json", bytes.NewReader(encodedPayload))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status %d: %s", response.StatusCode, response.Status)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	tokenResponse := &AccessTokenResponse{}
	err = json.Unmarshal(bodyBytes, tokenResponse)
	if err != nil {
		return nil, err
	}
	return tokenResponse, nil
}

func (c *OmadaClient) httpDoWrapped(request *http.Request, mapToJsonStructType interface{}) error {
	return c.internalHttpDoWithAuthContextAndJsonMarshalling(request, mapToJsonStructType, 1)
}

func (c *OmadaClient) internalHttpDoWithAuthContextAndJsonMarshalling(request *http.Request, mapToJsonStructType interface{}, tries int) error {
	if tries > 2 {
		return fmt.Errorf("could not perform request after refreshing token")
	}
	err := c.accessTokenCtx.initialiseAccessTokenIfNeeded(c)
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", fmt.Sprintf("AccessToken=%s", c.accessTokenCtx.getAccessToken()))
	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	// Should be a 200, even for errors
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response error: %d %s", response.StatusCode, response.Status)
	}

	// We need to check the body. An expired token still returns a 200, but the error is in the payload :(
	allBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// Check the response envelope for the expired token
	envelope := &EnvelopeResponse{}
	err = json.Unmarshal(allBytes, envelope)
	if err != nil {
		return err
	}
	if envelope.ErrorCode == -44112 {
		// Token expired, refresh the token and try again
		c.accessTokenCtx.resetAccessToken()
		return c.internalHttpDoWithAuthContextAndJsonMarshalling(request, mapToJsonStructType, tries+1)
	}

	// Finally, map to JSON
	err = json.Unmarshal(allBytes, mapToJsonStructType)
	if err != nil {
		return err
	}
	return nil

}

func (a *accessTokenCtx) getAccessToken() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.token
}

func (a *accessTokenCtx) resetAccessToken() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.tokenState = TokenStateUninitialised
	a.token = ""
}

func (a *accessTokenCtx) initialiseAccessTokenIfNeeded(c *OmadaClient) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.tokenState == TokenStateUninitialised {
		token, err := c.GetToken()
		if err != nil {
			return err
		}
		if token.ErrorCode != 0 {
			return fmt.Errorf("token error response: %d: %s", token.ErrorCode, token.Message)
		}
		a.tokenState = TokenStateActive
		a.token = token.Result.AccessToken
	}
	return nil
}
