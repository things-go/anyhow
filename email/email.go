package email

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/smtp"
	"strconv"

	"github.com/jordan-wright/email"
)

type Email = email.Email

type Config struct {
	Enabled   bool   `yaml:"enabled" json:"enabled"`
	Identity  string `yaml:"identity" json:"identity"`
	Host      string `yaml:"host" json:"host"`
	Port      int    `yaml:"port" json:"port"`
	Username  string `yaml:"username" json:"username"`
	Password  string `yaml:"password" json:"password"`
	From      string `yaml:"from" json:"from"`
	FromName  string `yaml:"fromName" json:"fromName"`
	Subject   string `yaml:"subject" json:"subject"`
	EnableSSL bool   `yaml:"enableSsl" json:"enableSsl"`
	CertFile  string `yaml:"certFile" json:"certFile"` // not used
	KeyFile   string `yaml:"keyFile" json:"keyFile"`   // not used
}

type Client struct {
	From      string
	Subject   string
	addr      string
	auth      smtp.Auth
	enableSSL bool
	tlsConfig *tls.Config
	tpl       *template.Template
}

// New creates a new email client.
func New(tpl *template.Template, c Config) (*Client, error) {
	if tpl == nil {
		return nil, errors.New("email: template is nil")
	}
	identity := c.Identity
	if identity == "" {
		identity = c.Username
	}
	from := c.From
	if c.FromName != "" {
		from = fmt.Sprintf("%s <%s>", c.FromName, c.From)
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         c.Host,
		Certificates:       nil,
	}
	if c.EnableSSL {
		if c.CertFile != "" {
			cert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
			if err != nil {
				return nil, fmt.Errorf("could not load key or cert file: %w", err)
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}
	}

	return &Client{
		from,
		c.Subject,
		net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		smtp.PlainAuth(identity, c.Username, c.Password, c.Host),
		c.EnableSSL,
		tlsConfig,
		tpl,
	}, nil
}

// Send sends the email.
func (sf *Client) Send(e *email.Email) error {
	if sf.enableSSL {
		return e.SendWithTLS(sf.addr, sf.auth, sf.tlsConfig)
	}
	return e.Send(sf.addr, sf.auth)
}

// SendEmail sends the email.
func (sf *Client) SendEmail(tpl string, data any, to ...string) error {
	if len(to) == 0 {
		return errors.New("email: to must not be empty")
	}
	buf := new(bytes.Buffer)
	err := sf.tpl.ExecuteTemplate(buf, tpl, data)
	if err != nil {
		return err
	}

	e := email.NewEmail()
	e.From = sf.From
	e.To = to
	e.Subject = sf.Subject
	e.HTML = buf.Bytes()
	return sf.Send(e)
}

// ExecuteTemplate executes the given template with data.
func (sf *Client) ExecuteTemplate(name string, data any) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := sf.tpl.ExecuteTemplate(buf, name, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// NewEmail creates a new email.
func NewEmail() *Email { return email.NewEmail() }

// NewEmailFromReader creates a new email from a reader.
func NewEmailFromReader(r io.Reader) (*Email, error) { return email.NewEmailFromReader(r) }
