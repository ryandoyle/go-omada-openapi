package omada

import (
	"fmt"
	"net/http"
)

func (c *OmadaClient) GetRoleList() (*GetRoleListResponse, error) {
	path := fmt.Sprintf("%s/openapi/v1/%s/roles", c.baseUrl, c.omadaCId)
	request, err := http.NewRequest("GET", path, nil)

	roleListResponse := &GetRoleListResponse{}
	err = c.httpDoWrapped(request, roleListResponse)
	if err != nil {
		return nil, err
	}
	return roleListResponse, nil
}

func (c *OmadaClient) GetRoleInfo(roleId string) (*GetRoleInfoResponse, error) {
	path := fmt.Sprintf("%s/openapi/v1/%s/roles/%s", c.baseUrl, c.omadaCId, roleId)
	request, err := http.NewRequest("GET", path, nil)

	roleInfoResponse := &GetRoleInfoResponse{}
	err = c.httpDoWrapped(request, roleInfoResponse)
	if err != nil {
		return nil, err
	}
	return roleInfoResponse, nil
}

type GetRoleListResponse struct {
	ErrorCode int                      `json:"errorCode"`
	Message   string                   `json:"msg"`
	Result    []ControllerRoleDetailVO `json:"result"`
}

type GetRoleInfoResponse struct {
	ErrorCode int                    `json:"errorCode"`
	Message   string                 `json:"msg"`
	Result    ControllerRoleDetailVO `json:"result"`
}

type ControllerRoleDetailVO struct {
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
}
