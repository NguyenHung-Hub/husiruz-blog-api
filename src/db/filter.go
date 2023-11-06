package db

type PostFilter struct {
	Status     string `json:"status" form:"status"`
	CategoryId string `json:"category_id" form:"category_id"`
}

func (f *PostFilter) Default() {
	if f.Status == "" {
		f.Status = string(PostStatusVisibility)
	}
}
