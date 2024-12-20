package sql

import (
	"fmt"
	"github.com/small-ek/antgo/utils/conv"
	"github.com/small-ek/antgo/utils/page"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

// Like Fuzzy search when there is value
func Like(key, value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if key != "" && value != "" {
			return db.Where(key+" LIKE ?", value+"%")
		}
		return db
	}
}

// Ilike Fuzzy search when there is value
func Ilike(key, value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if key != "" && value != "" {
			return db.Where(key+" ILIKE ?", value+"%")
		}
		return db
	}
}

// WhereIn WhereIn search when there is value
func WhereIn(key string, value interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch v := value.(type) {
		case []interface{}, []int, []int16, []int32, []int64, []uint16, []uint32, []uint64, []string, []float32, []float64:
			newValue := conv.Interfaces(v)
			if len(newValue) == 0 {
				return db
			}
		}

		if key != "" && value != nil && value != "" {
			return db.Where(fmt.Sprintf("%s IN ?", key), value)
		}
		return db
	}
}

// WhereNotIn WhereNotIn search when there is value
func WhereNotIn(key string, value interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch v := value.(type) {
		case []interface{}, []int, []int16, []int32, []int64, []uint16, []uint32, []uint64, []string, []float32, []float64:
			newValue := conv.Interfaces(v)
			if len(newValue) == 0 {
				return db
			}
		}

		if key != "" && value != nil && value != "" {
			return db.Where(""+key+" NOT IN ?", value)
		}
		return db
	}
}

// Where Where to search when there is value
func Where(key, conditions string, value interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		conditions = strings.ToUpper(conditions)
		if key == "" || conditions == "" || value == nil || value == "" || value == 0 {
			return db
		}

		switch v := value.(type) {
		case []interface{}, []int, []int16, []int32, []int64, []uint16, []uint32, []uint64, []string, []float32, []float64:
			newValue := conv.Interfaces(v)
			if (conditions == "BETWEEN" || conditions == "NOT BETWEEN") && len(newValue) == 2 {
				return db.Where(fmt.Sprintf("%s %s ? AND ?", key, conditions), newValue[0], newValue[1])
			}
			if (conditions == "IN" || conditions == "NOT IN") && len(newValue) > 0 {
				return db.Where(fmt.Sprintf("%s %s ?", key, conditions), newValue)
			}
		case string:
			if (conditions == "LIKE" || conditions == "ILIKE") && v != "" {
				return db.Where(fmt.Sprintf("%s %s ?", key, conditions), v+"%")
			} else {
				return db.Where(fmt.Sprintf("%s %s ?", key, conditions), value)
			}
		default:
			return db.Where(fmt.Sprintf("%s %s ?", key, conditions), value)
		}

		return db
	}
}

// Order Sort Can prevent injection sorting
func Order(str []string, desc []bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(str) > 0 && len(desc) > 0 && len(str) == len(desc) {
			for i := 0; i < len(str); i++ {
				db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: str[i]}, Desc: desc[i]})
			}
		}
		return db
	}
}

// Paginate 分页查询,默认最大10000，最大值可自定义
func Paginate(pageSize, currentPage int, maxSize ...int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(maxSize) > 0 && pageSize > maxSize[0] {
			pageSize = maxSize[0]
		} else if pageSize > 10000 {
			pageSize = 10000
		}

		return db.Limit(pageSize).Offset((currentPage - 1) * pageSize)
	}
}

// OnlyTrashed 显示软删除数据
func OnlyTrashed(res bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if res == true {
			return db.Unscoped().Where("deleted_at IS NOT NULL")
		} else {
			return db
		}
	}
}

// Filters constructs query filters
func Filters(filters interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var newFilter []page.Filter

		switch filter := filters.(type) {
		case string:
			if err := conv.UnmarshalJSON([]byte(filter), &newFilter); err != nil {
				return db
			}
		case []page.Filter:
			newFilter = filter
		default:
			return db
		}

		if len(newFilter) == 0 {
			return db
		}

		// Build WHERE conditions and arguments
		query, args := buildWhere(newFilter, "AND")
		return db.Where(query, args...)
	}
}

// buildWhere recursively constructs WHERE clause with AND/OR conditions
func buildWhere(filters []page.Filter, joinType string) (string, []interface{}) {
	var conditions []string
	var values []interface{}
	// Estimate the number of conditions to optimize slice allocation
	conditions = make([]string, 0, len(filters)*2) // Allocate space for conditions and subqueries

	for _, filter := range filters {
		// Handle OR conditions recursively
		if len(filter.Or) > 0 {
			subQuery, subValues := buildWhere(filter.Or, "OR")
			conditions = append(conditions, fmt.Sprintf("(%s)", subQuery))
			values = append(values, subValues...)
		}

		// Handle AND conditions recursively
		if len(filter.And) > 0 {
			subQuery, subValues := buildWhere(filter.And, "AND")
			conditions = append(conditions, fmt.Sprintf("(%s)", subQuery))
			values = append(values, subValues...)
		}

		// Handle basic filter conditions
		if filter.Field != "" && filter.Operator != "" && filter.Value != nil && isValidOperator(filter.Operator) {
			operator := strings.ToUpper(filter.Operator)
			condition, value := handleOperator(filter, operator)
			if condition != "" {
				conditions = append(conditions, condition)
				values = append(values, value...)
			}
		}
	}

	query := strings.Join(conditions, " "+joinType+" ")
	return query, values
}

// handleOperator handles different operator types to generate query conditions
func handleOperator(filter page.Filter, operator string) (string, []interface{}) {
	var condition string
	var values []interface{}

	switch operator {
	case "BETWEEN":
		// Handle BETWEEN operator
		value := conv.Strings(filter.Value)
		if len(value) == 2 {
			condition = fmt.Sprintf("%s BETWEEN ? AND ?", filter.Field)
			values = append(values, value[0], value[1])
		}
	case "IS NULL", "IS NOT NULL":
		// Handle IS NULL or IS NOT NULL
		condition = fmt.Sprintf("%s %s", filter.Field, operator)
	case "LIKE", "NOT LIKE":
		// Handle LIKE or NOT LIKE with escaping
		condition = fmt.Sprintf("%s %s ?", filter.Field, operator)
		values = append(values, fmt.Sprintf("%s%%", filter.Value))
	case "IN", "NOT IN":
		// Handle IN or NOT IN with values
		condition = fmt.Sprintf("%s %s (?)", filter.Field, operator)
		values = append(values, conv.Strings(filter.Value))
	default:
		// Handle other operators
		condition = fmt.Sprintf("%s %s ?", filter.Field, operator)
		values = append(values, filter.Value)
	}

	return condition, values
}

// isValidOperator checks if the operator is valid
func isValidOperator(operator string) bool {
	// Direct lookup for valid operators using map
	validOperators := map[string]bool{
		"=": true, ">": true, ">=": true, "<": true, "<=": true, "!=": true,
		"<>": true, "IN": true, "NOT IN": true, "LIKE": true, "NOT LIKE": true,
		"ILIKE": true, "RLIKE": true, "BETWEEN": true, "IS NULL": true, "IS NOT NULL": true,
	}
	return validOperators[strings.ToUpper(operator)]
}
