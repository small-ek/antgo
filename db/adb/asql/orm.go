package asql

import (
	"database/sql/driver"
	"reflect"
	"strings"

	"github.com/small-ek/antgo/utils/conv"
	"github.com/small-ek/antgo/utils/page"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type columnCondition struct {
	Column   clause.Column
	Operator string
	Values   []interface{}
}

func (c columnCondition) Build(builder clause.Builder) {
	operator := dialectOperator(builder, c.Operator)

	builder.WriteQuoted(c.Column)
	builder.WriteByte(' ')
	builder.WriteString(operator)

	switch operator {
	case "BETWEEN", "NOT BETWEEN":
		builder.WriteByte(' ')
		builder.AddVar(builder, c.Values[0])
		builder.WriteString(" AND ")
		builder.AddVar(builder, c.Values[1])
	case "IS NULL", "IS NOT NULL":
		return
	default:
		builder.WriteByte(' ')
		builder.AddVar(builder, c.Values...)
	}
}

// PostgreSQL uses POSIX regex operators; MySQL uses RLIKE. Other dialects keep
// the caller-provided operator so unsupported SQL fails visibly instead of being
// silently rewritten to a different meaning.
func dialectOperator(builder clause.Builder, operator string) string {
	if operator != "RLIKE" {
		return operator
	}
	if stmt, ok := builder.(*gorm.Statement); ok && stmt.DB != nil && stmt.DB.Dialector != nil {
		if stmt.DB.Dialector.Name() == "postgres" {
			return "~"
		}
	}
	return operator
}

func columnExpr(key, operator string, values ...interface{}) clause.Expression {
	return columnCondition{
		Column:   clause.Column{Name: key},
		Operator: strings.ToUpper(strings.TrimSpace(operator)),
		Values:   values,
	}
}

func pattern(value interface{}) string {
	return strings.TrimRight(conv.String(value), "%") + "%"
}

func isZeroFilterValue(value interface{}) bool {
	if value == nil {
		return true
	}
	if valuer, ok := value.(driver.Valuer); ok {
		v, err := valuer.Value()
		return err == nil && isZeroFilterValue(v)
	}

	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Chan, reflect.Map, reflect.Slice:
		return rv.IsNil() || rv.Len() == 0
	case reflect.Func, reflect.Interface, reflect.Ptr:
		return rv.IsNil()
	case reflect.Array:
		return rv.Len() == 0
	case reflect.String:
		return rv.Len() == 0
	default:
		return rv.IsZero()
	}
}

func sliceValues(value interface{}) []interface{} {
	if value == nil {
		return nil
	}
	if values, ok := value.([]interface{}); ok {
		return values
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil
	}

	values := make([]interface{}, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		values = append(values, rv.Index(i).Interface())
	}
	return values
}

// Like Fuzzy search when there is value
func Like(key, value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if key != "" && value != "" {
			return db.Where(columnExpr(key, "LIKE", pattern(value)))
		}
		return db
	}
}

// Ilike Fuzzy search when there is value
func Ilike(key, value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if key != "" && value != "" {
			return db.Where(columnExpr(key, "ILIKE", pattern(value)))
		}
		return db
	}
}

// WhereIn WhereIn search when there is value
func WhereIn(key string, value interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if key != "" && !isZeroFilterValue(value) {
			return db.Where(columnExpr(key, "IN", value))
		}
		return db
	}
}

// WhereNotIn WhereNotIn search when there is value
func WhereNotIn(key string, value interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if key != "" && !isZeroFilterValue(value) {
			return db.Where(columnExpr(key, "NOT IN", value))
		}
		return db
	}
}

// Where Where to search when there is value
func Where(key, conditions string, value interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		conditions = strings.ToUpper(strings.TrimSpace(conditions))
		if key == "" || conditions == "" || !isValidOperator(conditions) {
			return db
		}
		if conditions == "IS NULL" || conditions == "IS NOT NULL" {
			return db.Where(columnExpr(key, conditions))
		}
		if isZeroFilterValue(value) {
			return db
		}

		switch conditions {
		case "BETWEEN", "NOT BETWEEN":
			values := sliceValues(value)
			if len(values) == 2 {
				return db.Where(columnExpr(key, conditions, values[0], values[1]))
			}
		case "IN", "NOT IN":
			values := sliceValues(value)
			if len(values) > 0 {
				return db.Where(columnExpr(key, conditions, values))
			}
		case "LIKE", "NOT LIKE", "ILIKE", "NOT ILIKE":
			if v, ok := value.(string); ok {
				return db.Where(columnExpr(key, conditions, pattern(v)))
			}
			return db.Where(columnExpr(key, conditions, value))
		default:
			return db.Where(columnExpr(key, conditions, value))
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

		query := buildWhere(newFilter, "AND")
		if query == nil {
			return db
		}
		return db.Where(query)
	}
}

// buildWhere recursively constructs WHERE clause with AND/OR conditions
func buildWhere(filters []page.Filter, joinType string) clause.Expression {
	conditions := make([]clause.Expression, 0, len(filters)*2)

	for _, filter := range filters {
		if len(filter.Or) > 0 {
			if subQuery := buildWhere(filter.Or, "OR"); subQuery != nil {
				conditions = append(conditions, subQuery)
			}
		}

		if len(filter.And) > 0 {
			if subQuery := buildWhere(filter.And, "AND"); subQuery != nil {
				conditions = append(conditions, subQuery)
			}
		}

		operator := strings.ToUpper(strings.TrimSpace(filter.Operator))
		isNullOperator := operator == "IS NULL" || operator == "IS NOT NULL"
		if filter.Field != "" && operator != "" && isValidOperator(operator) && (isNullOperator || !isZeroFilterValue(filter.Value)) {
			if condition := handleOperator(filter, operator); condition != nil {
				conditions = append(conditions, condition)
			}
		}
	}

	if len(conditions) == 0 {
		return nil
	}
	if strings.ToUpper(joinType) == "OR" {
		return clause.Or(conditions...)
	}
	return clause.And(conditions...)
}

// handleOperator handles different operator types to generate query conditions
func handleOperator(filter page.Filter, operator string) clause.Expression {
	switch operator {
	case "BETWEEN", "NOT BETWEEN":
		values := sliceValues(filter.Value)
		if len(values) == 2 {
			return columnExpr(filter.Field, operator, values[0], values[1])
		}
	case "IS NULL", "IS NOT NULL":
		return columnExpr(filter.Field, operator)
	case "LIKE", "NOT LIKE", "ILIKE", "NOT ILIKE":
		return columnExpr(filter.Field, operator, pattern(filter.Value))
	case "IN", "NOT IN":
		values := sliceValues(filter.Value)
		if len(values) > 0 {
			return columnExpr(filter.Field, operator, values)
		}
	default:
		return columnExpr(filter.Field, operator, filter.Value)
	}

	return nil
}

// isValidOperator checks if the operator is valid
func isValidOperator(operator string) bool {
	validOperators := map[string]bool{
		"=": true, ">": true, ">=": true, "<": true, "<=": true, "!=": true,
		"<>": true, "IN": true, "NOT IN": true, "LIKE": true, "NOT LIKE": true,
		"ILIKE": true, "NOT ILIKE": true, "RLIKE": true, "BETWEEN": true, "NOT BETWEEN": true,
		"IS NULL": true, "IS NOT NULL": true,
	}
	return validOperators[strings.ToUpper(strings.TrimSpace(operator))]
}
