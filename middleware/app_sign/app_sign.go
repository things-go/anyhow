package app_sign

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"

	"github.com/things-go/clip/signature"
)

// 签名验证过程
// 拼接签名串 str = timestamp + method + url + cipherBody + app_id
// sign = Base64(HMACSHA256(app_secret, str))
//
// 请求头:
//  AppId: 应用标识
//  Timestamp: 时间戳, 单位ms
//  Sign: 签名, 即上述 sign

// Option 签名选项
type Option func(*Config)

// Config 签名配置
type Config struct {
	availWindow   time.Duration
	skip          func(*gin.Context) bool
	errorFallback func(*gin.Context, error)
	appKeyPairs   func(c *gin.Context) (string, string, error)
}

func WithAppKeyPairs(f func(c *gin.Context) (string, string, error)) Option {
	return func(sc *Config) {
		if f != nil {
			sc.appKeyPairs = f
		}
		sc.appKeyPairs = f
	}
}

// WithAvailWindow 有效窗口时间, 小于等于0表示不验证
func WithAvailWindow(t time.Duration) Option {
	return func(o *Config) {
		o.availWindow = t
	}
}

// WithSkip 忽略验证的接口
func WithSkip(skip func(c *gin.Context) bool) Option {
	return func(o *Config) {
		if skip != nil {
			o.skip = skip
		}
	}
}

// WithUnauthorizedFallback sets the fallback handler when requests are unauthorized.
func WithErrorFallback(f func(c *gin.Context, err error)) Option {
	return func(o *Config) {
		if f != nil {
			o.errorFallback = f
		}
	}
}

// VerifySign 签名验证器
func VerifySign(opts ...Option) gin.HandlerFunc {
	cfg := Config{
		skip: func(c *gin.Context) bool { return false },
		errorFallback: func(c *gin.Context, err error) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "无效请求",
				"detail":  err.Error(),
			})
		},
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return func(c *gin.Context) {
		if cfg.skip(c) {
			c.Next()
			return
		}

		// 时间戳验证
		timestamp := c.GetHeader("Timestamp") // 时间戳
		if timestamp == "" {
			cfg.errorFallback(c, errors.New("无效的timestamp"))
			return
		}
		milli, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			cfg.errorFallback(c, errors.New("无效的timestamp"))
			return
		}
		if cfg.availWindow > 0 && time.Since(time.UnixMilli(milli)) > cfg.availWindow {
			cfg.errorFallback(c, errors.New("该请求已过期失效"))
			return
		}

		appId, appSecret, err := cfg.appKeyPairs(c)
		if err != nil {
			cfg.errorFallback(c, err)
			return
		}

		// body处理
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			cfg.errorFallback(c, err)
			return
		}
		c.Request.Body.Close()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		str := timestamp + strings.ToUpper(c.Request.Method) + c.Request.RequestURI + *(*string)(unsafe.Pointer(&body)) + appId
		calcSign := signature.HmacSha256(appSecret, str)
		sign := c.GetHeader("Sign")
		if calcSign != sign {
			cfg.errorFallback(c, errors.New("无效的签名"))
			return
		}
		c.Next()
	}
}
