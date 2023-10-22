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

func (c *omadaClient) GetRoleList() (*GetRoleListResponse, error) {
	path := fmt.Sprintf("%s/openapi/v1/%s/roles", c.baseUrl, c.omadaCId)
	request, err := http.NewRequest("GET", path, nil)

	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("AccessToken=%s", token.Result.AccessToken))
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http response error code %d: %s", resp.StatusCode, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)

	roleListResponse := &GetRoleListResponse{}
	err = json.Unmarshal(body, roleListResponse)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return roleListResponse, nil
}

type GetRoleListResponse struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"msg"`
	Result    []struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Type        int    `json:"type"`
		DefaultRole bool   `json:"defaultRole"`
		Source      int    `json:"source"`
		Privilege   struct {
			License          int `json:"license"`
			GlobalDashboard  int `json:"globalDashboard"`
			Dashboard        int `json:"dashboard"`
			Devices          int `json:"devices"`
			Adopt            int `json:"adopt"`
			GlobalLog        int `json:"globalLog"`
			Log              int `json:"log"`
			LicenseBind      int `json:"licenseBind"`
			Users            int `json:"users"`
			Roles            int `json:"roles"`
			SamlUsers        int `json:"samlUsers"`
			SamlRoles        int `json:"samlRoles"`
			SamlSSOs         int `json:"samlSsos"`
			GlobalSetting    int `json:"globalSetting"`
			ExportData       int `json:"exportData"`
			GlobalExportData int `json:"globalExportData"`
			ExportGlobalLog  int `json:"exportGlobalLog"`
			Hotspot          int `json:"hotspot"`
			Statics          int `json:"statics"`
			Map              int `json:"map"`
			Clients          int `json:"clients"`
			Insight          int `json:"insight"`
			Report           int `json:"report"`
			Network          int `json:"network"`
			DeviceAccount    int `json:"deviceAccount"`
			Anomaly          int `json:"anomaly"`
			Analyze          int `json:"analyze"`
			SiteAnalyze      int `json:"siteAnalyze"`
		} `json:"privilege"`
	} `json:"result"`
}

//func (c *omadaClient) GetSiteList() (string, error) {
//	path := fmt.Sprintf("%s/openapi/v1/%s/sites?pageSize=100&page=1", c.baseUrl, c.omadaCId)
//	request, err := http.NewRequest("GET", path, nil)
//
//	token, err := c.GetToken()
//	if err != nil {
//		return "", err
//	}
//	request.Header.Set("Authorization", fmt.Sprintf("AccessToken=%s", token.Result.AccessToken))
//
//	resp, err := c.httpClient.Do(request)
//	if err != nil {
//		//log.Fatal(io.ReadAll(resp.Body))
//		return "", err
//	}
//	if resp.StatusCode != http.StatusOK {
//		body, err := io.ReadAll(resp.Body)
//		if err != nil {
//			return "", err
//		}
//		return "", fmt.Errorf("http response error code %d: %s %s", resp.StatusCode, resp.Status, body)
//	}
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return "", err
//	}
//	defer resp.Body.Close()
//	return string(body), nil
//}
