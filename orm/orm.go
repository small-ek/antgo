package orm

import (
	"encoding/json"
	. "github.com/small-ek/ginp/conv"
	"gorm.io/gorm"
	"strings"
)

//Like Fuzzy search when there is value
func Like(key, value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if key != "" && value != "" {
			return db.Where(key+" LIKE ?", value+"%")
		}
		return db
	}
}

//Ilike Fuzzy search when there is value
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
		switch value := value.(type) {
		case []int:
			if len(value) == 0 {
				return db
			}
		case []int32:
			if len(value) == 0 {
				return db
			}
		case []int64:
			if len(value) == 0 {
				return db
			}
		case []uint32:
			if len(value) == 0 {
				return db
			}
		case []uint64:
			if len(value) == 0 {
				return db
			}
		case []string:
			if len(value) == 0 {
				return db
			}
		}

		if key != "" && value != nil && value != "" {
			return db.Where(""+key+" IN (?)", value)
		}
		return db
	}
}

// Where Where to search when there is value
func Where(key, conditions string, value interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if key != "" && conditions != "" && value != "" {
			return db.Where(""+key+" "+conditions+" ?", value)
		}

		return db
	}
}

// WhereBuildQuery Build a where query
func WhereBuildQuery(where interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if where == nil {
			return db
		}

		whereArray := where.([]string)
		var column string
		var value []interface{}
		for _, v := range whereArray {
			var arr []interface{}
			json.Unmarshal([]byte(v), &arr)

			if len(arr) == 3 && arr[2] != "" {

				//判断是否需要拼接
				if column != "" {
					column = column + " AND "
				}

				//检索where条件
				if arr[1] == "like" || arr[1] == "notlike" || arr[1] == "ilike" || arr[1] == "rlike" {
					column = column + String(arr[0]) + " " + String(arr[1]) + " ?"
					value = append(value, String(arr[2])+"%")

				} else if arr[1] == "between" && arr[2] != "" { //搜索between
					var betweenStr []string
					json.Unmarshal(Bytes(arr[2]), &betweenStr)
					if len(betweenStr) > 1 {
						column = column + String(arr[0]) + " BETWEEN ? AND ?"
						value = append(value, betweenStr[0], betweenStr[1])
					}

				} else if strings.Index(" in not in", String(arr[1])) > -1 {
					column = column + String(arr[0]) + " " + String(arr[1]) + " (?)"
					value = append(value, arr[2])

				} else {
					column = column + String(arr[0]) + " " + String(arr[1]) + " ?"
					value = append(value, arr[2])
				}
			} else if strings.Index("is null is not null", String(arr[1])) > -1 {
				if column != "" {
					column = column + " AND "
				}
				column = column + String(arr[0]) + " " + String(arr[1])
			}
		}

		return db.Where(column, value...)
	}
}

// Order Sort
func Order(str string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if str == "" {
			return db
		} else {
			return db.Order(str)
		}
	}
}

// Paginate...
func Paginate(page_size, current_page int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(page_size).Offset((current_page - 1) * page_size)
	}
}
