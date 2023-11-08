package omada

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOmadaClient_GetClientList_ReturnsAValidClientList(t *testing.T) {
	mockMux := http.NewServeMux()
	mockMux.HandleFunc("/openapi/authorize/token", func(w http.ResponseWriter, r *http.Request) { mockValidTokenResponse(t, w, r) })
	mockMux.HandleFunc("/openapi/v1/my-cid/sites/me-site/clients", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "AccessToken=my-token", r.Header.Get("Authorization"))
		assert.Equal(t, "1", r.URL.Query().Get("page"))
		assert.Equal(t, "100", r.URL.Query().Get("pageSize"))
		_, err := w.Write([]byte(`{
		  "errorCode": 0,
		  "msg": "Success.",
		  "result": {
			"totalRows": 1,
			"currentPage": 1,
			"currentSize": 100,
			"data": [
			  {
				"id": "Some ID",
				"mac": "AA-BB-CC-DD-EE-FF",
				"name": "My Device",
				"hostName": "My Device Hostname",
				"vendor": "My Vendor",
				"deviceType": "unknown",
				"deviceCategory": "computer",
				"osName": "Some OS",
				"ip": "192.168.1.100",
				"ipv6List": [ "2001:0db8:85a3:0000:0000:8a2e:0370:7334" ],
				"connectType": 1,
				"connectDevType": "ap",
				"connectedToWirelessRouter": false,
				"wireless": true,
				"ssid": "My SSID",
				"signalLevel": 85,
				"healthScore": -1,
				"signalRank": 4,
				"wifiMode": 5,
				"apName": "My Access Point",
				"apMac": "11-22-33-44-55-66",
				"radioId": 1,
				"channel": 36,
				"rxRate": 780000,
				"txRate": 560000,
				"powerSave": false,
				"rssi": -56,
				"snr": 39,
				"switchMac": "44-55-44-55-44-55",
				"switchName": "My Switch",
				"gatewayMac": "66-77-66-77-66-77",
				"gatewayName": "My Gateway",
				"vid": 0,
				"networkName": "My Network Name",
				"dot1xIdentity": "802.1x identity",
				"dot1xVlan": 4,
				"port": 16,
				"lagId": 2,
				"activity": 246,
				"trafficDown": 159561759,
				"trafficUp": 121189972,
				"uptime": 291632,
				"lastSeen": 1698744232047,
				"authStatus": 0,
				"guest": false,
				"active": true,
				"manager": false,
				"ipSetting": {
					"useFixedAddr": true,
					"netId": "My netid",
					"ip": "10.2.2.2"
				},
				"downPacket": 1598847,
				"upPacket": 1546297,
				"rateLimit": {
					"mode": 1,
					"rateLimitProfileId": "Rate limit profileId",
					"customRateLimit": {
						"downLimitEnable": true,
						"downLimit": 8000000,
						"upLimitEnable": true,
						"upLimit": 6000000
					}
				},
				"clientLockToApSetting": {
					"enable": true,
					"aps": [
						{ "name": "Locked AP Name", "mac": "22-22-22-22-22-22" }
					]
				},
				"support5g2": true,
				"multiLink": [
					{
						"radioId": 5,
						"wifiMode": 2,
						"channel": 44,
						"rxRate": 80000000,
						"txRate": 70000000,
						"powerSave": true,
						"rssi": -27,
						"snr": 80,
						"signalLevel": 90,
						"signalRank": 7,
						"upPacket": 121231,
						"downPacket": 456456,
						"trafficDown": 456456456,
						"trafficUp": 123123123,
						"activity": 552234,
						"signalLevelAndRank": 77
					}
				],
				"unit": 82,
				"standardPort": "Std port string"
			  }
			],
			"clientStat": {
			  "total": 11,
			  "wireless": 12,
			  "wired": 13,
			  "num2g": 14,
			  "num5g": 15,
			  "num6g": 16,
			  "numUser": 17,
			  "numGuest": 18,
			  "numWirelessUser": 19,
			  "numWirelessGuest": 20,
			  "num2gUser": 21,
			  "num5gUser": 22,
			  "num6gUser": 23,
			  "num2gGuest": 24,
			  "num5gGuest": 25,
			  "num6gGuest": 26,
			  "poor": 27,
			  "fair": 28,
			  "noData": 29,
			  "good": 30
			}
		  }
		}`))
		assert.NoError(t, err)
	})
	server := httptest.NewServer(mockMux)
	defer server.Close()

	c := NewClient(server.URL, "my-cid", "my-client-id", "my-client-secret", true)
	clientList, err := c.GetClientList("me-site", 1)

	assert.NoError(t, err)
	assert.Equal(t, clientList.ErrorCode, 0)
	assert.Equal(t, clientList.Message, "Success.")
	assert.Equal(t, clientList.Result.TotalRows, int64(1))
	assert.Equal(t, clientList.Result.CurrentSize, int32(100))
	assert.Equal(t, clientList.Result.CurrentPage, int32(1))
	assert.Equal(t, clientList.Result.Data[0].Id, "Some ID")
	assert.Equal(t, clientList.Result.Data[0].MAC, "AA-BB-CC-DD-EE-FF")
	assert.Equal(t, clientList.Result.Data[0].Name, "My Device")
	assert.Equal(t, clientList.Result.Data[0].HostName, "My Device Hostname")
	assert.Equal(t, clientList.Result.Data[0].Vendor, "My Vendor")
	assert.Equal(t, clientList.Result.Data[0].DeviceType, "unknown")
	assert.Equal(t, clientList.Result.Data[0].DeviceCategory, "computer")
	assert.Equal(t, clientList.Result.Data[0].OsName, "Some OS")
	assert.Equal(t, clientList.Result.Data[0].IP, "192.168.1.100")
	assert.Equal(t, clientList.Result.Data[0].IPv6List[0], "2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	assert.Equal(t, clientList.Result.Data[0].ConnectType, 1)
	assert.Equal(t, clientList.Result.Data[0].ConnectDevType, "ap")
	assert.Equal(t, clientList.Result.Data[0].ConnectedToWirelessRouter, false)
	assert.Equal(t, clientList.Result.Data[0].Wireless, true)
	assert.Equal(t, clientList.Result.Data[0].SSID, "My SSID")
	assert.Equal(t, clientList.Result.Data[0].SignalLevel, 85)
	assert.Equal(t, clientList.Result.Data[0].HealthScore, -1)
	assert.Equal(t, clientList.Result.Data[0].SignalRank, 4)
	assert.Equal(t, clientList.Result.Data[0].WifiMode, 5)
	assert.Equal(t, clientList.Result.Data[0].APName, "My Access Point")
	assert.Equal(t, clientList.Result.Data[0].APMac, "11-22-33-44-55-66")
	assert.Equal(t, clientList.Result.Data[0].RadioId, 1)
	assert.Equal(t, clientList.Result.Data[0].Channel, 36)
	assert.Equal(t, clientList.Result.Data[0].RxRate, 780000)
	assert.Equal(t, clientList.Result.Data[0].TxRate, 560000)
	assert.Equal(t, clientList.Result.Data[0].PowerSave, false)
	assert.Equal(t, clientList.Result.Data[0].RSSI, -56)
	assert.Equal(t, clientList.Result.Data[0].SNR, 39)
	assert.Equal(t, clientList.Result.Data[0].SwitchMac, "44-55-44-55-44-55")
	assert.Equal(t, clientList.Result.Data[0].SwitchName, "My Switch")
	assert.Equal(t, clientList.Result.Data[0].GatewayMac, "66-77-66-77-66-77")
	assert.Equal(t, clientList.Result.Data[0].GatewayName, "My Gateway")
	assert.Equal(t, clientList.Result.Data[0].VID, 0)
	assert.Equal(t, clientList.Result.Data[0].NetworkName, "My Network Name")
	assert.Equal(t, clientList.Result.Data[0].Dot1xIdentity, "802.1x identity")
	assert.Equal(t, clientList.Result.Data[0].Dot1xVlan, 4)
	assert.Equal(t, clientList.Result.Data[0].Port, 16)
	assert.Equal(t, clientList.Result.Data[0].LagID, 2)
	assert.Equal(t, clientList.Result.Data[0].Activity, 246)
	assert.Equal(t, clientList.Result.Data[0].TrafficDown, 159561759)
	assert.Equal(t, clientList.Result.Data[0].TrafficUp, 121189972)
	assert.Equal(t, clientList.Result.Data[0].Uptime, 291632)
	assert.Equal(t, clientList.Result.Data[0].LastSeen, 1698744232047)
	assert.Equal(t, clientList.Result.Data[0].AuthStatus, 0)
	assert.Equal(t, clientList.Result.Data[0].Blocked, false)
	assert.Equal(t, clientList.Result.Data[0].Guest, false)
	assert.Equal(t, clientList.Result.Data[0].Active, true)
	assert.Equal(t, clientList.Result.Data[0].Manager, false)
	assert.Equal(t, clientList.Result.Data[0].IpSetting.UseFixedAddr, true)
	assert.Equal(t, clientList.Result.Data[0].IpSetting.NetId, "My netid")
	assert.Equal(t, clientList.Result.Data[0].IpSetting.IP, "10.2.2.2")
	assert.Equal(t, clientList.Result.Data[0].DownPacket, 1598847)
	assert.Equal(t, clientList.Result.Data[0].UpPacket, 1546297)
	assert.Equal(t, clientList.Result.Data[0].RateLimit.Mode, 1)
	assert.Equal(t, clientList.Result.Data[0].RateLimit.RateLimitProfileId, "Rate limit profileId")
	assert.Equal(t, clientList.Result.Data[0].RateLimit.CustomRateLimit.DownLimitEnable, true)
	assert.Equal(t, clientList.Result.Data[0].RateLimit.CustomRateLimit.DownLimit, 8000000)
	assert.Equal(t, clientList.Result.Data[0].RateLimit.CustomRateLimit.UpLimitEnable, true)
	assert.Equal(t, clientList.Result.Data[0].RateLimit.CustomRateLimit.UpLimit, 6000000)
	assert.Equal(t, clientList.Result.Data[0].ClientLockToApSetting.Enable, true)
	assert.Equal(t, clientList.Result.Data[0].ClientLockToApSetting.APs[0].Name, "Locked AP Name")
	assert.Equal(t, clientList.Result.Data[0].ClientLockToApSetting.APs[0].MAC, "22-22-22-22-22-22")
	assert.Equal(t, clientList.Result.Data[0].Support5g2, true)
	assert.Equal(t, clientList.Result.Data[0].MultiLink[0].RadioId, 5)
	assert.Equal(t, clientList.Result.Data[0].MultiLink[0].WifiMode, 2)
	assert.Equal(t, clientList.Result.Data[0].MultiLink[0].Channel, 44)
	assert.Equal(t, clientList.Result.Data[0].MultiLink[0].RxRate, 80000000)
	assert.Equal(t, clientList.Result.Data[0].MultiLink[0].TxRate, 70000000)
	assert.Equal(t, clientList.Result.Data[0].MultiLink[0].PowerSave, true)
	assert.Equal(t, clientList.Result.Data[0].MultiLink[0].RSSI, -27)
	assert.Equal(t, clientList.Result.Data[0].MultiLink[0].SNR, 80)
	assert.Equal(t, clientList.Result.Data[0].MultiLink[0].SignalLevel, 90)
	assert.Equal(t, clientList.Result.Data[0].MultiLink[0].SignalRank, 7)
	assert.Equal(t, clientList.Result.Data[0].MultiLink[0].UpPacket, 121231)
	assert.Equal(t, clientList.Result.Data[0].MultiLink[0].DownPacket, 456456)
	assert.Equal(t, clientList.Result.Data[0].MultiLink[0].TrafficDown, 456456456)
	assert.Equal(t, clientList.Result.Data[0].MultiLink[0].TrafficUp, 123123123)
	assert.Equal(t, clientList.Result.Data[0].MultiLink[0].Activity, 552234)
	assert.Equal(t, clientList.Result.Data[0].MultiLink[0].SignalLevelAndRank, 77)
	assert.Equal(t, clientList.Result.Data[0].Unit, 82)
	assert.Equal(t, clientList.Result.Data[0].StandardPort, "Std port string")
	assert.Equal(t, clientList.Result.ClientStat.Total, 11)
	assert.Equal(t, clientList.Result.ClientStat.Wireless, 12)
	assert.Equal(t, clientList.Result.ClientStat.Wired, 13)
	assert.Equal(t, clientList.Result.ClientStat.Num2g, 14)
	assert.Equal(t, clientList.Result.ClientStat.Num5g, 15)
	assert.Equal(t, clientList.Result.ClientStat.Num6g, 16)
	assert.Equal(t, clientList.Result.ClientStat.NumUser, 17)
	assert.Equal(t, clientList.Result.ClientStat.NumGuest, 18)
	assert.Equal(t, clientList.Result.ClientStat.NumWirelessUser, 19)
	assert.Equal(t, clientList.Result.ClientStat.NumWirelessGuest, 20)
	assert.Equal(t, clientList.Result.ClientStat.Num2gUser, 21)
	assert.Equal(t, clientList.Result.ClientStat.Num5gUser, 22)
	assert.Equal(t, clientList.Result.ClientStat.Num6gUser, 23)
	assert.Equal(t, clientList.Result.ClientStat.Num2gGuest, 24)
	assert.Equal(t, clientList.Result.ClientStat.Num5gGuest, 25)
	assert.Equal(t, clientList.Result.ClientStat.Num6gGuest, 26)
	assert.Equal(t, clientList.Result.ClientStat.Poor, 27)
	assert.Equal(t, clientList.Result.ClientStat.Fair, 28)
	assert.Equal(t, clientList.Result.ClientStat.NoData, 29)
	assert.Equal(t, clientList.Result.ClientStat.Good, 30)

}

