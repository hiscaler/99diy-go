package diy99

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"regexp"
)

type webImageEditorService service

const (
	ImageType                = "image"        // 图片
	TextType                 = "text"         // 文本
	BackgroundColorType      = "bgColor"      // 背景颜色
	BackgroundImageType      = "bgImage"      // 背景图
	BackgroundColorImageType = "bgColorImage" // 背景颜色图片
)

type WebImageEditorOrderItemData struct {
	Type      string `json:"type"`      // 资源类型
	Font      string `json:"font"`      // 字体
	Color     string `json:"color"`     // 字体颜色
	Content   string `json:"content"`   // 定制内容
	URL       string `json:"url"`       // 图片地址
	BgColor   string `json:"bgColor"`   // 背景颜色
	Direction string `json:"direction"` // 当前资源的字段名称
}

func (m WebImageEditorOrderItemData) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Type, validation.Required.Error("资源类型不能为空"), validation.In(ImageType, TextType, BackgroundColorType, BackgroundImageType, BackgroundColorImageType).Error("无效的资源类型")),
		validation.Field(&m.Font, validation.When(m.Type == TextType, validation.Required.Error("字体不能为空"))),
		validation.Field(&m.Color, validation.When(m.Type == TextType, validation.Required.Error("字体颜色不能为空"),
			validation.WithContext(func(ctx context.Context, value interface{}) error {
				s, ok := value.(string)
				if !ok {
					return errors.New("无效的字体颜色值")
				}
				if !regexp.MustCompile("^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$").MatchString(s) {
					return fmt.Errorf("无效的字体颜色值: %s", s)
				}
				return nil
			}),
		)),

		validation.Field(&m.Content, validation.When(m.Type == TextType, validation.Required.Error("定制内容不能为空"))),
		validation.Field(&m.URL, validation.When(m.Type == ImageType || m.Type == BackgroundColorImageType, validation.Required.Error("图片地址不能为空"), is.URL.Error("无效的图片地址"))),
		validation.Field(&m.BgColor, validation.When(m.Type == BackgroundColorType, validation.Required.Error("背景颜色值不能为空"))),
	)
}

type WebImageEditorOrderItem struct {
	OrderNumber           string                        `json:"orderNumber"` // 订单号
	OrderKey              string                        `json:"key"`         // 订单项 ID
	TemplateId            int                           `json:"templetId"`   // 模板
	PreviewViewPictureURL string                        `json:"preViewPic"`  // 预览图
	CallbackURL           string                        `json:"callBackUrl"` // 回调地址
	Data                  []WebImageEditorOrderItemData `json:"data"`        // 订单项数据
	State                 map[string]interface{}        `json:"state"`       // 附带信息
	Editor                struct {
		SKU                string  `json:"sku"`
		Width              float64 `json:"width"`
		Height             float64 `json:"height"`
		DPI                int     `json:"dpi"`
		Remark             string  `json:"remark"`
		ExtraCustomizeInfo string  `json:"extraCustomizeInfo"` // 额外的定制信息
	} `json:"editor"`
}

func (m WebImageEditorOrderItem) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OrderNumber, validation.Required.Error("订单号不能为空")),
		validation.Field(&m.OrderKey, validation.Required.Error("订单项目 ID 不能为空")),
		validation.Field(&m.TemplateId, validation.Required.Error("无效的模板")),
		validation.Field(&m.PreviewViewPictureURL,
			validation.Required.Error("预览图地址不能为空"),
			is.URL.Error("无效的预览图地址"),
		),
		validation.Field(&m.CallbackURL,
			validation.Required.Error("回调地址不能为空"),
			is.URL.Error("无效的回调地址"),
		),
		validation.Field(&m.Data, validation.Required.Error("订单项数据不能为空")),
	)
}

type WebImageEditorOrderRequest struct {
	Items []WebImageEditorOrderItem `json:"items"`
}

func (m WebImageEditorOrderRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Items, validation.Required.Error("订单项不能为空")),
		validation.Field(&m.Items, validation.When(len(m.Items) > 0, validation.Each(validation.WithContext(func(ctx context.Context, value interface{}) error {
			if item, ok := value.(WebImageEditorOrderItem); !ok {
				return errors.New("无效的数据")
			} else {
				return item.Validate()
			}
		})))),
	)
}

// PushOrders 推送订单
func (s webImageEditorService) PushOrders(req WebImageEditorOrderRequest) (orderId int, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	res := struct {
		Result  bool        `json:"result"`
		Code    int         `json:"code"`
		Message *string     `json:"msg"`
		Data    interface{} `json:"data"` // OrderId or other
	}{}
	for i, item := range req.Items {
		if item.State == nil {
			req.Items[i].State = make(map[string]interface{}, 0)
		}
	}
	resp, err := s.httpClient.R().
		SetBody(req.Items).
		Post("/order/createOrders")
	if err != nil {
		return
	}

	if err = json.Unmarshal(resp.Body(), &res); err != nil {
		return
	}
	if res.Result == false {
		message := ""
		if res.Message != nil {
			message = *res.Message
		}
		code := res.Code
		if code == 0 {
			code = resp.StatusCode()
		}
		err = ErrorWrap(code, message)
		if err != nil {
			return
		}
	}

	switch res.Data.(type) {
	case int:
		return res.Data.(int), nil
	case float64:
		return int(res.Data.(float64)), nil
	default:
		return 0, errors.New("未返回有效的订单号")
	}
}
