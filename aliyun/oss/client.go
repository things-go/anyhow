package oss

import (
	"io"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OSS 实例()
type Instance struct {
	*Config
	*oss.Client
}

// GetObject 获取对象
func (c *Instance) GetObject(key string) (io.ReadCloser, error) {
	if key == "" {
		return nil, ErrInvalidKey
	}
	bucket, err := c.Bucket(c.BucketName)
	if err != nil {
		return nil, err
	}
	return bucket.GetObject(key)
}

// GetObjects 获取对象列表
func (c *Instance) GetObjects(key []string) ([]io.ReadCloser, error) {
	if len(key) < 1 {
		return nil, ErrInvalidKey
	}

	bucket, err := c.Bucket(c.BucketName)
	if err != nil {
		return nil, err
	}
	rc := make([]io.ReadCloser, 0, len(key))
	for _, k := range key {
		if k == "" {
			continue
		}
		var r io.ReadCloser
		r, err = bucket.GetObject(k)
		if err != nil {
			break
		}
		rc = append(rc, r)
	}
	if err != nil {
		// 发生错误, 关闭所有已打开的文件
		for _, r := range rc {
			r.Close()
		}
		return nil, err
	}
	return rc, nil
}

// Upload 上传
func (c *Instance) Upload(key string, r io.Reader) error {
	bucket, err := c.Bucket(c.BucketName)
	if err != nil {
		return err
	}
	return bucket.PutObject(key, r)
}

// 获取资源 url(临时url)
func (c *Instance) SignURL(key string, method oss.HTTPMethod, expiredInSec int64, options ...oss.Option) (string, error) {
	bucket, err := c.Bucket(c.BucketName)
	if err != nil {
		return "", err
	}
	return bucket.SignURL(key, method, expiredInSec, options...)
}

// 获取资源 url, GET 方法 (临时url)
func (c *Instance) SignURLWithGet(key string, expiredInSec int64, options ...oss.Option) (string, error) {
	return c.SignURL(key, oss.HTTPGet, expiredInSec, options...)
}

// 获取资源 url, GET 方法 (临时url)
func (c *Instance) SignAttachmentURLWithGet(key, filename string, expiredInSec int64) (string, error) {
	return c.SignURL(key, oss.HTTPGet, expiredInSec, oss.ResponseContentDisposition("attachment; filename="+filename))
}

// PutObjectToTemporary 存储到临时文件夹, 返回 object key 并设置 tags: temporary:true
// oss服务器对 {Temporary} 目录下, 设有 tags为 temporary:true 一天后删除
// Temporary 的文件夹名为 200601 格式
func (c *Instance) PutObjectToTemporary(filename string, reader io.Reader) (string, error) {
	bucket, err := c.Bucket(c.BucketName)
	if err != nil {
		return "", err
	}
	objectKey := c.Temporary + "/" + time.Now().Format("200601") + "/" + filename
	err = bucket.PutObject(objectKey,
		reader,
		oss.SetTagging(oss.Tagging{Tags: []oss.Tag{{Key: "temporary", Value: "true"}}}),
	)
	if err != nil {
		return "", err
	}
	return objectKey, nil
}
