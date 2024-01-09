package page

// PageParam Paging parameters
type PageParam struct {
	CurrentPage int                    `form:"current_page" json:"current_page" bson:"current_page" xml:"current_page" yaml:"current_page"`
	PageSize    int                    `form:"page_size" json:"page_size" bson:"page_size" xml:"page_size" yaml:"page_size"`
	Total       int64                  `form:"total" json:"total" bson:"total" xml:"total" yaml:"total"`
	Filter      []string               `form:"filter[]" json:"filter[]" bson:"filter[]" xml:"filter[]" yaml:"filter[]"`
	Order       string                 `form:"order" json:"order" bson:"order" xml:"order" yaml:"order"`
	Select      []string               `form:"select[]" json:"select[]" bson:"select[]" xml:"select[]" yaml:"select[]"`
	Group       string                 `form:"group" json:"group" bson:"group" xml:"group" yaml:"group"`
	Omit        string                 `form:"omit" json:"omit" bson:"omit" xml:"omit" yaml:"omit"`
	Extra       map[string]interface{} `form:"extra" json:"extra" bson:"extra" xml:"extra" yaml:"extra"`
}

// New Default pagination
func New() PageParam {
	return PageParam{
		CurrentPage: 1,
		PageSize:    10,
		Extra:       map[string]interface{}{},
	}
}
