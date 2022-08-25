package captcha

import (
	"github.com/mojocn/base64Captcha"
	"github.com/things-go/clip/limit"
)

var _ limit.CaptchaProvider = (*Captcha)(nil)

type Captcha struct {
	d base64Captcha.Driver
}

func New(d base64Captcha.Driver) *Captcha {
	return &Captcha{
		d,
	}
}

func (c *Captcha) Name() string { return "base64Captcha" }

func (c *Captcha) GenerateQuestionAnswer() (*limit.QuestionAnswer, error) {
	id, q, a := c.d.GenerateIdQuestionAnswer()
	it, err := c.d.DrawCaptcha(q)
	if err != nil {
		return nil, err
	}
	return &limit.QuestionAnswer{
		Id:       id,
		Question: it.EncodeB64string(),
		Answer:   a,
	}, nil
}
