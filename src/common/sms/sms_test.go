package sms

import (
	"coding.net/baoquan2017/candy-backend/src/common/logger"
	"testing"
)

func TestSendVerificationCode(t *testing.T) {
	response := SendSMSCl("15658836559")
	if response.ErrorMsg != "" {
		logger.Error(response)
	}
	logger.Info("sms sent successfully")
}
