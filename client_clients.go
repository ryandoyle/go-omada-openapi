package omada

import (
	"fmt"
	"net/http"
)

func (c *OmadaClient) GetClientList(siteId string, page int) (*GetClientListResponse, error) {
	path := fmt.Sprintf("%s/openapi/v1/%s/sites/%s/clients?page=%d&pageSize=%d", c.baseUrl, c.omadaCId, siteId, page, c.PageSize)
	request, err := http.NewRequest("GET", path, nil)

	clientList := &GetClientListResponse{}
	err = c.httpDoWrapped(request, clientList)
	if err != nil {
		return nil, err
	}

	return clientList, nil
}

type GetClientListResponse struct {
	EnvelopeResponse
	Result struct {
		TotalRows   int64 `json:"totalRows"`
		CurrentPage int32 `json:"currentPage"`
		CurrentSize int32 `json:"currentSize"`
		Data        []struct {
			Id                        string   `json:"id"`
			MAC                       string   `json:"mac"`
			Name                      string   `json:"name"`
			HostName                  string   `json:"hostName"`
			Vendor                    string   `json:"vendor"`
			DeviceType                string   `json:"deviceType"`
			DeviceCategory            string   `json:"deviceCategory"`
			OsName                    string   `json:"osName"`
			IP                        string   `json:"ip"`
			IPv6List                  []string `json:"ipv6List"`
			ConnectType               int      `json:"connectType"`
			ConnectDevType            string   `json:"connectDevType"`
			ConnectedToWirelessRouter bool     `json:"connectedToWirelessRouter"`
			Wireless                  bool     `json:"wireless"`
			SSID                      string   `json:"ssid"`
			SignalLevel               int      `json:"signalLevel"`
			HealthScore               int      `json:"healthScore"`
			SignalRank                int      `json:"signalRank"`
			WifiMode                  int      `json:"wifiMode"`
			APName                    string   `json:"apName"`
			APMac                     string   `json:"apMac"`
			RadioId                   int      `json:"radioId"`
			Channel                   int      `json:"channel"`
			RxRate                    int      `json:"rxRate"`
			TxRate                    int      `json:"txRate"`
			PowerSave                 bool     `json:"powerSave"`
			RSSI                      int      `json:"rssi"`
			SNR                       int      `json:"snr"`
			SwitchMac                 string   `json:"switchMac"`
			SwitchName                string   `json:"switchName"`
			GatewayMac                string   `json:"gatewayMac"`
			GatewayName               string   `json:"gatewayName"`
			VID                       int      `json:"vid"`
			NetworkName               string   `json:"networkName"`
			Dot1xIdentity             string   `json:"dot1xIdentity"`
			Dot1xVlan                 int      `json:"dot1xVlan"`
			Port                      int      `json:"port"`
			LagID                     int      `json:"lagId"`
			Activity                  int      `json:"activity"`
			TrafficDown               int      `json:"trafficDown"`
			TrafficUp                 int      `json:"trafficUp"`
			Uptime                    int      `json:"uptime"`
			LastSeen                  int      `json:"lastSeen"`
			AuthStatus                int      `json:"authStatus"`
			Blocked                   bool     `json:"blocked"`
			Guest                     bool     `json:"guest"`
			Active                    bool     `json:"active"`
			Manager                   bool     `json:"manager"`
			IpSetting                 struct {
				UseFixedAddr bool   `json:"useFixedAddr"`
				NetId        string `json:"netId"`
				IP           string `json:"ip"`
			} `json:"ipSetting"`
			DownPacket int `json:"downPacket"`
			UpPacket   int `json:"upPacket"`
			RateLimit  struct {
				Mode               int    `json:"mode"`
				RateLimitProfileId string `json:"rateLimitProfileId"`
				CustomRateLimit    struct {
					DownLimit       int  `json:"downLimit"`
					DownLimitEnable bool `json:"downLimitEnable"`
					UpLimit         int  `json:"upLimit"`
					UpLimitEnable   bool `json:"upLimitEnable"`
				} `json:"customRateLimit"`
			} `json:"rateLimit"`
			ClientLockToApSetting struct {
				Enable bool `json:"enable"`
				APs    []struct {
					Name string `json:"name"`
					MAC  string `json:"mac"`
				} `json:"aps"`
			} `json:"clientLockToApSetting"`
			Support5g2 bool `json:"support5g2"`
			MultiLink  []struct {
				RadioId            int  `json:"radioId"`
				WifiMode           int  `json:"wifiMode"`
				Channel            int  `json:"channel"`
				RxRate             int  `json:"rxRate"`
				TxRate             int  `json:"txRate"`
				PowerSave          bool `json:"powerSave"`
				RSSI               int  `json:"rssi"`
				SNR                int  `json:"snr"`
				SignalLevel        int  `json:"signalLevel"`
				SignalRank         int  `json:"signalRank"`
				UpPacket           int  `json:"upPacket"`
				DownPacket         int  `json:"downPacket"`
				TrafficDown        int  `json:"trafficDown"`
				TrafficUp          int  `json:"trafficUp"`
				Activity           int  `json:"activity"`
				SignalLevelAndRank int  `json:"signalLevelAndRank"`
			} `json:"multiLink"`
			Unit         int    `json:"unit"`
			StandardPort string `json:"standardPort"`
		} `json:"data"`
		ClientStat struct {
			Total            int `json:"total"`
			Wireless         int `json:"wireless"`
			Wired            int `json:"wired"`
			Num2g            int `json:"num2g"`
			Num5g            int `json:"num5g"`
			Num6g            int `json:"num6g"`
			NumUser          int `json:"numUser"`
			NumGuest         int `json:"numGuest"`
			NumWirelessUser  int `json:"numWirelessUser"`
			NumWirelessGuest int `json:"numWirelessGuest"`
			Num2gUser        int `json:"num2gUser"`
			Num5gUser        int `json:"num5gUser"`
			Num6gUser        int `json:"num6gUser"`
			Num2gGuest       int `json:"num2gGuest"`
			Num5gGuest       int `json:"num5gGuest"`
			Num6gGuest       int `json:"num6gGuest"`
			Poor             int `json:"poor"`
			Fair             int `json:"fair"`
			NoData           int `json:"noData"`
			Good             int `json:"good"`
		} `json:"clientStat"`
	} `json:"result"`
}
