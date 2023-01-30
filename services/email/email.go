package email

import (
	"fmt"

	"github.com/pnnh/multiverse-server/config"

	"golang.org/x/time/rate"

	gomail "gopkg.in/gomail.v2"
)

type Service struct {
	dialer  *gomail.Dialer
	limiter *rate.Limiter
}

func NewService() *Service {
	serv := &Service{}
	d := gomail.NewDialer(config.MailHost, config.MailPort, config.MailUser, config.MailPassword)
	serv.dialer = d

	serv.limiter = rate.NewLimiter(1, 1)
	return serv
}

func (s *Service) Init() error {
	return nil
}

func (s *Service) SendMessage(from, subject, mailBody string, to ...string) error {
	if !s.limiter.Allow() {
		return fmt.Errorf("邮件发送频率过高")
	}
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", mailBody)

	return s.dialer.DialAndSend(m)
}

func (s *Service) Start() error {
	return nil
}
