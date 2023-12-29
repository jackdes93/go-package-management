package common

type OrderBy struct {
	Key    string
	IsDesc bool
}

type Paging struct {
	Limit int `json:"limit" form:"limit"`
	Total int `json:"total" form:"total"`
	Page  int `json:"page" form:"page"`
}

func (p *Paging) FullFill() {
	if p.Limit <= 0 {
		p.Limit = 25
	}

	if p.Page <= 0 {
		p.Page = 1
	}
}
