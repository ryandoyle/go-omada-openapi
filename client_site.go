package omada

import (
	"fmt"
	"net/http"
)

func (c *OmadaClient) GetSiteList(page int) (*GetSiteListResponse, error) {
	path := fmt.Sprintf("%s/openapi/v1/%s/sites?pageSize=%d&page=%d", c.baseUrl, c.omadaCId, c.PageSize, page)
	request, err := http.NewRequest("GET", path, nil)

	siteList := &GetSiteListResponse{}
	err = c.httpDoWrapped(request, siteList)
	if err != nil {
		return nil, err
	}

	return siteList, nil
}

func (c *OmadaClient) GetSiteInfo(site string) (*GetSiteInfoResponse, error) {
	path := fmt.Sprintf("%s/openapi/v1/%s/sites/%s", c.baseUrl, c.omadaCId, site)
	request, err := http.NewRequest("GET", path, nil)

	siteInfo := &GetSiteInfoResponse{}
	err = c.httpDoWrapped(request, siteInfo)
	if err != nil {
		return nil, err
	}

	return siteInfo, nil
}

type GetSiteListResponse struct {
	EnvelopeResponse
	Result struct {
		TotalRows   int          `json:"totalRows"`
		CurrentPage int          `json:"currentPage"`
		CurrentSize int          `json:"currentSize"`
		Data        []SiteEntity `json:"data"`
	} `json:"result"`
}

type GetSiteInfoResponse struct {
	EnvelopeResponse
	Result SiteEntity `json:"result"`
}

type SiteEntity struct {
	SiteId    string  `json:"siteId"`
	Name      string  `json:"name"`
	Region    string  `json:"region"`
	TimeZone  string  `json:"timeZone"`
	Scenario  string  `json:"scenario"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Address   string  `json:"address"`
	Type      int     `json:"type"`
}

func (c *OmadaClient) GetScenarioList() (*GetScenarioListResponse, error) {
	path := fmt.Sprintf("%s/openapi/v1/%s/scenarios", c.baseUrl, c.omadaCId)
	request, err := http.NewRequest("GET", path, nil)

	scenario := &GetScenarioListResponse{}
	err = c.httpDoWrapped(request, scenario)
	if err != nil {
		return nil, err
	}

	return scenario, nil
}

type GetScenarioListResponse struct {
	EnvelopeResponse
	Result []string `json:"result"`
}

func (c *OmadaClient) GetSiteDeviceAccountSetting(siteId string) (*GetSiteDeviceAccountSettingResponse, error) {
	path := fmt.Sprintf("%s/openapi/v1/%s/sites/%s/device-account", c.baseUrl, c.omadaCId, siteId)
	request, err := http.NewRequest("GET", path, nil)

	scenario := &GetSiteDeviceAccountSettingResponse{}
	err = c.httpDoWrapped(request, scenario)
	if err != nil {
		return nil, err
	}

	return scenario, nil
}

type GetSiteDeviceAccountSettingResponse struct {
	EnvelopeResponse
	Result struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"result"`
}
