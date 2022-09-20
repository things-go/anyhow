package sms

type Dummy struct{}

func NewDummy() *Dummy                                      { return &Dummy{} }
func (*Dummy) Name() string                                 { return "DummySms " }
func (*Dummy) SendMessage(string, SendMessageRequest) error { return nil }
