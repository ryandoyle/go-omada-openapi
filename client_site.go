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

type GetSiteListResponse struct {
	Envelope
	Result struct {
		TotalRows   int `json:"totalRows"`
		CurrentPage int `json:"currentPage"`
		CurrentSize int `json:"currentSize"`
		Data        []struct {
			SiteId   string `json:"siteId"`
			Name     string `json:"name"`
			Region   string `json:"region"`
			TimeZone string `json:"timeZone"`
			Scenario string `json:"scenario"`
			Type     int    `json:"type"`
		} `json:"data"`
	} `json:"result"`
}
