package category

type ConsumerMessageV1 struct {
	Uuid    string `json:"uuid"`
	Subject string `json:"subject"`
	Event   string `json:"event"`
	Version string `json:"version"`
	Payload struct {
		Id          string `json:"id"`
		ParentId    string `json:"parent_id,omitempty"`
		IsActive    bool   `json:"is_active"`
		Name        string `json:"name"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Img         string `json:"img,omitempty"`
	} `json:"payload"`
}
