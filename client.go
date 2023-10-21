package omada

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

//type OmadaClient interface {
//}

type omadaClient struct {
	httpClient   *http.Client
	omadaCId     string
	baseUrl      string
	clientId     string
	clientSecret string
}

func NewClient(baseUrl, omadaCId, clientId, clientSecret string, disableCertVerification bool) *omadaClient {
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: disableCertVerification},
	}}

	c := omadaClient{
		httpClient:   client,
		omadaCId:     omadaCId,
		baseUrl:      baseUrl,
		clientSecret: clientSecret,
		clientId:     clientId,
	}
	return &c
}

type AccessTokenResponse struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"msg"`
	Result    struct {
		AccessToken  string `json:"accessToken"`
		TokenType    string `json:"tokenType"`
		ExpiresIn    int    `json:"expiresIn"`
		RefreshToken string `json:"refreshToken"`
	} `json:"result"`
}

func (c *omadaClient) GetToken() (*AccessTokenResponse, error) {
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

func (c *omadaClient) GetRoleList() (string, error) {
	path := fmt.Sprintf("%s/openapi/v1/%s/roles", c.baseUrl, c.omadaCId)
	resp, err := c.httpClient.Get(path)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http response error code %d: %s", resp.StatusCode, resp.Status)
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return string(bytes), nil
}
