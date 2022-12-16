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
	Content string `json:"content"`
	URL     string `json:"url"`
}

type WebImageEditorOrderItem struct {
	OrderNumber string                        `json:"order_number"`
	OrderKey    string                        `json:"order_key"`
	PreviewURL  string                        `json:"preview_url"`
	CallbackURL string                        `json:"callback_url"`
	Data        []WebImageEditorOrderItemData `json:"data"`
}

type WebImageEditorOrderRequest struct {
	CallbackURL string                    `json:"callback_url"`
	Items       []WebImageEditorOrderItem `json:"items"`
}

func (m WebImageEditorOrderRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CallbackURL, validation.Required.Error("回调地址不能为空")),
		validation.Field(&m.Items, validation.Required.Error("订单项不能为空")),
	)
}

// PushOrders 推送订单
func (s webImageEditorService) PushOrders(req WebImageEditorOrderRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	res := struct {
		Data struct{} `json:"data"`
	}{}
	resp, err := s.httpClient.R().
		SetBody(req).
		Post("/orders/createOrder")
	if err != nil {
		return err
	}

	if err = json.Unmarshal(resp.Body(), &res); err != nil {
		return err
	}
	return nil
}
