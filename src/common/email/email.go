package email

import (
	"bytes"
	"github.com/yspk/frame/src/common/constant"
	"github.com/yspk/frame/src/common/email/send"
	"github.com/yspk/frame/src/common/email/validate"
	"github.com/yspk/frame/src/common/logger"
	"github.com/yspk/frame/src/config"
	"net/mail"
	"sync"
)

var (
	codeV    = validate.NewCodeValidate(validate.NewMemoryStore(constant.DefaultGCInterval))
	valEmail config.ValEmail
)

func init() {
	valEmail = config.GetValidateEmail()
}

func SendEmail(email string) {
	code, err := codeV.Generate(email)
	if err != nil {
		logger.Error(err)
		return
	}

	sender, err := send.NewSmtpSender(valEmail.Host, mail.Address{valEmail.Name, valEmail.Account}, valEmail.Password)
	if err != nil {
		logger.Error(err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	msg := &send.Message{
		Subject: "糖果小镇注册验证",
		Content: bytes.NewBufferString(code),
		To:      []string{email},
	}
	err = sender.AsyncSend(msg, false, func(err error) {
		defer wg.Done()
		if err != nil {
			logger.Error(err)
		}
	})
	if err != nil {
		logger.Error(err)
	}
	wg.Wait()
}

func EmailCodeValidate(code, email string) bool {
	isValid, err := codeV.Validate(email, code)
	if err != nil {
		logger.Error(err)
	}
	return isValid
}