func TestOmadaClient_GetClientInfo_ReturnsAValidClientInfo(t *testing.T) {
	mockMux := http.NewServeMux()
	mockMux.HandleFunc("/openapi/authorize/token", func(w http.ResponseWriter, r *http.Request) { mockValidTokenResponse(t, w, r) })
	mockMux.HandleFunc("/openapi/v1/my-cid/sites/me-site/clients/11-22-33-44-55-66", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "AccessToken=my-token", r.Header.Get("Authorization"))
		_, err := w.Write([]byte(`{
		  "errorCode": 0,
		  "msg": "Success.",
		  "result": {
				"id": "Some ID",
				"mac": "AA-BB-CC-DD-EE-FF",
				"name": "My Device",
				"hostName": "My Device Hostname",
				"vendor": "My Vendor",
				"deviceType": "unknown",
				"deviceCategory": "computer",
				"osName": "Some OS",
				"ip": "192.168.1.100",
				"ipv6List": [ "2001:0db8:85a3:0000:0000:8a2e:0370:7334" ],
				"connectType": 1,
				"connectDevType": "ap",
				"connectedToWirelessRouter": false,
				"wireless": true,
				"ssid": "My SSID",
				"signalLevel": 85,
				"healthScore": -1,
				"signalRank": 4,
				"wifiMode": 5,
				"apName": "My Access Point",
				"apMac": "11-22-33-44-55-66",
				"radioId": 1,
				"channel": 36,
				"rxRate": 780000,
				"txRate": 560000,
				"powerSave": false,
				"rssi": -56,
				"snr": 39,
				"switchMac": "44-55-44-55-44-55",
				"switchName": "My Switch",
				"gatewayMac": "66-77-66-77-66-77",
				"gatewayName": "My Gateway",
				"vid": 0,
				"networkName": "My Network Name",
				"dot1xIdentity": "802.1x identity",
				"dot1xVlan": 4,
				"port": 16,
				"lagId": 2,
				"activity": 246,
				"trafficDown": 159561759,
				"trafficUp": 121189972,
				"uptime": 291632,
				"lastSeen": 1698744232047,
				"authStatus": 0,
				"guest": false,
				"active": true,
				"manager": false,
				"ipSetting": {
					"useFixedAddr": true,
					"netId": "My netid",
					"ip": "10.2.2.2"
				},
				"downPacket": 1598847,
				"upPacket": 1546297,
				"rateLimit": {
					"mode": 1,
					"rateLimitProfileId": "Rate limit profileId",
					"customRateLimit": {
						"downLimitEnable": true,
						"downLimit": 8000000,
						"upLimitEnable": true,
						"upLimit": 6000000
					}
				},
				"clientLockToApSetting": {
					"enable": true,
					"aps": [
						{ "name": "Locked AP Name", "mac": "22-22-22-22-22-22" }
					]
				},
				"support5g2": true,
				"multiLink": [
					{
						"radioId": 5,
						"wifiMode": 2,
						"channel": 44,
						"rxRate": 80000000,
						"txRate": 70000000,
						"powerSave": true,
						"rssi": -27,
						"snr": 80,
						"signalLevel": 90,
						"signalRank": 7,
						"upPacket": 121231,
						"downPacket": 456456,
						"trafficDown": 456456456,
						"trafficUp": 123123123,
						"activity": 552234,
						"signalLevelAndRank": 77
					}
				],
				"unit": 82,
				"standardPort": "Std port string"
		  }
		}`))
		assert.NoError(t, err)
	})
	server := httptest.NewServer(mockMux)
	defer server.Close()

	c := NewClient(server.URL, "my-cid", "my-client-id", "my-client-secret", true)
	clientInfo, err := c.GetClientInfo("me-site", "11-22-33-44-55-66")

	assert.NoError(t, err)
	assert.Equal(t, clientInfo.ErrorCode, 0)
	assert.Equal(t, clientInfo.Message, "Success.")
	assert.Equal(t, clientInfo.Result.Id, "Some ID")
	assert.Equal(t, clientInfo.Result.MAC, "AA-BB-CC-DD-EE-FF")
	assert.Equal(t, clientInfo.Result.Name, "My Device")
	assert.Equal(t, clientInfo.Result.HostName, "My Device Hostname")
	assert.Equal(t, clientInfo.Result.Vendor, "My Vendor")
	assert.Equal(t, clientInfo.Result.DeviceType, "unknown")
	assert.Equal(t, clientInfo.Result.DeviceCategory, "computer")
	assert.Equal(t, clientInfo.Result.OsName, "Some OS")
	assert.Equal(t, clientInfo.Result.IP, "192.168.1.100")
	assert.Equal(t, clientInfo.Result.IPv6List[0], "2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	assert.Equal(t, clientInfo.Result.ConnectType, 1)
	assert.Equal(t, clientInfo.Result.ConnectDevType, "ap")
	assert.Equal(t, clientInfo.Result.ConnectedToWirelessRouter, false)
	assert.Equal(t, clientInfo.Result.Wireless, true)
	assert.Equal(t, clientInfo.Result.SSID, "My SSID")
	assert.Equal(t, clientInfo.Result.SignalLevel, 85)
	assert.Equal(t, clientInfo.Result.HealthScore, -1)
	assert.Equal(t, clientInfo.Result.SignalRank, 4)
	assert.Equal(t, clientInfo.Result.WifiMode, 5)
	assert.Equal(t, clientInfo.Result.APName, "My Access Point")
	assert.Equal(t, clientInfo.Result.APMac, "11-22-33-44-55-66")
	assert.Equal(t, clientInfo.Result.RadioId, 1)
	assert.Equal(t, clientInfo.Result.Channel, 36)
	assert.Equal(t, clientInfo.Result.RxRate, 780000)
	assert.Equal(t, clientInfo.Result.TxRate, 560000)
	assert.Equal(t, clientInfo.Result.PowerSave, false)
	assert.Equal(t, clientInfo.Result.RSSI, -56)
	assert.Equal(t, clientInfo.Result.SNR, 39)
	assert.Equal(t, clientInfo.Result.SwitchMac, "44-55-44-55-44-55")
	assert.Equal(t, clientInfo.Result.SwitchName, "My Switch")
	assert.Equal(t, clientInfo.Result.GatewayMac, "66-77-66-77-66-77")
	assert.Equal(t, clientInfo.Result.GatewayName, "My Gateway")
	assert.Equal(t, clientInfo.Result.VID, 0)
	assert.Equal(t, clientInfo.Result.NetworkName, "My Network Name")
	assert.Equal(t, clientInfo.Result.Dot1xIdentity, "802.1x identity")
	assert.Equal(t, clientInfo.Result.Dot1xVlan, 4)
	assert.Equal(t, clientInfo.Result.Port, 16)
	assert.Equal(t, clientInfo.Result.LagID, 2)
	assert.Equal(t, clientInfo.Result.Activity, 246)
	assert.Equal(t, clientInfo.Result.TrafficDown, 159561759)
	assert.Equal(t, clientInfo.Result.TrafficUp, 121189972)
	assert.Equal(t, clientInfo.Result.Uptime, 291632)
	assert.Equal(t, clientInfo.Result.LastSeen, 1698744232047)
	assert.Equal(t, clientInfo.Result.AuthStatus, 0)
	assert.Equal(t, clientInfo.Result.Blocked, false)
	assert.Equal(t, clientInfo.Result.Guest, false)
	assert.Equal(t, clientInfo.Result.Active, true)
	assert.Equal(t, clientInfo.Result.Manager, false)
	assert.Equal(t, clientInfo.Result.IpSetting.UseFixedAddr, true)
	assert.Equal(t, clientInfo.Result.IpSetting.NetId, "My netid")
	assert.Equal(t, clientInfo.Result.IpSetting.IP, "10.2.2.2")
	assert.Equal(t, clientInfo.Result.DownPacket, 1598847)
	assert.Equal(t, clientInfo.Result.UpPacket, 1546297)
	assert.Equal(t, clientInfo.Result.RateLimit.Mode, 1)
	assert.Equal(t, clientInfo.Result.RateLimit.RateLimitProfileId, "Rate limit profileId")
	assert.Equal(t, clientInfo.Result.RateLimit.CustomRateLimit.DownLimitEnable, true)
	assert.Equal(t, clientInfo.Result.RateLimit.CustomRateLimit.DownLimit, 8000000)
	assert.Equal(t, clientInfo.Result.RateLimit.CustomRateLimit.UpLimitEnable, true)
	assert.Equal(t, clientInfo.Result.RateLimit.CustomRateLimit.UpLimit, 6000000)
	assert.Equal(t, clientInfo.Result.ClientLockToApSetting.Enable, true)
	assert.Equal(t, clientInfo.Result.ClientLockToApSetting.APs[0].Name, "Locked AP Name")
	assert.Equal(t, clientInfo.Result.ClientLockToApSetting.APs[0].MAC, "22-22-22-22-22-22")
	assert.Equal(t, clientInfo.Result.Support5g2, true)
	assert.Equal(t, clientInfo.Result.MultiLink[0].RadioId, 5)
	assert.Equal(t, clientInfo.Result.MultiLink[0].WifiMode, 2)
	assert.Equal(t, clientInfo.Result.MultiLink[0].Channel, 44)
	assert.Equal(t, clientInfo.Result.MultiLink[0].RxRate, 80000000)
	assert.Equal(t, clientInfo.Result.MultiLink[0].TxRate, 70000000)
	assert.Equal(t, clientInfo.Result.MultiLink[0].PowerSave, true)
	assert.Equal(t, clientInfo.Result.MultiLink[0].RSSI, -27)
	assert.Equal(t, clientInfo.Result.MultiLink[0].SNR, 80)
	assert.Equal(t, clientInfo.Result.MultiLink[0].SignalLevel, 90)
	assert.Equal(t, clientInfo.Result.MultiLink[0].SignalRank, 7)
	assert.Equal(t, clientInfo.Result.MultiLink[0].UpPacket, 121231)
	assert.Equal(t, clientInfo.Result.MultiLink[0].DownPacket, 456456)
	assert.Equal(t, clientInfo.Result.MultiLink[0].TrafficDown, 456456456)
	assert.Equal(t, clientInfo.Result.MultiLink[0].TrafficUp, 123123123)
	assert.Equal(t, clientInfo.Result.MultiLink[0].Activity, 552234)
	assert.Equal(t, clientInfo.Result.MultiLink[0].SignalLevelAndRank, 77)
	assert.Equal(t, clientInfo.Result.Unit, 82)
}
