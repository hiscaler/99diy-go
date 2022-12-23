package diy99

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type webImageEditorService service

const (
	ImageType = "image"
	TextType  = "text"
)

type WebImageEditorOrderItemData struct {
	Type    string `json:"type"`
	Font    string `json:"font"`
	Color   string `json:"color"`
	Content string `json:"content"`
	URL     string `json:"url"`
}

type WebImageEditorOrderItem struct {
	OrderNumber           string                        `json:"orderNumber"` // 订单号
	OrderKey              string                        `json:"key"`         // 订单项 ID
	TemplateId            int                           `json:"templetId"`   // 模板
	PreviewViewPictureURL string                        `json:"preViewPic"`  // 预览图
	CallbackURL           string                        `json:"callBackUrl"` // 回调地址
	Data                  []WebImageEditorOrderItemData `json:"data"`        // 订单项数据
}

type WebImageEditorOrderRequest struct {
	Items []WebImageEditorOrderItem `json:"items"`
}

func (m WebImageEditorOrderRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Items, validation.Required.Error("订单项不能为空")),
	)
}

// PushOrders 推送订单
func (s webImageEditorService) PushOrders(req WebImageEditorOrderRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	res := struct {
		Data interface{} `json:"data"`
	}{}
	resp, err := s.httpClient.R().
		SetBody(req.Items).
		Post("/order/createOrders")
	if err != nil {
		return err
	}

	if err = json.Unmarshal(resp.Body(), &res); err != nil {
		return err
	}
	return nil
}
