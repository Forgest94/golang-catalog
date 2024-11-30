package product

type Product struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	Code          string  `json:"code"`
	Price         float32 `json:"price"`
	BasePrice     float32 `json:"base_price"`
	ParentProduct struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	}
	Categories []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	}
	Properties []struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		Code   string `json:"code"`
		Values []any  `json:"values"`
	}
}

type ConsumerMessageV1 struct {
	Uuid    string `json:"uuid"`
	Subject string `json:"subject"`
	Event   string `json:"event"`
	Version string `json:"version"`
	Payload struct {
		Id                 string   `json:"id"`
		ParentId           string   `json:"parent_id,omitempty"`
		IsActive           bool     `json:"is_active"`
		Name               string   `json:"name"`
		Code               string   `json:"code"`
		Description        string   `json:"description"`
		PreviewDescription string   `json:"preview_description"`
		Img                string   `json:"img"`
		Price              float32  `json:"price"`
		BasePrice          float32  `json:"base_price"`
		Categories         []string `json:"categories,omitempty"`
		Properties         []struct {
			Id     string `json:"id"`
			Values []any  `json:"values"`
		} `json:"properties,omitempty"`
	} `json:"payload"`
}
