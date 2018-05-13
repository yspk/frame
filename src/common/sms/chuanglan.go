package sms

import (
	"bytes"
	"coding.net/baoquan2017/candy-backend/src/common/email/validate"
	"coding.net/baoquan2017/candy-backend/src/common/logger"
	"coding.net/baoquan2017/candy-backend/src/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	codeV     = validate.NewCodeValidate(validate.NewMemoryStore(60))
	chuangLan config.ChuangLan
)

func init() {
	chuangLan = config.GetChuangLanConfig()
}

type smsClRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Msg      string `json:"msg"`
	Phone    string `json:"phone"`
	Sendtime string `json:"sendtime"`
	Report   bool   `json:"report"`
	Extend   string `json:"extend"`
	Uid      string `json:"uid"`
}

type SmsClResponse struct {
	Code     int    `json:"code"`
	MsgId    string `json:"msgId"`
	ErrorMsg string `json:"errorMsg"`
	Time     string `json:"time"`
}

/*
	This method allows to call Chuanglan to send SMS
*/
func SendSMSCl(mobile string) *SmsClResponse {
	code, err := codeV.Generate(mobile)
	if err != nil {
		logger.Error(err)
		return nil
	}

	// generate api url
	apiUrl := "http://smsssh1.253.com/msg/send/json"

	content := fmt.Sprintf(chuangLan.Template, code)

	request := smsClRequest{
		Account:  chuangLan.Account,
		Password: chuangLan.Password,
		Msg:      content,
		Phone:    mobile,
		Sendtime: time.Now().Format("201704101400"),
		Report:   false,
		Extend:   "",
		Uid:      chuangLan.Uid,
	}

	b, _ := json.Marshal(request)
	body := bytes.NewReader(b)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", apiUrl, body)
	req.Header.Set("Content-Type", "application/json")

	var response SmsClResponse
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return nil
	}

	defer resp.Body.Close()
	rb, err := ioutil.ReadAll(resp.Body)
	logger.Debug("sms response :" + string(rb))
	if err != nil {
		logger.Error(err)
		return nil
	}

	if err = json.Unmarshal(rb, &response); err != nil {
		logger.Error(err)
		return nil
	}
	return &response
}

func MobileCodeValidate(code, mobile string) bool {
	isValid, err := codeV.Validate(mobile, code)
	if err != nil {
		logger.Error(err)
	}
	return isValid
}
