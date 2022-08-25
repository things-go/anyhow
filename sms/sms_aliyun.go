package sms

import (
	"errors"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

var (
	ErrSendTooFrequent = errors.New("发送短信过于频繁")
	ErrSendFailed      = errors.New("发送短信失败")
)

type Client struct {
	*dysmsapi.Client
}

func NewAliyun(c *dysmsapi.Client) *Client {
	return &Client{
		c,
	}
}
func (*Client) Name() string { return "AliyunSms " }

func (sf *Client) SendMessage(mobile string, req SendMessageRequest) error {
	request := dysmsapi.CreateSendSmsRequest()
	request.PhoneNumbers = mobile             // 目标手机号
	request.SignName = req.SignName           // 短信签名名称
	request.TemplateCode = req.TemplateCode   // 短信模板id
	request.TemplateParam = req.TemplateParam // 短信模板变量对应的实际值，JSON格式

	response, err := sf.Client.SendSms(request)
	if err != nil {
		return err
	}
	// 默认流控：使用同一个签名，对同一个手机号码发送短信验证码，
	// 支持1条/分钟，5条/小时 ，累计10条/天
	if response.Code == "isv.BUSINESS_LIMIT_CONTROL" {
		return ErrSendTooFrequent
	}
	if response.Code != "OK" {
		return ErrSendFailed
	}
	return nil
}
