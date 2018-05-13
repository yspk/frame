package util

import (
	"coding.net/baoquan2017/candy-backend/src/common/logger"
	"net/url"
)

/*
	This method allows to call cmn api to get area_code.
*/
func GetAreaByIp(ip, cmnBindAddr string) (string, string) {
	type listResp struct {
		ErrCode int    `json:"error_code"`
		Country string `json:"country"`
		Data    string `json:"data"`
	}

	values := url.Values{}
	values.Add("ip", ip)
	apiUrl := cmnBindAddr + "/cmn/v1.0/ip/area?" + values.Encode()

	var response listResp

	if err := Get(apiUrl, &response, nil); err != nil {
		logger.Error(err)
		return "", ""
	}

	return response.Data, response.Country
}
