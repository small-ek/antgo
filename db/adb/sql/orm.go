package sql

import (
	"encoding/json"
	"github.com/small-ek/antgo/utils/conv"
	"gorm.io/gorm"
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
			return db.Where(""+key+" IN (?)", value)
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
			return db.Where(""+key+" NOT IN (?)", value)
		}
		return db
	}
}

// Where Where to search when there is value
func Where(key, conditions string, value interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if key != "" && conditions != "" && value != "" && value != nil && value != 0 {
			switch v := value.(type) {
			case []interface{}, []int, []int16, []int32, []int64, []uint16, []uint32, []uint64, []string, []float32, []float64:
				newValue := conv.Interfaces(v)
				if (conditions == "BETWEEN" || conditions == "between") && len(newValue) == 2 {
					return db.Where(""+key+" "+conditions+" ? AND ?", newValue[0], newValue[1])
				}
				if (conditions == "NOT BETWEEN" || conditions == "not between") && len(newValue) == 2 {
					return db.Where(""+key+" "+conditions+" ? AND ?", newValue[0], newValue[1])
				}
				if (conditions == "IN" || conditions == "in") && len(newValue) > 0 {
					return db.Where(""+key+" IN (?)", value)
				}
				if (conditions == "NOT IN" || conditions == "not in") && len(newValue) > 0 {
					return db.Where(""+key+" NOT IN (?)", value)
				}
			default:
				return db.Where(""+key+" "+conditions+" ?", value)
			}

		}
		return db
	}
}

// Order Sort
func Order(str string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if str == "" {
			return db
		}
		return db.Order(str)
	}
}

// Paginate ...
func Paginate(pageSize, currentPage int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(pageSize).Offset((currentPage - 1) * pageSize)
	}
}

// OnlyTrashed ...
func OnlyTrashed(res bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if res == true {
			return db.Unscoped().Where("deleted_at IS NOT NULL")
		} else {
			return db
		}
	}
}

// Filters ...
func Filters(where interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if where == nil {
			return db
		}
		whereArray := where.([]string)
		for _, v := range whereArray {
			var arr []interface{}
			json.Unmarshal([]byte(v), &arr)
			db = buildWhere(arr, db)
		}
		return db
	}
}

// buildWhere
func buildWhere(arr []interface{}, db *gorm.DB) *gorm.DB {
	if len(arr) == 3 && arr[2] != "" && arr[2] != nil {
		db = and(arr[0].(string), arr[1].(string), arr[2], db)
	}
	if len(arr) == 4 && arr[3] != "" && arr[3] != nil && arr[0] == "or" {
		db = or(arr[1].(string), arr[2].(string), arr[3], db)
	}
	return db
}

// and
func and(key, condition string, value interface{}, db *gorm.DB) *gorm.DB {
	switch condition {
	case "like", "LIKE", "Like", "notlike", "NOTLIKE", "Notlike", "ilike", "ILIKE", "Ilike", "rlike", "RLIKE", "Rlike":
		db = db.Where(key+" "+condition+" ?", conv.String(value)+"%")
	case "in", "IN", "In", "not in", "NOT IN", "Not in", "notin", "NOTIN", "NotIn", "Notin":
		db = db.Where(key+" "+condition+" (?)", value)
	case "between", "BETWEEN":
		var betweenStr []string
		json.Unmarshal(conv.Bytes(value), &betweenStr)
		if len(betweenStr) > 1 {
			db = db.Where(key+" "+condition+" ? and ?", betweenStr[0], betweenStr[1])
		}
	case "<", "<=", ">", ">=", "=", "<>":
		var values = conv.String(value)
		if strings.Index("is null is not null", values) > -1 {
			db = db.Where(key + " " + values)
		} else {
			db = db.Where(key+" "+condition+" ?", values)
		}
	}
	return db
}

// or
func or(key, condition string, value interface{}, db *gorm.DB) *gorm.DB {
	switch condition {
	case "like", "LIKE", "Like", "notlike", "NOTLIKE", "Notlike", "ilike", "ILIKE", "Ilike", "rlike", "RLIKE", "Rlike":
		db = db.Or(key+" "+condition+" ?", conv.String(value)+"%")
	case "in", "IN", "In", "not in", "NOT IN", "Not in", "notin", "NOTIN", "NotIn", "Notin":
		db = db.Or(key+" "+condition+" (?)", value.([]interface{}))
	case "between", "BETWEEN", "Between":
		var betweenStr []string
		json.Unmarshal(conv.Bytes(value), &betweenStr)
		if len(betweenStr) > 1 {
			db = db.Or(key+" "+condition+" ? and ?", betweenStr[0], betweenStr[1])
		}
	case "<", "<=", ">", ">=", "=", "<>":
		var values = conv.String(value)
		if strings.Index("is null is not null", values) > -1 {
			db = db.Or(key + " " + values)
		} else {
			db = db.Or(key+" "+condition+" ?", values)
		}
	}
	return db
}
