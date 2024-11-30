package property

type ConsumerMessageV1 struct {
	Uuid    string `json:"uuid"`
	Subject string `json:"subject"`
	Event   string `json:"event"`
	Version string `json:"version"`
	Payload struct {
		Id                string `json:"id"`
		IsActive          bool   `json:"is_active"`
		Name              string `json:"name"`
		Code              string `json:"code"`
		Hint              string `json:"hint,omitempty"`
		Type              string `json:"type"`
		ShowFilter        bool   `json:"show_filter"`
		ShowProductList   bool   `json:"show_product_list"`
		ShowDetailProduct bool   `json:"show_detail_product"`
	} `json:"payload"`
}
