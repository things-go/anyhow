package oss

import (
	"errors"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
)

var ErrInvalidKey = errors.New("invalid key")

type Config struct {
	AccessKeyId     string `json:"access_key_id" yaml:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret" yaml:"access_key_secret"`
	// 角色权限
	RoleArn string `json:"role_arn" yaml:"role_arn"`
	// 外网端点: 例如 https://oss-cn-shenzhen.aliyuncs.com
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	// 内网端点: 例如 https://oss-cn-shenzhen-internal.aliyuncs.com
	InternalEndpoint string `json:"internal_endpoint" yaml:"internal_endpoint"`
	// 桶名称
	BucketName string `json:"bucket_name" yaml:"bucket_name"`
	// 根目录
	Root string `json:"root" yaml:"root"`
	// 临时目录
	Temporary string `json:"temporary" yaml:"temporary"`
}

type Client struct {
	Config
	outer *oss.Client // 外网实例
	inner *oss.Client // 内网实例
}

func New(c Config) (*Client, error) {
	client, err := oss.New(c.Endpoint, c.AccessKeyId, c.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	innerClient, err := oss.New(c.InternalEndpoint, c.AccessKeyId, c.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	return &Client{
		Config: c,
		outer:  client,
		inner:  innerClient,
	}, nil
}

// R 获取oss客户端,
// isIntranet: 通过哪个端点通道传输
//
//	true: 内网
//	false: 外网
func (c *Client) R(isIntranet bool) *Instance {
	client := c.outer
	if isIntranet {
		client = c.inner
	}
	return &Instance{
		&c.Config,
		client,
	}
}

// GenerateObjectKey 生成对象的key
// 格式: {root}/{dir}/{uuid}[.{suffix}]
func (c *Client) GenerateObjectKey(dir, suffix string) string {
	uq := uuid.New().String()
	b := strings.Builder{}
	b.Grow(len(c.Root) + len(dir) + len(uq) + len(suffix) + 3)

	if c.Root != "" {
		b.WriteString(c.Root)
		b.WriteString("/")
	}
	b.WriteString(dir)
	b.WriteString("/")
	b.WriteString(uq)
	if suffix != "" {
		b.WriteString(".")
		b.WriteString(suffix)
	}
	return b.String()
}
