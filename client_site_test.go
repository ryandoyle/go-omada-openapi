package omada

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
				"longitude": 123.345,
				"latitude": 55.55,
				"address": "123 fake street",
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
	assert.Equal(t, siteList.Result.Data[0].Longitude, 123.345)
	assert.Equal(t, siteList.Result.Data[0].Latitude, 55.55)
	assert.Equal(t, siteList.Result.Data[0].Address, "123 fake street")
	assert.Equal(t, siteList.Result.Data[0].Type, 0)
}

func TestOmadaClient_GetSiteInfo_ReturnsAValidSiteInfo(t *testing.T) {
	mockMux := http.NewServeMux()
	mockMux.HandleFunc("/openapi/authorize/token", func(w http.ResponseWriter, r *http.Request) { mockValidTokenResponse(t, w, r) })
	mockMux.HandleFunc("/openapi/v1/my-cid/sites/me-site", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "AccessToken=my-token", r.Header.Get("Authorization"))
		_, err := w.Write([]byte(`{
		  "errorCode": 0,
		  "msg": "Success.",
		  "result": 
		  	{
				"siteId": "65335c80b0dd10259f9ec5b",
				"name": "me-site",
				"region": "United States",
				"timeZone": "UTC",
				"scenario": "Hotel",
				"longitude": 123.345,
				"latitude": 55.55,
				"address": "123 fake street",
				"type": 0
			}
		}`))
		assert.NoError(t, err)
	})
	server := httptest.NewServer(mockMux)
	defer server.Close()

	c := NewClient(server.URL, "my-cid", "my-client-id", "my-client-secret", true)
	siteInfo, err := c.GetSiteInfo("me-site")

	assert.NoError(t, err)
	assert.Equal(t, siteInfo.ErrorCode, 0)
	assert.Equal(t, siteInfo.Message, "Success.")
	assert.Equal(t, siteInfo.Result.SiteId, "65335c80b0dd10259f9ec5b")
	assert.Equal(t, siteInfo.Result.Name, "me-site")
	assert.Equal(t, siteInfo.Result.Region, "United States")
	assert.Equal(t, siteInfo.Result.TimeZone, "UTC")
	assert.Equal(t, siteInfo.Result.Scenario, "Hotel")
	assert.Equal(t, siteInfo.Result.Longitude, 123.345)
	assert.Equal(t, siteInfo.Result.Latitude, 55.55)
	assert.Equal(t, siteInfo.Result.Address, "123 fake street")
	assert.Equal(t, siteInfo.Result.Type, 0)
}

func TestOmadaClient_GetScenarioList_ReturnsAValidScenarioList(t *testing.T) {
	mockMux := http.NewServeMux()
	mockMux.HandleFunc("/openapi/authorize/token", func(w http.ResponseWriter, r *http.Request) { mockValidTokenResponse(t, w, r) })
	mockMux.HandleFunc("/openapi/v1/my-cid/scenarios", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "AccessToken=my-token", r.Header.Get("Authorization"))
		_, err := w.Write([]byte(`{
		  "errorCode": 0,
		  "msg": "Success.",
		  "result": ["Hotel", "Shopping Mall"]
		}`))
		assert.NoError(t, err)
	})
	server := httptest.NewServer(mockMux)
	defer server.Close()

	c := NewClient(server.URL, "my-cid", "my-client-id", "my-client-secret", true)
	siteInfo, err := c.GetScenarioList()

	assert.NoError(t, err)
	assert.Equal(t, siteInfo.ErrorCode, 0)
	assert.Equal(t, siteInfo.Message, "Success.")
	assert.Contains(t, siteInfo.Result, "Hotel")
	assert.Contains(t, siteInfo.Result, "Shopping Mall")
}

func TestOmadaClient_GetSiteDeviceAccountSetting_ReturnsAValidSiteDeviceAccountSetting(t *testing.T) {
	mockMux := http.NewServeMux()
	mockMux.HandleFunc("/openapi/authorize/token", func(w http.ResponseWriter, r *http.Request) { mockValidTokenResponse(t, w, r) })
	mockMux.HandleFunc("/openapi/v1/my-cid/sites/my-site/device-account", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "AccessToken=my-token", r.Header.Get("Authorization"))
		_, err := w.Write([]byte(`{
		  "errorCode": 0,
		  "msg": "Success.",
		  "result": {"username": "me-user", "password": "me-password"}
		}`))
		assert.NoError(t, err)
	})
	server := httptest.NewServer(mockMux)
	defer server.Close()

	c := NewClient(server.URL, "my-cid", "my-client-id", "my-client-secret", true)
	deviceAccountSetting, err := c.GetSiteDeviceAccountSetting("my-site")

	assert.NoError(t, err)
	assert.Equal(t, deviceAccountSetting.ErrorCode, 0)
	assert.Equal(t, deviceAccountSetting.Message, "Success.")
	assert.Equal(t, deviceAccountSetting.Result.Username, "me-user")
	assert.Equal(t, deviceAccountSetting.Result.Password, "me-password")
}
