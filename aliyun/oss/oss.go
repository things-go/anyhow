package oss

import (
	"errors"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
)

var ErrInvalidKey = errors.New("invalid key")

type Config struct {
	// required
	AccessKeyId string `json:"accessKeyId" yaml:"accessKeyId" binding:"required"`
	// required
	AccessKeySecret string `json:"accessKeySecret" yaml:"accessKeySecret" binding:"required"`
	// 角色权限
	// required
	RoleArn string `json:"roleArn" yaml:"roleArn" binding:"required"`
	// 区域: oss-cn-guangzhou
	// required
	Region string `yaml:"region" json:"region" binding:"required"`
	// 外网端点: 例如 https://oss-cn-shenzhen.aliyuncs.com
	// required
	Endpoint string `json:"endpoint" yaml:"endpoint" binding:"required"`
	// 内网端点: 例如 https://oss-cn-shenzhen-internal.aliyuncs.com
	// required
	InternalEndpoint string `json:"internalEndpoint" yaml:"internalEndpoint" binding:"required"`
	// 桶名称
	// required
	BucketName string `json:"bucketName" yaml:"bucketName" binding:"required"`
	// cdn: cdn 域名地址
	CdnDomain string `yaml:"cdnDomain" json:"cdnDomain"`
	// 根目录
	Root string `json:"root" yaml:"root"`
	// 临时目录
	Temporary string `json:"temporary" yaml:"temporary"`
}

// GenerateObjectKey 生成对象的key
// 格式: {root}/{dir}/{uuid}[.{suffix}]
func (c *Config) GenerateObjectKey(dir, suffix string) string {
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

// CdnUrl 获取 key 的 cdn url
func (c *Config) CdnUrl(key string) string {
	if key == "" || c.CdnDomain == "" {
		return ""
	}
	return c.CdnDomain + "/" + key
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
