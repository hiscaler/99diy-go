package diy99

import (
	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/99diy-go/config"
)

type service struct {
	config     *config.Config // Config
	logger     Logger         // Logger
	httpClient *resty.Client  // HTTP client
}

// API Services
type services struct {
	WebImageEditor webImageEditorService // Web 图片编辑器
}
