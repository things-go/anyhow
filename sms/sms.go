package sms

type SendMessageRequest struct {
	SignName      string // required, 签名
	TemplateCode  string // required, 模板代号, 如 SMS_xxxxxx
	TemplateParam string // optional, 模板的参数, 如 {"code": "123456"}
}
type Provider interface {
	Name() string
	SendMessage(mobile string, req SendMessageRequest) error
}

type Sms struct {
	p Provider
}

func New(p Provider) *Sms { return &Sms{p: p} }

func (sf *Sms) SendMessage(mobile string, req SendMessageRequest) error {
	return sf.p.SendMessage(mobile, req)
}
