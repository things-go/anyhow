package sts

import (
	"errors"
	"net/http"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
)

// https://help.aliyun.com/document_detail/66053.html?spm=a2c4g.11186623.2.12.151338dfkKJJKh
const (
	// STS服务的接入地址
	RegionId = "cn-hangzhou" // https://sts.aliyuncs.com/
	// 用户自定义参数.此参数用来区分不同的令牌,可用于用户级别的访问审计
	RoleSessionName = "token_server"
	// 权限策略,生成STS Token时可以指定一个额外的权限策略，以进一步限制STS Token的权限
	Policy = ""
	// 过期时间: 单位为秒
	//      最小值为900秒,
	//      最大值为 MaxSessionDuration 设置的时间. 默认值为3600秒。
	TokenExpireSeconds = 900
)

type Config struct {
	// required
	AccessKeyId string `yaml:"access_key_id" json:"access_key_id" binding:"required"` // 阿里的key
	// required
	AccessKeySecret string `yaml:"access_key_secret" json:"access_key_secret" binding:"required"` // 阿里的secret
	// required
	RoleArn string `yaml:"role_arn" json:"role_arn" binding:"required"` // 角色
	// required
	RegionId        string `yaml:"region_id" json:"region_id" binding:"required"` // STS服务的接入地址
	RoleSessionName string `yaml:"role_session_name" json:"role_session_name"`    // 用户自定义参数.此参数用来区分不同的令牌,可用于用户级别的访问审计
	Policy          string `yaml:"policy" json:"policy"`                          // 权限策略, 默认: 空
	// required
	Expiration int `yaml:"expiration" json:"expiration" binding:"gte=900"` // 过期时间(单位为秒),最小值为900秒, 最大值为MaxSessionDuration设置的时间. 默认值为3600秒。
}

type Client struct {
	*sts.Client
	Config
}

func NewClient(c Config) (*Client, error) {
	client, err := sts.NewClientWithAccessKey(c.RegionId, c.AccessKeyId, c.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	return &Client{
		Client: client,
		Config: c,
	}, nil
}

func (sf *Client) AssumeRole() (*sts.Credentials, error) {
	req := sts.CreateAssumeRoleRequest()
	req.Scheme = "https"
	req.RoleArn = sf.RoleArn
	req.RoleSessionName = sf.RoleSessionName
	req.Policy = sf.Policy
	req.DurationSeconds = requests.NewInteger(sf.Expiration)

	rsp, err := sf.Client.AssumeRole(req)
	if err != nil {
		return nil, err
	}

	if rsp.GetHttpStatus() != http.StatusOK {
		return nil, errors.New("获取sts令牌失败, 请稍后再试")
	}
	return &rsp.Credentials, nil
}
