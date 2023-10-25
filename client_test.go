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

func TestOmadaClient_GetRoleList_ReturnsARoleList(t *testing.T) {
	mockMux := http.NewServeMux()
	mockMux.HandleFunc("/openapi/authorize/token", func(w http.ResponseWriter, r *http.Request) { mockValidTokenResponse(t, w, r) })
	mockMux.HandleFunc("/openapi/v1/my-cid/roles", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "AccessToken=my-token", r.Header.Get("Authorization"))
		_, err := w.Write([]byte(`{
		  "errorCode": 0,
		  "msg": "Success.",
		  "result": [
			{
			  "id": "master_admin_id",
			  "name": "Main Administrator",
			  "type": 0,
			  "defaultRole": true,
			  "source": 0,
			  "privilege": {
				"license": 2,
				"globalDashboard": 2,
				"dashboard": 2,
				"devices": 2,
				"adopt": 2,
				"globalLog": 2,
				"log": 2,
				"licenseBind": 2,
				"users": 2,
				"roles": 2,
				"samlUsers": 2,
				"samlRoles": 2,
				"samlSsos": 2,
				"globalSetting": 2,
				"exportData": 2,
				"globalExportData": 2,
				"exportGlobalLog": 2,
				"hotspot": 2,
				"statics": 2,
				"map": 2,
				"clients": 2,
				"insight": 2,
				"report": 2,
				"network": 2,
				"deviceAccount": 2,
				"anomaly": 2,
				"analyze": 2,
				"siteAnalyze": 2
			  }
			}]}`))
		assert.NoError(t, err)
	})
	server := httptest.NewServer(mockMux)
	defer server.Close()

	c := NewClient(server.URL, "my-cid", "my-client-id", "my-client-secret", true)
	roleList, err := c.GetRoleList()

	assert.NoError(t, err)
	assert.Equal(t, roleList.ErrorCode, 0)
	assert.Equal(t, roleList.Message, "Success.")
	assert.Equal(t, roleList.Result[0].Id, "master_admin_id")
	assert.Equal(t, roleList.Result[0].Name, "Main Administrator")
	assert.Equal(t, roleList.Result[0].Type, 0)
	assert.Equal(t, roleList.Result[0].DefaultRole, true)
	assert.Equal(t, roleList.Result[0].Source, 0)
	assert.Equal(t, roleList.Result[0].Privilege.License, 2)
	assert.Equal(t, roleList.Result[0].Privilege.GlobalDashboard, 2)
	assert.Equal(t, roleList.Result[0].Privilege.Dashboard, 2)
	assert.Equal(t, roleList.Result[0].Privilege.Devices, 2)
	assert.Equal(t, roleList.Result[0].Privilege.Adopt, 2)
	assert.Equal(t, roleList.Result[0].Privilege.GlobalLog, 2)
	assert.Equal(t, roleList.Result[0].Privilege.Log, 2)
	assert.Equal(t, roleList.Result[0].Privilege.LicenseBind, 2)
	assert.Equal(t, roleList.Result[0].Privilege.Users, 2)
	assert.Equal(t, roleList.Result[0].Privilege.Roles, 2)
	assert.Equal(t, roleList.Result[0].Privilege.SamlUsers, 2)
	assert.Equal(t, roleList.Result[0].Privilege.SamlRoles, 2)
	assert.Equal(t, roleList.Result[0].Privilege.SamlSSOs, 2)
	assert.Equal(t, roleList.Result[0].Privilege.GlobalSetting, 2)
	assert.Equal(t, roleList.Result[0].Privilege.ExportData, 2)
	assert.Equal(t, roleList.Result[0].Privilege.ExportGlobalLog, 2)
	assert.Equal(t, roleList.Result[0].Privilege.Hotspot, 2)
	assert.Equal(t, roleList.Result[0].Privilege.Statics, 2)
	assert.Equal(t, roleList.Result[0].Privilege.Map, 2)
	assert.Equal(t, roleList.Result[0].Privilege.Clients, 2)
	assert.Equal(t, roleList.Result[0].Privilege.Insight, 2)
	assert.Equal(t, roleList.Result[0].Privilege.Report, 2)
	assert.Equal(t, roleList.Result[0].Privilege.Network, 2)
	assert.Equal(t, roleList.Result[0].Privilege.DeviceAccount, 2)
	assert.Equal(t, roleList.Result[0].Privilege.Anomaly, 2)
	assert.Equal(t, roleList.Result[0].Privilege.Analyze, 2)
	assert.Equal(t, roleList.Result[0].Privilege.SiteAnalyze, 2)

}
func TestOmadaClient_GetRoleInfo_ReturnsRoleInfo(t *testing.T) {
	mockMux := http.NewServeMux()
	mockMux.HandleFunc("/openapi/authorize/token", func(w http.ResponseWriter, r *http.Request) { mockValidTokenResponse(t, w, r) })
	mockMux.HandleFunc("/openapi/v1/my-cid/roles/master_admin_id", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "AccessToken=my-token", r.Header.Get("Authorization"))
		_, err := w.Write([]byte(`{
		  "errorCode": 0,
		  "msg": "Success.",
		  "result":
			{
			  "id": "master_admin_id",
			  "name": "Main Administrator",
			  "type": 0,
			  "defaultRole": true,
			  "source": 0,
			  "privilege": {
				"license": 2,
				"globalDashboard": 2,
				"dashboard": 2,
				"devices": 2,
				"adopt": 2,
				"globalLog": 2,
				"log": 2,
				"licenseBind": 2,
				"users": 2,
				"roles": 2,
				"samlUsers": 2,
				"samlRoles": 2,
				"samlSsos": 2,
				"globalSetting": 2,
				"exportData": 2,
				"globalExportData": 2,
				"exportGlobalLog": 2,
				"hotspot": 2,
				"statics": 2,
				"map": 2,
				"clients": 2,
				"insight": 2,
				"report": 2,
				"network": 2,
				"deviceAccount": 2,
				"anomaly": 2,
				"analyze": 2,
				"siteAnalyze": 2
			  }
			}}`))
		assert.NoError(t, err)
	})
	server := httptest.NewServer(mockMux)
	defer server.Close()

	c := NewClient(server.URL, "my-cid", "my-client-id", "my-client-secret", true)
	roleList, err := c.GetRoleInfo("master_admin_id")

	assert.NoError(t, err)
	assert.Equal(t, roleList.ErrorCode, 0)
	assert.Equal(t, roleList.Message, "Success.")
	assert.Equal(t, roleList.Result.Id, "master_admin_id")
	assert.Equal(t, roleList.Result.Name, "Main Administrator")
	assert.Equal(t, roleList.Result.Type, 0)
	assert.Equal(t, roleList.Result.DefaultRole, true)
	assert.Equal(t, roleList.Result.Source, 0)
	assert.Equal(t, roleList.Result.Privilege.License, 2)
	assert.Equal(t, roleList.Result.Privilege.GlobalDashboard, 2)
	assert.Equal(t, roleList.Result.Privilege.Dashboard, 2)
	assert.Equal(t, roleList.Result.Privilege.Devices, 2)
	assert.Equal(t, roleList.Result.Privilege.Adopt, 2)
	assert.Equal(t, roleList.Result.Privilege.GlobalLog, 2)
	assert.Equal(t, roleList.Result.Privilege.Log, 2)
	assert.Equal(t, roleList.Result.Privilege.LicenseBind, 2)
	assert.Equal(t, roleList.Result.Privilege.Users, 2)
	assert.Equal(t, roleList.Result.Privilege.Roles, 2)
	assert.Equal(t, roleList.Result.Privilege.SamlUsers, 2)
	assert.Equal(t, roleList.Result.Privilege.SamlRoles, 2)
	assert.Equal(t, roleList.Result.Privilege.SamlSSOs, 2)
	assert.Equal(t, roleList.Result.Privilege.GlobalSetting, 2)
	assert.Equal(t, roleList.Result.Privilege.ExportData, 2)
	assert.Equal(t, roleList.Result.Privilege.ExportGlobalLog, 2)
	assert.Equal(t, roleList.Result.Privilege.Hotspot, 2)
	assert.Equal(t, roleList.Result.Privilege.Statics, 2)
	assert.Equal(t, roleList.Result.Privilege.Map, 2)
	assert.Equal(t, roleList.Result.Privilege.Clients, 2)
	assert.Equal(t, roleList.Result.Privilege.Insight, 2)
	assert.Equal(t, roleList.Result.Privilege.Report, 2)
	assert.Equal(t, roleList.Result.Privilege.Network, 2)
	assert.Equal(t, roleList.Result.Privilege.DeviceAccount, 2)
	assert.Equal(t, roleList.Result.Privilege.Anomaly, 2)
	assert.Equal(t, roleList.Result.Privilege.Analyze, 2)
	assert.Equal(t, roleList.Result.Privilege.SiteAnalyze, 2)
}

