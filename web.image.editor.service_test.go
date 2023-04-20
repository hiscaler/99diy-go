package diy99

import "testing"

func Test_webImageEditorService_PushOrders(t *testing.T) {
	tests := []struct {
		name    string
		req     WebImageEditorOrderRequest
		wantErr bool
	}{
		{"tag1", WebImageEditorOrderRequest{
			Items: []WebImageEditorOrderItem{
				{
					OrderNumber:           "PO001",
					OrderKey:              "1",
					TemplateId:            3,
					PreviewViewPictureURL: "https://www.example.com/1.jpg",
					CallbackURL:           "https://api.example.com/callback",
					Data: []WebImageEditorOrderItemData{
						{Type: TextType, Font: "Arial", Color: "#000000", Content: "Hello"},
						{Type: ImageType, URL: "https://www.example.com/2.jpg"},
					},
				},
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := diy99Client.Services.WebImageEditor.PushOrders(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("PushOrders() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
