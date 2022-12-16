package diy99

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/99diy-go/config"
	"github.com/hiscaler/gox/bytex"
	"strings"
	"time"
)

const (
	Version   = "0.0.1"
	userAgent = "99Diy API Client-Golang/" + Version + " (https://github.com/hiscaler/99diy-go)"
)

const (
	OK                   = 200 // 无错误
	ServiceNotFoundError = 400 // 服务不存在
	InternalError        = 500 // 内部错误，数据库异常
)

type Diy99 struct {
	config     *config.Config // 配置
	logger     Logger
	httpClient *resty.Client // Resty Client
	forceToken bool          // 强制获取 Token
	Services   services      // API Services
}

func NewDiy99(cfg config.Config) *Diy99 {
	diy99Client := &Diy99{
		config: &cfg,
		logger: createLogger(),
	}
	httpClient := resty.
		New().
		SetDebug(diy99Client.config.Debug).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
			"User-Agent":   userAgent,
		})
	if cfg.Sandbox {
		httpClient.SetBaseURL("http://8.129.218.196:9199")
	} else {
		httpClient.SetBaseURL("http://8.129.218.196:9199")
	}

	httpClient.
		SetTimeout(time.Duration(cfg.Timeout) * time.Second).
		OnAfterResponse(func(client *resty.Client, response *resty.Response) (err error) {
			if response.IsError() {
				return fmt.Errorf("%s: %s", response.Status(), bytex.ToString(response.Body()))
			}

			r := struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{}
			if err = json.Unmarshal(response.Body(), &r); err == nil {
				err = ErrorWrap(r.Code, r.Message)
			}
			return
		}).
		SetRetryCount(2).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second)
	diy99Client.httpClient = httpClient
	xService := service{
		config:     &cfg,
		logger:     diy99Client.logger,
		httpClient: diy99Client.httpClient,
	}
	diy99Client.Services = services{
		WebImageEditor: (webImageEditor)(xService),
	}
	return diy99Client
}

// SetDebug 设置是否开启调试模式
func (diy *Diy99) SetDebug(v bool) *Diy99 {
	diy.config.Debug = v
	diy.httpClient.SetDebug(v)
	return diy
}

// SetLogger 设置日志器
func (diy *Diy99) SetLogger(logger Logger) *Diy99 {
	diy.logger = logger
	return diy
}

// ErrorWrap 错误包装
func ErrorWrap(code int, message string) error {
	if code == OK || code == 0 {
		return nil
	}

	switch code {
	case ServiceNotFoundError:
		message = "服务不存在"
	default:
		if code == InternalError {
			if message == "" {
				message = "内部错误，请联系 99Diy"
			}
		} else {
			message = strings.TrimSpace(message)
			if message == "" {
				message = "未知错误"
			}
		}
	}
	return fmt.Errorf("%d: %s", code, message)
}
