package page

// PageParam 分页参数结构体，适用于GET请求
// PageParam is a struct for pagination parameters, typically used in GET requests.
type PageParam struct {
	CurrentPage int            `form:"current_page" json:"current_page"` // 当前页码，默认为1 / Current page number, default is 1
	PageSize    int            `form:"page_size" json:"page_size"`       // 每页显示的数据条数，默认为10 / Number of items per page, default is 10
	Filter      string         `form:"filter" json:"filter"`             // 复杂过滤条件，适用于复杂查询 / Complex filtering conditions, suitable for complex queries
	FilterMap   map[string]any `form:"filter_map" json:"filter_map"`     // 过滤条件的另一种形式，使用键值对表示 / Another form of filtering conditions, represented as key-value pairs
	Omit        string         `form:"omit" json:"omit"`                 // 需要忽略的字段 / Fields to be omitted
	Extra       map[string]any `form:"extra" json:"extra"`               // 额外的参数，用于传递其他自定义信息 / Additional parameters for passing custom information
	Order       []string       `form:"order[]" json:"order[]"`           // 排序字段列表 / List of fields to sort by
	Desc        []bool         `form:"desc[]" json:"desc[]"`             // 排序方式列表，true表示降序，false表示升序 / List of sorting orders, true for descending, false for ascending
}

// PageJson 分页参数结构体，适用于POST请求，支持更复杂的查询条件
// PageJson is a struct for pagination parameters, typically used in POST requests, supporting more complex query conditions.
type PageJson struct {
	CurrentPage int            `form:"current_page" json:"current_page" bson:"current_page" xml:"current_page" yaml:"current_page"` // 当前页码，默认为1 / Current page number, default is 1
	PageSize    int            `form:"page_size" json:"page_size" bson:"page_size" xml:"page_size" yaml:"page_size"`                // 每页显示的数据条数，默认为10 / Number of items per page, default is 10
	Filter      []Filter       `form:"filter" json:"filter" bson:"filter" xml:"filter" yaml:"filter"`                               // 复杂过滤条件列表 / List of complex filtering conditions
	FilterMap   map[string]any `form:"filter_map" json:"filter_map" bson:"filter_map" xml:"filter_map" yaml:"filter_map"`           // 过滤条件的另一种形式，使用键值对表示 / Another form of filtering conditions, represented as key-value pairs
	Omit        string         `form:"omit" json:"omit" bson:"omit" xml:"omit" yaml:"omit"`                                         // 需要忽略的字段 / Fields to be omitted
	Extra       map[string]any `form:"extra" json:"extra" bson:"extra" xml:"extra" yaml:"extra"`                                    // 额外的参数，用于传递其他自定义信息 / Additional parameters for passing custom information
	Order       []string       `form:"order" json:"order" bson:"order" xml:"order" yaml:"order"`                                    // 排序字段列表 / List of fields to sort by
	Desc        []bool         `form:"desc" json:"desc" bson:"desc" xml:"desc" yaml:"desc"`                                         // 排序方式列表，true表示降序，false表示升序 / List of sorting orders, true for descending, false for ascending
}

// Filter 过滤条件结构体，用于定义复杂的查询条件
// Filter is a struct for defining complex query conditions.
type Filter struct {
	Field    string   `form:"field" json:"field" bson:"field" xml:"field" yaml:"field"`                // 字段名 / Field name
	Operator string   `form:"operator" json:"operator" bson:"operator" xml:"operator" yaml:"operator"` // 操作符，如 "=", ">", "<" 等 / Operator, such as "=", ">", "<", etc.
	Value    any      `form:"value" json:"value" bson:"value" xml:"value" yaml:"value"`                // 字段值 / Field value
	And      []Filter `form:"and" json:"and" bson:"and" xml:"and" yaml:"and"`                          // AND 条件列表，用于组合多个过滤条件 / List of AND conditions for combining multiple filters
	Or       []Filter `form:"or" json:"or" bson:"or" xml:"or" yaml:"or"`                               // OR 条件列表，用于组合多个过滤条件 / List of OR conditions for combining multiple filters
}

// New 创建一个默认的分页参数实例，适用于GET请求
// New creates a default instance of PageParam, suitable for GET requests.
func New() PageParam {
	return PageParam{
		CurrentPage: 1,
		PageSize:    10,
	}
}

// NewJson 创建一个默认的分页参数实例，适用于POST请求
// NewJson creates a default instance of PageJson, suitable for POST requests.
func NewJson() PageJson {
	return PageJson{
		CurrentPage: 1,
		PageSize:    10,
	}
}