func TestOmadaClient_GetSiteList_ReturnsAValidSiteList(t *testing.T) {
	mockMux := http.NewServeMux()
	mockMux.HandleFunc("/openapi/authorize/token", func(w http.ResponseWriter, r *http.Request) { mockValidTokenResponse(t, w, r) })
	mockMux.HandleFunc("/openapi/v1/my-cid/sites", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "AccessToken=my-token", r.Header.Get("Authorization"))
		assert.Equal(t, "100", r.URL.Query().Get("pageSize"))
		assert.Equal(t, "1", r.URL.Query().Get("page"))
		_, err := w.Write([]byte(`{
		  "errorCode": 0,
		  "msg": "Success.",
		  "result": {
			"totalRows": 2,
			"currentPage": 1,
			"currentSize": 100,
			"data": [
			  {
				"siteId": "65335c80b0dd10259f9ec5b",
				"name": "me-site",
				"region": "United States",
				"timeZone": "UTC",
				"scenario": "Hotel",
				"type": 0
			  },
			  {
				"siteId": "6534e46e0b0dd10259f9ec84",
				"name": "site2",
				"region": "United States",
				"timeZone": "UTC",
				"scenario": "Shopping Mall",
				"type": 0
			  }
			]
		  }
		}`))
		assert.NoError(t, err)
	})
	server := httptest.NewServer(mockMux)
	defer server.Close()

	c := NewClient(server.URL, "my-cid", "my-client-id", "my-client-secret", true)
	siteList, err := c.GetSiteList(1)

	assert.NoError(t, err)
	assert.Equal(t, siteList.ErrorCode, 0)
	assert.Equal(t, siteList.Message, "Success.")
	assert.Equal(t, siteList.Result.TotalRows, 2)
	assert.Equal(t, siteList.Result.CurrentPage, 1)
	assert.Equal(t, siteList.Result.CurrentSize, 100)
	assert.Equal(t, siteList.Result.Data[0].SiteId, "65335c80b0dd10259f9ec5b")
	assert.Equal(t, siteList.Result.Data[0].Name, "me-site")
	assert.Equal(t, siteList.Result.Data[0].Region, "United States")
	assert.Equal(t, siteList.Result.Data[0].TimeZone, "UTC")
	assert.Equal(t, siteList.Result.Data[0].Scenario, "Hotel")
	assert.Equal(t, siteList.Result.Data[0].Type, 0)
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
