package page

// PageParam Paging parameters
type PageParam struct {
	CurrentPage int                    `form:"current_page" json:"current_page" bson:"current_page" xml:"current_page" yaml:"current_page"` //当前页
	PageSize    int                    `form:"page_size" json:"page_size" bson:"page_size" xml:"page_size" yaml:"page_size"`                //每页显示数量
	Filter      []string               `form:"filter[]" json:"filter[]" bson:"filter[]" xml:"filter[]" yaml:"filter[]"`                     //过滤条件
	Select      []string               `form:"select[]" json:"select[]" bson:"select[]" xml:"select[]" yaml:"select[]"`                     //选择字段
	Omit        string                 `form:"omit" json:"omit" bson:"omit" xml:"omit" yaml:"omit"`                                         //忽略字段
	Extra       map[string]interface{} `form:"extra" json:"extra" bson:"extra" xml:"extra" yaml:"extra"`                                    //额外参数
}

// New Default pagination
func New() PageParam {
	return PageParam{
		CurrentPage: 1,
		PageSize:    10,
		Extra:       map[string]interface{}{},
	}
}
