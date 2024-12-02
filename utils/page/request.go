package page

// PageParam Paging parameters 分页参数
type PageParam struct {
	CurrentPage int            `form:"current_page" json:"current_page" bson:"current_page" xml:"current_page" yaml:"current_page"` //当前页
	PageSize    int            `form:"page_size" json:"page_size" bson:"page_size" xml:"page_size" yaml:"page_size"`                //每页显示数量
	Filter      []Filter       `form:"filter[]" json:"filter[]" bson:"filter[]" xml:"filter[]" yaml:"filter[]"`                     //过滤条件
	FilterMap   map[string]any `form:"filter_map" json:"filter_map" bson:"filter_map" xml:"filter_map" yaml:"filter_map"`           //过滤条件形式2
	Omit        string         `form:"omit" json:"omit" bson:"omit" xml:"omit" yaml:"omit"`                                         //忽略字段
	Extra       map[string]any `form:"extra" json:"extra" bson:"extra" xml:"extra" yaml:"extra"`                                    //额外参数
	Order       []string       `form:"order[]" json:"order[]" bson:"order[]" xml:"order[]" yaml:"order[]"`                          //排序字段
	Desc        []bool         `form:"desc[]" json:"desc[]" bson:"desc[]" xml:"desc[]" yaml:"desc[]"`                               //排序方式 true 降序 false 升序
}

// Filter 过滤条件
type Filter struct {
	Field    string   `form:"field" json:"field" bson:"field" xml:"field" yaml:"field"`                //字段
	Operator string   `form:"operator" json:"operator" bson:"operator" xml:"operator" yaml:"operator"` //操作符
	Value    any      `form:"value" json:"value" bson:"value" xml:"value" yaml:"value"`                //值
	And      []Filter `form:"and" json:"and" bson:"and" xml:"and" yaml:"and"`                          //处理And条件
	Or       []Filter `form:"or" json:"or" bson:"or" xml:"or" yaml:"or"`                               //处理Or条件
}

// New Default pagination 默认分页参数,查询的时候可以直接使用
func New() PageParam {
	return PageParam{
		CurrentPage: 1,
		PageSize:    10,
	}
}
